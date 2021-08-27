package conn

import (
	"encoding/binary"
	YPPdu "github.com/kainhuck/yao-proxy/internal/pdu"
	"io"
	"net"
)

// Conn 这个链接用于处理本地和远程之间的通信
type Conn struct {
	net.Conn

	CDataChan chan []byte
}

func NewConn(c net.Conn) *Conn {
	cc := &Conn{
		Conn:      c,
		CDataChan: make(chan []byte, 3),
	}
	go cc.read()
	return cc
}

func Dial(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		Conn:      conn,
		CDataChan: make(chan []byte, 3),
	}
	go c.read()
	return c, nil
}

func (c *Conn) Write(type_ uint8, data []byte) error {

	pdu := YPPdu.NewPDU(0, type_, data)

	bts, err := pdu.Encode()
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(bts)
	return err
}

func (c *Conn) read() {
	for {
		head := make([]byte, 9)
		_, err := io.ReadFull(c.Conn, head)
		if err != nil {
			return
		}

		length := int(binary.BigEndian.Uint16(head[7:]))
		data := make([]byte, length+2)
		_, err = io.ReadFull(c.Conn, data)
		if err != nil {
			return
		}
		c.CDataChan <- data[:length]
	}
}
