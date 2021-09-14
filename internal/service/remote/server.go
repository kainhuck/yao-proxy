package remote

import (
	"encoding/binary"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"net"
	"os"
	"strconv"
)

type Server struct {
	logger    log.Logger
	localAddr string
	cipher    *YPCipher.Cipher
}

func NewServer(localAddr string, logger log.Logger, cipher *YPCipher.Cipher) *Server {
	return &Server{
		logger:    logger,
		localAddr: localAddr,
		cipher:    cipher,
	}
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", s.localAddr)
	if err != nil {
		s.logger.Errorf("listen failed: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				s.logger.Errorf("accept failed: %v", err)
				continue
			}

			go s.handleConn(conn)
		}
	}()

	s.logger.Infof("listen on %v success", lis.Addr())
	select {}
}

func (s *Server) handleConn(conn net.Conn) {
	localConn := YPConn.NewConn(conn, s.cipher.Copy())
	defer func() {
		_ = localConn.Close()
	}()

	// 1. 读出target地址
	host, err := s.getTargetAddr(localConn)
	if err != nil {
		s.logger.Errorf("getTargetAddr error: %v", err)
		return
	}

	// 2. 和目标地址建立连接
	targetConn, err := net.Dial("tcp", host)
	if err != nil {
		s.logger.Errorf("dial target error: %v", err)
		return
	}

	defer func() {
		_ = targetConn.Close()
	}()

	// 3. 转发targetConn和localConn之间的数据
	YPConn.Forward(targetConn, localConn)
}

// 获取目标地址 string 类型
func (s *Server) getTargetAddr(conn *YPConn.Conn) (string, error) {
	buff := make([]byte, 256)

	_, err := conn.Read(buff)
	if err != nil {
		return "", err
	}

	length := buff[0] // ip + 端口(2) 的长度

	// 判断类型
	var host string
	switch buff[1] {
	case 1: // IPv4
		host = net.IP(buff[2:length]).String()
	case 3: // Domain
		host = string(buff[2:length])
	case 4: // IPv6
		host = net.IP(buff[2:length]).String()
	}

	port := binary.BigEndian.Uint16(buff[length : length+2])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))

	return host, nil
}
