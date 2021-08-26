package local

import (
	"fmt"
	"log"
	"math"
	"net"
	"time"
	YPCipher "yao-proxy/internal/cipher"
	YPConn "yao-proxy/internal/conn"
	YPPdu "yao-proxy/internal/pdu"
)

var cipher YPCipher.Cipher
var remoteAddr string

func Main() {
	var err error
	// 本地启动一个服务用于接收来自浏览器的请求

	// 参数 todo 后期改成从配置文件或环境变量中读取
	port := 20808
	key := []byte("1234567890qwerty")
	remoteHost := "127.0.0.1"
	remotePort := 20807

	remoteAddr = fmt.Sprintf("%s:%d", remoteHost, remotePort)
	cipher, err = YPCipher.NewCipher(key)
	if err != nil {
		log.Fatalf("[ERROR] new cipher error: %v", err)
	}

	// 启动服务
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("[ERROR] listen failed: %v", err)
	}

	go func() {
		var id uint32
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Printf("[ERROR] accept failed: %v", err)
				continue
			}

			go handleConn(&YPConn.Conn{
				Id:   id,
				Conn: conn,
			})
			if id < math.MaxUint32 {
				id++
			} else {
				id = 0
			}
		}
	}()

	log.Printf("[INFO] listen on %v success", lis.Addr())
	select {}
}

// 处理这个请求
func handleConn(conn *YPConn.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	// 1. 处理发过来的socks5请求, 并返回响应
	type_, host, err := handShake(conn)
	if err != nil {
		log.Printf("[ERROR] handshake error: %v", err)
		return
	}

	// 2. 解析出真正的IP，组装成我们的格式加密后发送给远程 host = ip + port
	cipherHost, err := cipher.Encrypt(host)
	if err != nil {
		log.Printf("[ERROR] Encrypt failed: %v", err)
		return
	}

	pdu := YPPdu.NewPDU(conn.Id, type_, cipherHost)
	bts, err := pdu.Encode()
	if err != nil {
		log.Printf("[ERROR] encode error: %v", bts)
		return
	}

	// 发给远程 地址信息
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("[ERROR] dial remote error: %v", err)
		return
	}
	_, err = remoteConn.Write(bts)
	if err != nil {
		log.Printf("[ERROR] write to remote error: %v", err)
		return
	}

	log.Printf("[INFO] 成功将数据信息发送给远程代理服务器")

	// 3. 将远程服务器发给我们的数据解密后转发给浏览器
	go sendToBrowser(remoteConn, conn)

	// 访问成功需要给客户端返回响应 假设访问成功
	/*
		   returns a reply formed as follows:
			+----+-----+-------+------+----------+----------+
			|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
			+----+-----+-------+------+----------+----------+
			| 1  |  1  | X'00' |  1   | Variable |    2     |
			+----+-----+-------+------+----------+----------+
	*/
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Printf("[ERROR] write back error: %v", err)
		return
	}

	// 4. 将浏览器的数据发送给远程
	sendToRemote(conn, remoteConn)
}

// handShake 和浏览器握手：处理socks5请求
func handShake(conn *YPConn.Conn) (uint8, []byte, error) {
	data, err := YPConn.Read(conn, 300*time.Second)
	if err != nil {
		return 0, nil, err
	}
	// 第一个包 用于选择协商方法
	/*
	   +----+----------+----------+
	   |VER | NMETHODS | METHODS  |
	   +----+----------+----------+
	   | 1  |    1     | 1 to 255 |
	   +----+----------+----------+
	*/
	if data[0] != 0x05 { // 必须确保是socks5协议
		return 0, nil, fmt.Errorf("unSupport socks version: %d", data[0])
	}
	// 给浏览器返回 0 表示无需验证
	/*
	 +----+--------+
	 |VER | METHOD |
	 +----+--------+
	 | 1  |   1    |
	 +----+--------+

	          o  X'00' NO AUTHENTICATION REQUIRED
	          o  X'01' GSSAPI
	          o  X'02' USERNAME/PASSWORD
	          o  X'03' to X'7F' IANA ASSIGNED
	          o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	          o  X'FF' NO ACCEPTABLE METHODS

	*/
	_, err = conn.Write([]byte{5, 0})
	if err != nil {
		return 0, nil, err
	}
	// 读取浏览器发送的第二个包 [5 1 0 1 104 16 249 249 1 187]
	/*
		        +----+-----+-------+------+----------+----------+
		        |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		        +----+-----+-------+------+----------+----------+
		        | 1  |  1  | X'00' |  1   | Variable |    2     |
		        +----+-----+-------+------+----------+----------+
			从中解析出地址和端口，其他字段暂不验证
	*/
	data, err = YPConn.Read(conn, 300*time.Second)
	if err != nil {
		return 0, nil, err
	}
	/*
		ATYP   address type of following address
			 o  IP V4 address: X'01'
			 o  DOMAINNAME: X'03'
			 o  IP V6 address: X'04'
	*/

	return data[3], data[4:], nil
}

// 将远程数据发送给浏览器
func sendToBrowser(remoteConn net.Conn, browserConn *YPConn.Conn) {
	for {
		cipherData, err := YPConn.Read(remoteConn, 300*time.Second)
		if err != nil {
			return
		}
		log.Printf("[INFO] 成功收到来自远程代理服务器的消息")
		data, err := cipher.Decrypt(cipherData)
		if err != nil {
			log.Printf("[ERROR] Decrypt error: %v", err)
			return
		}

		_, err = browserConn.Write(data)
		if err != nil {
			return
		}
		log.Printf("[INFO] 成功将远程代理收到的消息发送给本地浏览器")
	}
}

// 将浏览器数据发送给远程
func sendToRemote(browserConn *YPConn.Conn, remoteConn net.Conn) {
	for {
		data, err := YPConn.Read(browserConn, 300*time.Second)
		if err != nil {
			return
		}
		log.Printf("[INFO] 成功读取到浏览器的数据:%v", data)
		cipherData, err := cipher.Encrypt(data)
		if err != nil {
			log.Printf("[ERROR] Encrypt error: %v", err)
			return
		}

		pdu := YPPdu.NewPDU(browserConn.Id, YPPdu.DIRECT, cipherData)
		bts, err := pdu.Encode()
		if err != nil {
			log.Printf("[ERROR] pdu encode error: %v", err)
			return
		}

		_, err = remoteConn.Write(bts)
		if err != nil {
			return
		}
		log.Printf("[INFO] 成功将浏览器的数据转发给远程服务器")
	}
}
