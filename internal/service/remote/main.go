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
)

var cipher *YPCipher.Cipher
var localAddr string
var logger log.Logger

func init() {
	var err error
	// 参数
	var configFile string
	flag.StringVar(&configFile, "c", "/etc/yao-proxy/config.json", "go run main.go -c configFile")
	flag.Parse()
	cfg := ReadConfig(configFile)

	cipher, err = YPCipher.NewCipher(cfg.Method, cfg.Key)
	if err != nil {
		logger.Errorf("new cipher error: %v", err)
		os.Exit(1)
	}

	localAddr = fmt.Sprintf(":%d", cfg.Port)

	logger = log.NewLogger(cfg.Debug)
}

func Main() {
	// 启动服务
	lis, err := net.Listen("tcp", localAddr)
	if err != nil {
		logger.Errorf("listen failed: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				logger.Errorf("accept failed: %v", err)
				continue
			}

			go handleConn(conn)
		}
	}()

	logger.Infof("listen on %v success", lis.Addr())
	select {}
}

func handleConn(conn net.Conn) {
	localConn := YPConn.NewConn(conn, cipher)
	defer func() {
		_ = localConn.Close()
	}()

	// 1. 读出target地址
	host, err := getTargetAddr(localConn)
	if err != nil {
		logger.Errorf("getTargetAddr error: %v", err)
		return
	}

	// 2. 和目标地址建立连接
	targetConn, err := net.Dial("tcp", host)
	if err != nil {
		logger.Errorf("dial remote error: %v", err)
		return
	}

	defer func() {
		_ = targetConn.Close()
	}()

	// 3. 转发targetConn和localConn之间的数据
	errChan := make(chan error, 2)
	go func() {
		errChan <- YPConn.Copy(targetConn, localConn)
	}()
	go func() {
		errChan <- YPConn.Copy(localConn, targetConn)
	}()

	select {
	case <-errChan:
		return
	}
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
