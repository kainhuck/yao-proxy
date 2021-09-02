package utils

import (
	"crypto/md5"
	"time"
)

// ExponentialBackoff 指数退避
func ExponentialBackoff(f1 func() error, f2 func()) (err error) {
	sleepTime := time.Second
	for i := 0; i < 10; i++ {
		if err = f1(); err != nil {
			sleepTime <<= i
			time.Sleep(sleepTime)
			f2()
		} else {
			return nil
		}
	}
	return err
}

func MD5(bts []byte) []byte {
	h := md5.New()
	h.Write(bts)
	return h.Sum(nil)
}
