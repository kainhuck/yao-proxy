package conn

import (
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"io"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	*YPCipher.Cipher
	readBuf  []byte
	writeBuf []byte
}

func NewConn(conn net.Conn, cipher *YPCipher.Cipher) *Conn {
	return &Conn{
		Conn:     conn,
		Cipher:   cipher,
		readBuf:  bufferPoolGet(),
		writeBuf: bufferPoolGet(),
	}
}

func (c *Conn) Close() error {
	bufferPoolPut(c.readBuf)
	bufferPoolPut(c.writeBuf)
	return c.Conn.Close()
}

func DialAndSend(addr string, cipher *YPCipher.Cipher, data []byte) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := NewConn(conn, cipher)

	_, err = c.Write(data)

	return c, err
}

func Dial(addr string, cipher *YPCipher.Cipher, timeout time.Duration) (*Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)

	return NewConn(conn, cipher), err
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.Dec == nil {
		iv := make([]byte, c.Info.IvLen)
		if _, err = io.ReadFull(c.Conn, iv); err != nil {
			return
		}
		if err = c.InitDecrypt(iv); err != nil {
			return
		}
	}

	cipherData := c.readBuf
	if len(b) > len(cipherData) {
		cipherData = make([]byte, len(b))
	} else {
		cipherData = cipherData[:len(b)]
	}

	n, err = c.Conn.Read(cipherData)
	if n > 0 {
		c.Decrypt(b[0:n], cipherData[0:n])
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	var iv []byte
	if c.Enc == nil {
		iv, err = c.InitEncrypt()
		if err != nil {
			return
		}
	}

	cipherData := c.writeBuf
	dataSize := len(b) + len(iv)
	if dataSize > len(cipherData) {
		cipherData = make([]byte, dataSize)
	} else {
		cipherData = cipherData[:dataSize]
	}

	if iv != nil {
		copy(cipherData, iv)
	}

	c.Encrypt(cipherData[len(iv):], b)
	n, err = c.Conn.Write(cipherData)
	return
}
