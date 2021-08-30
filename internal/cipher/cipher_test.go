package cipher

import "testing"

func TestCipher(t *testing.T) {
	ci, err := NewAes([]byte("1234567890123456"))
	if err != nil {
		t.Fatal(err)
	}

	rawData := []byte("kainhuck adfadadf")
	cipherData := make([]byte, 0)

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

	cipherData, err = ci.Encrypt(rawData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cipherData), len(cipherData))
	rawData, err = ci.Decrypt(cipherData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawData), len(rawData))

}
