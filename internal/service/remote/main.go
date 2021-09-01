package remote

import (
	"encoding/binary"
	"flag"
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"net"
	"os"
	"strconv"
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

// 获取目标地址 string 类型
func getTargetAddr(conn *YPConn.Conn) (string, error) {
	buff := make([]byte, 256)

	_, err := conn.Read(buff)
	if err != nil {
		return "", err
	}

	length := buff[0] // ip + 端口(2) 的长度

	// 判断类型
	var host string
	switch buff[1] {
	case 1: // IPv4
		host = net.IP(buff[2:length]).String()
	case 3: // Domain
		host = string(buff[2:length])
	case 4: // IPv6
		host = net.IP(buff[2:length]).String()
	}

	port := binary.BigEndian.Uint16(buff[length : length+2])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))

	return host, nil
}
