package local

import (
	"flag"
	"fmt"
	"github.com/kainhuck/yao-proxy/internal/log"
	"os"
)

func Main() {
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	// 参数
	cfg, err := ReadConfig(configFile)
	if err != nil {
		fmt.Printf("read config file error: %v", err)
		os.Exit(1)
	}

	logger := log.NewLogger(cfg.Debug)
	localAddr := fmt.Sprintf(":%d", cfg.Port)

	server := NewServer(localAddr, logger, cfg.RemoteInfos)
	server.Run()
}
