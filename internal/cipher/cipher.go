package cipher

import (
	"crypto/aes"
	C "crypto/cipher"
)

type cipher struct {
	C.Block
}

func (c *cipher) Encrypt(originData []byte) (cipherData []byte, err error) {
	length := (len(originData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, originData)
	pad := byte(len(plain) - len(originData))
	for i := len(originData); i < len(plain); i++ {
		plain[i] = pad
	}
	cipherData = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, c.BlockSize(); bs <= len(originData); bs, be = bs+c.BlockSize(), be+c.BlockSize() {
		c.Block.Encrypt(cipherData[bs:be], plain[bs:be])
	}

	return
}

func (c *cipher) Decrypt(cipherData []byte) (originData []byte, err error) {

	originData = make([]byte, len(cipherData))

	for bs, be := 0, c.BlockSize(); bs < len(cipherData); bs, be = bs+c.BlockSize(), be+c.BlockSize() {
		c.Block.Decrypt(originData[bs:be], cipherData[bs:be])
	}

	trim := 0
	if len(originData) > 0 {
		trim = len(originData) - int(originData[len(originData)-1])
	}

	return originData[:trim], nil
}

func NewCipher(key []byte) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &cipher{block}, nil
}
