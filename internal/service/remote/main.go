package remote

import (
	"flag"
	"fmt"
	"log"
	"net"
	YPCipher "yao-proxy/internal/cipher"
	"yao-proxy/internal/utils"
)

var cipher YPCipher.Cipher

func Main() {
	var err error
	// 参数
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	cfg := ReadConfig(configFile)

	cipher, err = YPCipher.NewCipher(utils.MD5(cfg.Key))
	if err != nil {
		log.Fatalf("[ERROR] new cipher error: %v", err)
	}

	// 启动服务
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
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
