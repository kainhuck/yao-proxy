package middle

import (
	"context"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"net"
	"os"
	"time"
)

// NextNode 下一条节点信息
type NextNode struct {
	cipher *YPCipher.Cipher
	addr   string
	next   *NextNode // 循环链表，链接下一个
}

type Server struct {
	ctx       context.Context
	logger    log.Logger
	localAddr string
	nextNodes *NextNode
	index     int
	crLength  int
	pool      chan *YPConn.Conn
	cipher    *YPCipher.Cipher
}

func NewServer(ctx context.Context, localAddr string, logger log.Logger, next []NextNodeInfo, cipher *YPCipher.Cipher) *Server {
	s := &Server{
		ctx:       ctx,
		logger:    logger,
		localAddr: localAddr,
		index:     0,
		crLength:  len(next),
		pool:      make(chan *YPConn.Conn, 10),
		cipher:    cipher,
	}

	ci, err := YPCipher.NewCipher(next[0].Method, next[0].Key)
	if err != nil {
		s.logger.Errorf("new cipher error: %v", err)
		os.Exit(1)
	}

	node := &NextNode{
		cipher: ci,
		addr:   next[0].Addr,
		next:   nil,
	}

	s.nextNodes = node

	for i := 1; i < s.crLength; i++ {
		ci, err := YPCipher.NewCipher(next[i].Method, next[i].Key)
		if err != nil {
			s.logger.Errorf("new cipher error: %v", err)
			os.Exit(1)
		}
		node.next = &NextNode{
			cipher: ci,
			addr:   next[i].Addr,
			next:   nil,
		}

		node = node.next
	}

	node.next = s.nextNodes

	return s
}

func (s *Server) Run() {
	// 启动服务
	lis, err := net.Listen("tcp", s.localAddr)
	if err != nil {
		s.logger.Errorf("listen failed: %v", err)
		os.Exit(1)
	}

	// 事先和远程建立连接
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case s.pool <- s.getNextNodeConn():
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				conn, err := lis.Accept()
				if err != nil {
					s.logger.Errorf("accept failed: %v", err)
					continue
				}

				go s.handleConn(conn)
			}
		}
	}()

	s.logger.Infof("listen on %v success", lis.Addr())
	<-s.ctx.Done()
}

// handleConn 处理链接
func (s *Server) handleConn(conn net.Conn) {
	aheadConn := YPConn.NewConn(conn, s.cipher.Copy())
	defer func() {
		_ = aheadConn.Close()
	}()

	// 和远程建立链接，直接转发给远程
	remoteConn := new(YPConn.Conn)
	select {
	case remoteConn = <-s.pool:
	case <-time.After(10 * time.Second):
		s.logger.Errorf("dial remote time out")
		return
	}

	YPConn.Forward(remoteConn, aheadConn)
}

// getNextNode 从下一个节点链表中顺序选择一个
func (s *Server) getNextNode() *NextNode {
	s.nextNodes = s.nextNodes.next

	return s.nextNodes
}

// getNextNodeConn 获取下一个节点的链接
func (s *Server) getNextNodeConn() *YPConn.Conn {
	node := s.getNextNode()
	conn, err := YPConn.Dial(node.addr, node.cipher.Copy(), 10*time.Second)
	if err != nil {
		// 如果当前的节点不可用，就递归拿下一个节点  todo 这种做法是否有隐患?
		return s.getNextNodeConn()
	}

	return conn
}
