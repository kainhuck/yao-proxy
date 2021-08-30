package cipher

type Empty struct {
}

func (e Empty) Encrypt(originData []byte) (cipherData []byte, err error) {
	return originData, nil
}

func (e Empty) Decrypt(cipherData []byte) (originData []byte, err error) {
	return cipherData, nil
}

func NewEmpty() (Cipher, error) {
	return &Empty{}, nil
}
