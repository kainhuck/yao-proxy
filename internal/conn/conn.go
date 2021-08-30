package conn

import (
	"fmt"
	"github.com/kainhuck/yao-proxy/internal/cipher"
	"net"
	"sync"
	"time"
)

const bufSize = 65535

var buffPool sync.Pool

func init() {
	buffPool.New = func() interface{} {
		return make([]byte, bufSize)
	}
}

func bufferPoolGet() []byte {
	return buffPool.Get().([]byte)
}

func bufferPoolPut(b []byte) {
	buffPool.Put(b)
}

func DecryptRead(conn net.Conn, timeout time.Duration, ci cipher.Cipher) ([]byte, error) {
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}
	buff := bufferPoolGet()
	defer bufferPoolPut(buff)
	receiveData := make([]byte, 0)
	receiveSize := 0
	if timeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return nil, err
		}
	}
	for {
		n, err := conn.Read(buff)
		if err != nil {
			return nil, err
		}

		if n > 0 {
			receiveSize += n
			receiveData = append(receiveData, buff[:n]...)
			if n < bufSize {
				break
			}
		}
	}
	response := make([]byte, receiveSize)
	copy(response, receiveData)

	return ci.Decrypt(response)
}

func EncryptWrite(conn net.Conn, ci cipher.Cipher, data []byte) error {
	cipherData, err := ci.Encrypt(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(cipherData)

	return err
}

// EncryptCopy 从src 读出数据 加密后发给dst
func EncryptCopy(dst net.Conn, src net.Conn, ci cipher.Cipher) error {
	buff := bufferPoolGet()
	defer bufferPoolPut(buff)
	for {
		n, err := src.Read(buff)
		if n > 0 {
			cipherData, err := ci.Encrypt(buff[:n])
			if err != nil {
				return err
			}
			_, err = dst.Write(cipherData)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}

// DecryptCopy 从src 读出数据 解密后发给dst
func DecryptCopy(dst net.Conn, src net.Conn, ci cipher.Cipher) error {
	buff := bufferPoolGet()
	defer bufferPoolPut(buff)
	for {
		n, err := src.Read(buff)
		if n > 0 {
			rawData, err := ci.Decrypt(buff[:n])
			if err != nil {
				return err
			}
			_, err = dst.Write(rawData)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
}
