package remote

import (
	"flag"
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"github.com/kainhuck/yao-proxy/internal/log"
	"os"
	"os/signal"
	"syscall"
)

func Main() {
	defer func() {
		fmt.Printf("[YAO-PROXY] remote agent exit successfully !")
	}()
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	cfg := ReadConfig(configFile)

	logger := log.NewLogger(cfg.Debug)

	for _, info := range cfg.ServerInfos {
		localAddr := fmt.Sprintf(":%d", info.Port)
		cipher, err := YPCipher.NewCipher(info.Method, info.Key)
		if err != nil {
			logger.Errorf("new cipher error: %v", err)
			os.Exit(1)
		}

		server := NewServer(localAddr, logger, cipher)
		go server.Run()
	}

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, os.Interrupt)

	<-stopSignalCh
}
