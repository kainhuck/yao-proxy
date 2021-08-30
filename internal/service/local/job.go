package local

import (
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
	BrowserConn net.Conn // 和浏览器的链接
	RemoteConn  net.Conn // 和远程服务器的链接

	logger  log.Logger
	timeout time.Duration
	ci      YPCipher.Cipher
}

func NewJob(c net.Conn, remoteAddr string, ci YPCipher.Cipher, debug bool) (*Job, error) {
	rc, err := net.Dial("tcp", remoteAddr)
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

	errChan := make(chan error, 2)
	go func() {
		errChan <- YPConn.EncryptCopy(j.RemoteConn, j.BrowserConn, j.ci)
	}()

	go func() {
		errChan <- YPConn.DecryptCopy(j.BrowserConn, j.RemoteConn, j.ci)
	}()

	select {
	case err := <-errChan:
		if err != io.EOF {
			j.logger.Error(err)
		}
		return
	}
}
