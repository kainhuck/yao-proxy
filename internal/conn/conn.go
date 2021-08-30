package conn

import (
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"net"
)

type Conn struct {
	net.Conn
	YPCipher.Cipher
	readBuff  []byte
}

func NewConn(conn net.Conn, cipher YPCipher.Cipher) *Conn {
	return &Conn{
		Conn:      conn,
		Cipher:    cipher,
		readBuff:  bufferPoolGet(),
	}
}

func (c *Conn) Close() {
	bufferPoolPut(c.readBuff)
	_ = c.Conn.Close()
}

func (c *Conn) Read(b []byte) (int, error) {
	cipherData := c.readBuff
	if len(b) > len(cipherData) {
		cipherData = make([]byte, len(b))
	} else {
		cipherData = cipherData[:len(b)]
	}

	n, err := c.Conn.Read(cipherData)
	if n > 0 {
		data, _ := c.Decrypt(cipherData[:n])
		copy(b, data)
	}

	return len(b), err
}

func (c *Conn) Write(b []byte) (int, error) {
	cipherData, _ := c.Encrypt(b)
	return c.Conn.Write(cipherData)
}