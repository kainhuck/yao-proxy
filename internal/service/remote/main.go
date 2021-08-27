package remote

import (
	"fmt"
	"log"
	"net"
	YPCipher "yao-proxy/internal/cipher"
)

var cipher YPCipher.Cipher

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

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Printf("[ERROR] accept failed: %v", err)
				continue
			}

			job := NewJob(conn, cipher)
			go job.Run()
		}
	}()

	log.Printf("[INFO] listen on %v success", lis.Addr())
	select {}
}
