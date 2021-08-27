package local

import (
	"flag"
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"github.com/kainhuck/yao-proxy/internal/utils"
	"log"
	"net"
)

var cipher YPCipher.Cipher
var remoteAddr string

func Main() {
	var err error
	// 本地启动一个服务用于接收来自浏览器的请求

	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	// 参数
	cfg, err := ReadConfig(configFile)
	if err != nil {
		log.Fatalf("config file error")
	}

	remoteAddr = fmt.Sprintf("%s:%d", cfg.RemoteHost, cfg.RemotePort)
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

			job, err := NewJob(conn, remoteAddr, cipher, cfg.Debug)
			if err != nil {
				continue
			}

			go job.Run()
		}
	}()

	log.Printf("[INFO] listen on %v success", lis.Addr())
	select {}
}
