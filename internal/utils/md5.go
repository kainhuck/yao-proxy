package utils

import "crypto/md5"

func MD5(s string) []byte {
	b := md5.Sum([]byte(s))

	return b[:]
}
