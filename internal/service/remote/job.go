package remote

import (
	"encoding/binary"
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	YPConn "github.com/kainhuck/yao-proxy/internal/conn"
	"github.com/kainhuck/yao-proxy/internal/log"
	YPPdu "github.com/kainhuck/yao-proxy/internal/pdu"
	"net"
	"time"
)

// Job 每接收一个本地代理的请求，就开启一个任务
//     任务中包含两个链接，一个和目标网站链接，一个和本地代理链接
type Job struct {
	LocalConn  *YPConn.Conn // 和本地代理链接
	TargetConn net.Conn     // 和目标网站链接

	logger  log.Logger
	timeout time.Duration
	ci      YPCipher.Cipher
}

func NewJob(c net.Conn, ci YPCipher.Cipher, debug bool) *Job {
	return &Job{
		LocalConn: YPConn.NewConn(c),
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
	data, err := YPConn.Read(j.LocalConn.Conn, j.timeout)
	if err != nil {
		j.logger.Errorf("read from local error: %v", err)
		return
	}

	pdu := &YPPdu.PDU{}
	err = pdu.Decode(data)
	if err != nil {
		j.logger.Errorf("Decode error: %v", err)
		return
	}

	host, err := j.ci.Decrypt(pdu.Data)
	if err != nil {
		j.logger.Errorf("Decrypt error: %v", err)
		return
	}

	var addr *net.TCPAddr
	switch pdu.Type {
	case YPPdu.IPv4:
		addr = &net.TCPAddr{
			IP:   host[:4],
			Port: int(binary.BigEndian.Uint16(host[4:])),
		}
	case YPPdu.IPv6:
		addr = &net.TCPAddr{
			IP:   host[:16],
			Port: int(binary.BigEndian.Uint16(host[16:])),
		}
	case YPPdu.DOMAIN:
		ipAddr, err := net.ResolveIPAddr("ip", string(host[:len(host)-2]))
		if err != nil {
			return
		}
		addr = &net.TCPAddr{
			IP:   ipAddr.IP,
			Port: int(binary.BigEndian.Uint16(host[len(host)-2:])),
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

	j.logger.Debugf("dial remote success")

	go func() {
		for {
			select {
			case data := <-j.LocalConn.CDataChan:
				rawData, err := j.ci.Decrypt(data)
				if err != nil {
					return
				}
				j.logger.Debugf("read from local success")
				_, err = j.TargetConn.Write(rawData)
				if err != nil {
					return
				}
				j.logger.Debugf("send to target success")
			}
		}
	}()

	for {
		data, err := YPConn.Read(j.TargetConn, j.timeout)
		if err != nil {
			return
		}
		j.logger.Debugf("read from target success")
		cData, err := j.ci.Encrypt(data)
		if err != nil {
			return
		}

		err = j.LocalConn.Write(0, cData)
		if err != nil {
			return
		}
		j.logger.Debugf("send to local success")
	}
}
