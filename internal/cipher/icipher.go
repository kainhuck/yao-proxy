package cipher

// Cipher 加解密
type Cipher interface {
	Encrypt(originData []byte) (cipherData []byte, err error)
	Decrypt(cipherData []byte) (originData []byte, err error)
}

func NewCipher(key []byte) (Cipher, error) {
	return NewAes(key)
}
