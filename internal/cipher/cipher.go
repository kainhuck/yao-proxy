package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Cipher 加解密
type Cipher struct {
	Enc  cipher.Stream
	Dec  cipher.Stream
	Key  []byte
	Info *cipherInfo
}

func NewCipher(method, password string) (c *Cipher, err error) {
	if password == "" {
		return nil, fmt.Errorf("password is empty")
	}
	mi, ok := cipherMethod[method]
	if !ok {
		return nil, errors.New("Unsupported encryption method: " + method)
	}

	key := evpBytesToKey(password, mi.KeyLen)

	c = &Cipher{Key: key, Info: mi}

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cipher) InitEncrypt() (iv []byte, err error) {
	iv = make([]byte, c.Info.IvLen)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	c.Enc, err = c.Info.NewStream(c.Key, iv, Encrypt)
	return
}

func (c *Cipher) InitDecrypt(iv []byte) (err error) {
	c.Dec, err = c.Info.NewStream(c.Key, iv, Decrypt)
	return
}

func (c *Cipher) Encrypt(dst, src []byte) {
	c.Enc.XORKeyStream(dst, src)
}

func (c *Cipher) Decrypt(dst, src []byte) {
	c.Dec.XORKeyStream(dst, src)
}

func (c *Cipher) Copy() *Cipher {
	nc := *c
	nc.Enc = nil
	nc.Dec = nil
	return &nc
}
