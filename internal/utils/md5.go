package utils

import "crypto/md5"

func MD5(bts []byte) []byte {
	h := md5.New()
	h.Write(bts)
	return h.Sum(nil)
}
