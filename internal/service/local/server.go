package local

import (
	"context"
	"encoding/binary"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type CipherRemote struct {
	cipher     *YPCipher.Cipher
	RemoteAddr string
	next       *CipherRemote
}

type Server struct {
	ctx           context.Context
	logger        log.Logger
	localAddr     string
	cipherRemotes *CipherRemote
	index         int
	crLength      int
	remotePool    chan *YPConn.Conn
	filter        *Filter
}

func NewServer(ctx context.Context, localAddr string, logger log.Logger, infos []RemoteInfo, filter *Filter) *Server {
	s := &Server{
		ctx:        ctx,
		logger:     logger,
		localAddr:  localAddr,
		crLength:   len(infos),
		index:      0,
		remotePool: make(chan *YPConn.Conn, 10),
		filter:     filter,
	}

	// infos 不可能为 0
	cipher, err := YPCipher.NewCipher(infos[0].Method, infos[0].Key)
	if err != nil {
		s.logger.Errorf("new cipher error: %v", err)
		os.Exit(1)
	}

	cr := &CipherRemote{
		cipher:     cipher,
		RemoteAddr: infos[0].RemoteAddr,
		next:       nil,
	}

	s.cipherRemotes = cr

	for i := 1; i < s.crLength; i++ {
		cipher, err := YPCipher.NewCipher(infos[i].Method, infos[i].Key)
		if err != nil {
			s.logger.Errorf("new cipher error: %v", err)
			os.Exit(1)
		}
		cr.next = &CipherRemote{
			cipher:     cipher,
			RemoteAddr: infos[i].RemoteAddr,
			next:       nil,
		}

		cr = cr.next
	}

	cr.next = s.cipherRemotes

	return s
}

func (s *Server) Run() {
	// 启动服务
	lis, err := net.Listen("tcp", s.localAddr)
	if err != nil {
		s.logger.Errorf("listen failed: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case s.remotePool <- s.getRemoteConn():
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

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	// 1. 握手
	if err := s.handShake(conn); err != nil {
		s.logger.Errorf("handShake error: %v", err)
		return
	}
	// 2. 获取真实的地址
	addr, host, err := s.getTargetAddr(conn)
	if err != nil {
		s.logger.Errorf("getTargetAddr error: %v", err)
		return
	}

	// 3. 给浏览器发送成功响应
	go func() {
		/*
		 The SOCKS request information is sent by the client as soon as it has
		   established a connection to the SOCKS server, and completed the
		   authentication negotiations.  The server evaluates the request, and
		   returns a reply formed as follows:

		        +----+-----+-------+------+----------+----------+
		        |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		        +----+-----+-------+------+----------+----------+
		        | 1  |  1  | X'00' |  1   | Variable |    2     |
		        +----+-----+-------+------+----------+----------+

		     Where:

		          o  VER    protocol version: X'05'
		          o  REP    Reply field:
		             o  X'00' succeeded
		             o  X'01' general SOCKS server failure
		             o  X'02' connection not allowed by ruleset
		             o  X'03' Network unreachable
		             o  X'04' Host unreachable
		             o  X'05' Connection refused
		             o  X'06' TTL expired
		             o  X'07' Command not supported
		             o  X'08' Address type not supported
		             o  X'09' to X'FF' unassigned
		          o  RSV    RESERVED
		          o  ATYP   address type of following address
		             o  IP V4 address: X'01'
		             o  DOMAINNAME: X'03'
		             o  IP V6 address: X'04'
		          o  BND.ADDR       server bound address
		          o  BND.PORT       server bound port in network octet order
		*/
		if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil {
			s.logger.Errorf("reply to browser error: %v", err)
			return
		}
	}()

	if !s.filter.Check(host) {
		// 4. 和远程建立链接并将目标地址发送给远程
		remoteConn := new(YPConn.Conn)
		select {
		case remoteConn = <-s.remotePool:
		case <-time.After(10 * time.Second):
			s.logger.Errorf("dial remote time out")
			return
		}

		_, err = remoteConn.Write(addr)
		if err != nil {
			s.logger.Errorf("DialAndSend error: %v", err)
			return
		}
		defer func() {
			_ = remoteConn.Close()
		}()

		// 5. 将RemoteConn的数据和conn的数据进行转发
		YPConn.Forward(remoteConn, conn)
	} else {
		// 直接访问目标地址

		targetConn, err := net.Dial("tcp", host)
		if err != nil {
			s.logger.Errorf("dial target error: %v", err)
			return
		}

		defer func() {
			_ = targetConn.Close()
		}()

		// 3. 转发targetConn和localConn之间的数据
		YPConn.Forward(targetConn, conn)
	}
}

func (s *Server) handShake(conn net.Conn) error {
	// 和浏览器进行s5握手
	/*
		The client connects to the server, and sends a version
		   identifier/method selection message:

				   +----+----------+----------+
				   |VER | NMETHODS | METHODS  |
				   +----+----------+----------+
				   | 1  |    1     | 1 to 255 |
				   +----+----------+----------+
	*/
	buff := make([]byte, 257)
	_, err := conn.Read(buff)
	if err != nil {
		return err
	}

	/*
		The server selects from one of the methods given in METHODS, and
		   sends a METHOD selection message:

		                         +----+--------+
		                         |VER | METHOD |
		                         +----+--------+
		                         | 1  |   1    |
		                         +----+--------+

		   If the selected METHOD is X'FF', none of the methods listed by the
		   client are acceptable, and the client MUST close the connection.

		   The values currently defined for METHOD are:

		          o  X'00' NO AUTHENTICATION REQUIRED
		          o  X'01' GSSAPI
		          o  X'02' USERNAME/PASSWORD
		          o  X'03' to X'7F' IANA ASSIGNED
		          o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
		          o  X'FF' NO ACCEPTABLE METHODS
	*/
	_, err = conn.Write([]byte{5, 0})
	return err
}

func (s *Server) getTargetAddr(conn net.Conn) ([]byte, string, error) {
	/*
		The SOCKS request is formed as follows:

		        +----+-----+-------+------+----------+----------+
		        |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		        +----+-----+-------+------+----------+----------+
		        | 1  |  1  | X'00' |  1   | Variable |    2     |
		        +----+-----+-------+------+----------+----------+

		     Where:

		          o  VER    protocol version: X'05'
		          o  CMD
		             o  CONNECT X'01'
		             o  BIND X'02'
		             o  UDP ASSOCIATE X'03'
		          o  RSV    RESERVED
		          o  ATYP   address type of following address
		             o  IP V4 address: X'01'
		             o  DOMAINNAME: X'03'
		             o  IP V6 address: X'04'
		          o  DST.ADDR       desired destination address
		          o  DST.PORT desired destination port in network octet
		             order
	*/
	head := make([]byte, 5)
	if _, err := io.ReadFull(conn, head); err != nil {
		return nil, "", err
	}

	var addr []byte
	var host string
	// 判断类型
	switch head[3] {
	case 1: // IPV4
		addr = make([]byte, net.IPv4len+4)
		addr[0] = net.IPv4len + 2
		addr[1] = 1
		addr[2] = head[4]
		if _, err := io.ReadFull(conn, addr[3:]); err != nil {
			return nil, "", err
		}
		host = net.IP(addr[2:addr[0]]).String()
	case 3: // Domain
		addr = make([]byte, head[4]+4)
		addr[0] = head[4] + 2
		addr[1] = 3
		if _, err := io.ReadFull(conn, addr[2:]); err != nil {
			return nil, "", err
		}
		host = string(addr[2:addr[0]])
	case 4: // IPV6
		addr = make([]byte, net.IPv6len+4)
		addr[0] = net.IPv6len + 2
		addr[1] = 4
		addr[2] = head[4]
		if _, err := io.ReadFull(conn, addr[3:]); err != nil {
			return nil, "", err
		}
		host = net.IP(addr[2:addr[0]]).String()
	}

	port := binary.BigEndian.Uint16(addr[addr[0] : addr[0]+2])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))

	return addr, host, nil
}

// 顺序选择一个远程服务器
func (s *Server) getCipherRemote() *CipherRemote {
	s.cipherRemotes = s.cipherRemotes.next

	return s.cipherRemotes
}

func (s *Server) getRemoteConn() *YPConn.Conn {
	cr := s.getCipherRemote()
	conn, err := YPConn.Dial(cr.RemoteAddr, cr.cipher.Copy(), 10*time.Second)
	if err != nil {
		return s.getRemoteConn()
	}

	return conn
}
