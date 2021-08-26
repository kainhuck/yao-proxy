package remote

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
	YPCipher "yao-proxy/internal/cipher"
	YPConn "yao-proxy/internal/conn"
	YPPdu "yao-proxy/internal/pdu"
)

var cipher YPCipher.Cipher

var broker map[uint32]chan []byte

func Main() {
	var err error
	// 参数 todo 后期改成从配置文件或环境变量中读取
	port := 20807
	key := []byte("1234567890qwerty")

	cipher, err = YPCipher.NewCipher(key)
	if err != nil {
		log.Fatalf("[ERROR] new cipher error: %v", err)
	}

	// 启动服务
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("[ERROR] listen failed: %v", err)
	}

	broker = make(map[uint32]chan []byte)

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Printf("[ERROR] accept failed: %v", err)
				continue
			}

			go handleConn(conn)
		}
	}()

	log.Printf("[INFO] listen on %v success", lis.Addr())
	select {}
}

func handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	bts, err := YPConn.Read(conn, 300*time.Second)
	if err != nil {
		log.Printf("[ERROR] failed to read conn: %v", err)
		return
	}

	pdu := &YPPdu.PDU{}
	err = pdu.Decode(bts)
	if err != nil {
		log.Printf("[ERROR] error pdu: %v", err)
		return
	}

	data, err := cipher.Decrypt(pdu.Data)
	if err != nil {
		log.Printf("[ERROR] Decrypt failed: %v", err)
		return
	}

	// 先解析type
	var addr *net.TCPAddr
	switch pdu.Type {
	case YPPdu.IPv4:
		// 解析出ipv4地址
		if len(data) != 6 {
			log.Printf("[ERROR] error data length for ipv4, length: %v", len(data))
			return
		}
		addr = &net.TCPAddr{
			IP:   data[:4],
			Port: int(binary.BigEndian.Uint16(data[4:])),
		}
	case YPPdu.IPv6:
		if len(data) != 18 {
			log.Printf("[ERROR] error data length for ipv6, length: %v", len(data))
			return
		}
		addr = &net.TCPAddr{
			IP:   data[:16],
			Port: int(binary.BigEndian.Uint16(data[16:])),
		}
	case YPPdu.DOMAIN:
		if len(data) < 3 {
			log.Printf("[ERROR] error data length for Domain, length: %v", len(data))
			return
		}
		ipAddr, err := net.ResolveIPAddr("ip", string(data[:len(data)-2]))
		if err != nil {
			log.Printf("[ERROR] ResolveIPAddr error: %v", err)
			return
		}
		addr = &net.TCPAddr{
			IP:   ipAddr.IP,
			Port: int(binary.BigEndian.Uint16(data[len(data)-2:])),
		}
	}

	// 和远程建立链接
	targetConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Printf("和远程建立链接失败")
	}
	defer func() {
		_ = targetConn.Close()
	}()

	go func() {
		// 将真正的数据加密后发给 localConn
		for {
			bts, err := YPConn.Read(targetConn, 300*time.Second)
			if err != nil {
				log.Printf("[ERROR] READ from Target error: %v", err)
				return
			}
			log.Printf("[INFO] 成功从目标网站读到数据")
			cipherData, err := cipher.Encrypt(bts)
			if err != nil {
				log.Printf("[ERROR] encrypt error: %v", err)
				return
			}
			_, err = conn.Write(cipherData)
			if err != nil {
				log.Printf("[ERROR] write to local error: %v", err)
				return
			}

			log.Printf("[INFO] 成功将数据转发给本地端")
		}
	}()

	for {
		bts, err := YPConn.Read(conn, 300*time.Second)
		if err != nil {
			log.Printf("[ERROR] read error: %v", err)
			return
		}
		pdu := &YPPdu.PDU{}
		err = pdu.Decode(bts)
		if err != nil {
			log.Printf("[ERROR] pdu error: %v", err)
			return
		}

		data, err := cipher.Decrypt(pdu.Data)
		if err != nil {
			log.Printf("[ERROR] Decrypt error: %v", err)
			return
		}

		_, err = targetConn.Write(data)
		if err != nil {
			log.Printf("[ERROR] Write error: %v", err)
			return
		}
	}
}
