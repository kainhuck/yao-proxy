package local

import (
	"fmt"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"io"
	"net"
	"time"
)

// Job 每接收一个浏览器的请求，就开启一个任务
//     任务中包含两个链接，一个和浏览器链接，一个和远程链接
type Job struct {
	BrowserConn net.Conn     // 和浏览器的链接
	RemoteConn  *YPConn.Conn // 和远程服务器的链接

	logger  log.Logger
	timeout time.Duration
	ci      YPCipher.Cipher
}

func NewJob(c net.Conn, remoteAddr string, ci YPCipher.Cipher, debug bool) (*Job, error) {
	rc, err := YPConn.Dial(remoteAddr)
	if err != nil {
		return nil, err
	}

	return &Job{
		BrowserConn: c,
		RemoteConn:  rc,
		logger:      log.NewLogger(debug),
		timeout:     300 * time.Second,
		ci:          ci,
	}, nil
}

func (j *Job) Run() {
	defer func() {
		err := recover()
		if err != nil {
			j.logger.Error(err)
		}
	}()
	defer func() {
		_ = j.BrowserConn.Close()
		_ = j.RemoteConn.Close()
	}()
	// 1. 握手🤝
	t, host, err := j.HandShake()
	if err != nil {
		j.logger.Errorf("handshake error: %v", err)
		return
	}

	j.logger.Debugf("handshake success")

	cHost, err := j.ci.Encrypt(host)
	if err != nil {
		j.logger.Errorf("Encrypt error: %v", err)
		return
	}

	err = j.RemoteConn.Write(t, cHost)
	if err != nil {
		j.logger.Errorf("Write error: %v", err)
		return
	}

	// 2. 转发消息
	go func() {
		for {
			data, err := YPConn.Read(j.BrowserConn, j.timeout)
			if err != nil {
				return
			}
			j.logger.Debugf("Read from browser success")
			cData, err := j.ci.Encrypt(data)
			if err != nil {
				return
			}

			err = j.RemoteConn.Write(0, cData)
			if err != nil {
				return
			}
			j.logger.Debugf("send to remote success")
		}
	}()

	for {
		select {
		case data := <-j.RemoteConn.CDataChan:
			rawData, err := j.ci.Decrypt(data)
			if err != nil {
				return
			}
			j.logger.Debugf("read from remote success")
			_, err = j.BrowserConn.Write(rawData)
			if err != nil {
				return
			}
			j.logger.Debugf("send to browser success")
		}
	}
}

func (j *Job) HandShake() (uint8, []byte, error) {
	data, err := YPConn.Read(j.BrowserConn, j.timeout)
	if err != nil {
		if err != io.EOF {
			return 0, nil, nil
		}
	}

	if data[0] != 0x05 {
		return 0, nil, fmt.Errorf("only support socks5")
	}

	_, err = j.BrowserConn.Write([]byte{5, 0})
	if err != nil {
		return 0, nil, err
	}

	data, err = YPConn.Read(j.BrowserConn, j.timeout)
	if err != nil {
		if err != io.EOF {
			return 0, nil, err
		}
	}

	_, err = j.BrowserConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return 0, nil, err
	}

	return data[3], data[4:], nil
}
