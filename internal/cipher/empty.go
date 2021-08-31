package cipher

type Empty struct {
}

func (e Empty) Encrypt(originData []byte) (cipherData []byte) {
	return originData
}

func (e Empty) Decrypt(cipherData []byte) (originData []byte) {
	return cipherData
}

func NewEmpty() (Cipher, error) {
	return &Empty{}, nil
}
