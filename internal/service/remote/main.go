package remote

import (
	"flag"
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"github.com/kainhuck/yao-proxy/internal/log"
	"os"
	"sync"
)

func Main() {
	// 参数
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	cfg := ReadConfig(configFile)

	logger := log.NewLogger(cfg.Debug)

	var wg sync.WaitGroup
	for _, info := range cfg.ServerInfos {
		localAddr := fmt.Sprintf(":%d", info.Port)
		cipher, err := YPCipher.NewCipher(info.Method, info.Key)
		if err != nil {
			logger.Errorf("new cipher error: %v", err)
			os.Exit(1)
		}

		server := NewServer(localAddr, logger, cipher)
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.Run()
		}()
	}

	wg.Wait()
}
