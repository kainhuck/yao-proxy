package conn

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const bufSize = 256

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

func Read(conn net.Conn, timeout time.Duration) ([]byte, error) {
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

	return response, nil
}
