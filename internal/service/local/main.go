package local

import (
	"flag"
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"github.com/kainhuck/yao-proxy/internal/utils"
	"io"
	"net"
	"os"
	"time"
)

var cipher YPCipher.Cipher
var remoteAddr string
var logger log.Logger
var localAddr string

func init() {
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	// 参数
	cfg, err := ReadConfig(configFile)
	if err != nil {
		logger.Errorf("read config file error")
		os.Exit(1)
	}

	remoteAddr = fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort)
	localAddr = fmt.Sprintf(":%d", cfg.Port)

	cipher, err = YPCipher.NewCipher(utils.MD5(cfg.Key))
	if err != nil {
		logger.Errorf("new cipher error: %v", err)
		os.Exit(1)
	}

	logger = log.NewLogger(cfg.Debug)
}

func Main() {
	// 启动服务
	lis, err := net.Listen("tcp", localAddr)
	if err != nil {
		logger.Errorf("listen failed: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				logger.Errorf("accept failed: %v", err)
				continue
			}

			go handleConn(conn)
		}
	}()

	logger.Infof("listen on %v success", lis.Addr())
	select {}
}

func handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetReadDeadline(time.Now().Add(600 * time.Second))
	// 1. 握手
	if err := handShake(conn); err != nil {
		logger.Errorf("handShake error: %v", err)
		return
	}
	// 2. 获取真实的地址
	addr, err := getTargetAddr(conn)
	if err != nil {
		logger.Errorf("getTargetAddr error: %v", err)
		return
	}

	// 3. 给浏览器发送成功响应
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



		Leech, et al                Standards Track                     [Page 5]

		RFC 1928                SOCKS Protocol Version 5              March 1996


		             o  IP V4 address: X'01'
		             o  DOMAINNAME: X'03'
		             o  IP V6 address: X'04'
		          o  BND.ADDR       server bound address
		          o  BND.PORT       server bound port in network octet order
	*/
	if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil {
		logger.Errorf("reply to browser error: %v", err)
		return
	}

	// 4. 和远程建立链接并将目标地址发送给远程
	remoteConn, err := YPConn.DialAndSend(remoteAddr, cipher, addr)
	if err != nil {
		logger.Errorf("DialAndSend error: %v", err)
		return
	}
	defer func() {
		_ = remoteConn.Close()
	}()

	// 5. 将RemoteConn的数据和conn的数据进行转发
	errChan := make(chan error, 2)
	go func() {
		errChan <- YPConn.Copy(remoteConn, conn)
	}()
	go func() {
		errChan <- YPConn.Copy(conn, remoteConn)
	}()

	select {
	case <-errChan:
		return
	}
}

func handShake(conn net.Conn) error {
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

func getTargetAddr(conn net.Conn) ([]byte, error) {
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
		return nil, err
	}

	var addr []byte
	// 判断类型
	switch head[3] {
	case 1: // IPV4
		addr = make([]byte, net.IPv4len+4)
		addr[0] = net.IPv4len + 2
		addr[1] = 1
		addr[2] = head[4]
		if _, err := io.ReadFull(conn, addr[3:]); err != nil {
			return nil, err
		}
	case 3: // Domain
		addr = make([]byte, head[4]+4)
		addr[0] = head[4] + 2
		addr[1] = 3
		if _, err := io.ReadFull(conn, addr[2:]); err != nil {
			return nil, err
		}
	case 4: // IPV6
		addr = make([]byte, net.IPv6len+4)
		addr[0] = net.IPv6len + 2
		addr[1] = 1
		addr[2] = head[4]
		if _, err := io.ReadFull(conn, addr[3:]); err != nil {
			return nil, err
		}
	}

	return addr, nil
}
