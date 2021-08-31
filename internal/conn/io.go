package conn

import (
	"net"
	"sync"
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

// Copy src -> dst
func Copy(dst, src net.Conn) error {
	buff := bufferPoolGet()
	defer bufferPoolPut(buff)

	for {
		n, err := src.Read(buff)
		if n > 0 {
			if _, err := dst.Write(buff[:n]); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
}
