package local

import (
	"fmt"
	"log"
	"net"
	YPCipher "yao-proxy/internal/cipher"
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
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Printf("[ERROR] accept failed: %v", err)
				continue
			}

			job, err := NewJob(conn, remoteAddr, cipher)
			if err != nil {
				continue
			}

			go job.Run()
		}
	}()

	log.Printf("[INFO] listen on %v success", lis.Addr())
	select {}
}
