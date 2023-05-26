package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

type Password struct {
	Text 	string
	Key 	string
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
"!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`"

func GenSecureKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return []byte{}, err
	}

	return key, nil
}

func Encrypt(key []byte, text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherSlice := make([]byte, aes.BlockSize+len(text))
	iv := cipherSlice[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherSlice[aes.BlockSize:], text)

	return hex.EncodeToString(cipherSlice), nil
}

func Decrypt(key []byte, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(text) < aes.BlockSize {
		return []byte{}, err
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(text, text)

	return text, nil
}

func GenSecureString() (string, error) {
	charBytes := make([]byte, 64)
	_, err := rand.Read(charBytes)
	if err != nil {
		return "", err
	}

	for i, b := range charBytes {
		random := int(b) % len(charset)
		charBytes[i] = charset[random]
	}

	return string(charBytes), nil
}

func GenPassword () (*Password, error) {
	key, err := GenSecureKey()
	if err != nil {
		return &Password{}, err
	}

	password, err := GenSecureString()
	if err != nil {
		return &Password{}, err
	}

	cipherText, err := Encrypt(key, []byte(password))
	if err != nil {
		return &Password{}, err
	}

	return &Password{Text: cipherText, Key: hex.EncodeToString(key)}, nil
}