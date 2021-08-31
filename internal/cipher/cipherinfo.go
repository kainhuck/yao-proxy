package cipher

import "crypto/cipher"

type DecOrEnc int

const (
	Decrypt DecOrEnc = iota
	Encrypt
)

type cipherInfo struct {
	KeyLen    int
	IvLen     int
	NewStream func(key, iv []byte, doe DecOrEnc) (cipher.Stream, error)
}

var cipherMethod = map[string]*cipherInfo{
	"aes-128-cfb": {16, 16, newAESCFBStream},
	"aes-192-cfb": {24, 16, newAESCFBStream},
	"aes-256-cfb": {32, 16, newAESCFBStream},
	"aes-128-ctr": {16, 16, newAESCTRStream},
	"aes-192-ctr": {24, 16, newAESCTRStream},
	"aes-256-ctr": {32, 16, newAESCTRStream},
	"des-cfb":     {8, 8, newDESStream},
	"rc4-md5":     {16, 16, newRC4MD5Stream},
	"rc4-md5-6":   {16, 6, newRC4MD5Stream},
}
