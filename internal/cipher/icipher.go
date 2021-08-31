package cipher

// Cipher 加解密
type Cipher interface {
	Encrypt(originData []byte) (cipherData []byte)
	Decrypt(cipherData []byte) (originData []byte)
}

func NewCipher(key []byte) (Cipher, error) {
	return NewAes(key)
}
