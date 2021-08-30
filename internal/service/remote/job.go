package remote

import (
	"encoding/binary"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	"io"
	"net"
	"time"
)

// Job 每接收一个本地代理的请求，就开启一个任务
//     任务中包含两个链接，一个和目标网站链接，一个和本地代理链接
type Job struct {
	LocalConn  net.Conn // 和本地代理链接
	TargetConn net.Conn // 和目标网站链接

	logger  log.Logger
	timeout time.Duration
	ci      YPCipher.Cipher
}

func NewJob(c net.Conn, ci YPCipher.Cipher, debug bool) *Job {
	return &Job{
		LocalConn: c,
		logger:    log.NewLogger(debug),
		timeout:   300 * time.Second,
		ci:        ci,
	}
}

func (j *Job) Run() {
	defer func() {
		err := recover()
		if err != nil {
			j.logger.Error(err)
		}
	}()
	defer func() {
		_ = j.LocalConn.Close()
	}()

	data, err := YPConn.DecryptRead(j.LocalConn, j.timeout, j.ci)
	if err != nil {
		j.logger.Errorf("read from local error: %v", err)
		return
	}

	if data[0] != 0x05 {
		return
	}

	err = YPConn.EncryptWrite(j.LocalConn, j.ci, []byte{5, 0})
	if err != nil {
		return
	}

	data, err = YPConn.DecryptRead(j.LocalConn, j.timeout, j.ci)
	if err != nil {
		return
	}

	var addr *net.TCPAddr
	switch data[3] {
	case 1:
		addr = &net.TCPAddr{
			IP:   data[4 : 4+net.IPv4len],
			Port: int(binary.BigEndian.Uint16(data[len(data)-2:])),
		}
	case 4:
		addr = &net.TCPAddr{
			IP:   data[4 : 4+net.IPv6len],
			Port: int(binary.BigEndian.Uint16(data[len(data)-2:])),
		}
	case 3:
		ipAddr, err := net.ResolveIPAddr("ip", string(data[4:len(data)-2]))
		if err != nil {
			return
		}
		addr = &net.TCPAddr{
			IP:   ipAddr.IP,
			Port: int(binary.BigEndian.Uint16(data[len(data)-2:])),
		}
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		j.logger.Errorf("dial remote error: %v", err)
		return
	}

	j.TargetConn = conn
	defer func() {
		_ = j.TargetConn.Close()
	}()
	err = YPConn.EncryptWrite(j.LocalConn, j.ci, []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return
	}
	j.logger.Debugf("dial remote success")

	errChan := make(chan error, 2)
	go func() { errChan <- YPConn.EncryptCopy(j.LocalConn, j.TargetConn, j.ci) }()
	go func() { errChan <- YPConn.DecryptCopy(j.TargetConn, j.LocalConn, j.ci) }()

	select {
	case err := <-errChan:
		if err != io.EOF {
			j.logger.Error(err)
		}
		return
	}
}
