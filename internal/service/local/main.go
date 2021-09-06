package local

import (
	"flag"
	"fmt"
	"github.com/kainhuck/yao-proxy/internal/log"
	"os"
	"sync"
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
	var wg sync.WaitGroup

	for _, info := range cfg.ServerInfos {
		localAddr := fmt.Sprintf(":%d", info.Port)

		server := NewServer(localAddr, logger, info.RemoteInfos)

		wg.Add(1)

		go func() {
			defer wg.Done()
			server.Run()
		}()
	}

	wg.Wait()
}
