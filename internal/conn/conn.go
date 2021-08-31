package conn

import (
	YPCipher "github.com/kainhuck/yao-proxy/internal/cipher"
	"net"
)

type Conn struct {
	net.Conn
	YPCipher.Cipher
}

func NewConn(conn net.Conn, cipher YPCipher.Cipher) *Conn {
	return &Conn{
		Conn:   conn,
		Cipher: cipher,
	}
}

func DialAndSend(addr string, cipher YPCipher.Cipher, data []byte) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		Conn:   conn,
		Cipher: cipher,
	}

	_, err = c.Write(data)

	return c, err
}

func (c *Conn) Read(b []byte) (n int, err error) {
	cipherData := make([]byte, len(b))
	n, err = c.Conn.Read(cipherData)
	if n > 0 {
		rawData := c.Decrypt(cipherData[:n])
		if len(rawData) > len(b) {
			b = make([]byte, len(rawData))
		}

		copy(b, rawData)
		n = len(rawData)
	}

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	cipherData := c.Encrypt(b)

	return c.Conn.Write(cipherData)
}
