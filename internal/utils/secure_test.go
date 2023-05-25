package secure

import (
	"testing"
)

func TestGenSecureKey(t *testing.T) {
	_, err := GenSecureKey()

	if err != nil {
		t.Error("Error generating secure key")
	}
}

func TestEncrypt(t *testing.T) {
	key, err := GenSecureKey()
	if err != nil {
		t.Error("Error generating secure key")
	}

	_, err = Encrypt(key, []byte("test"))
	if err != nil {
		t.Error("Error encrypting text")
	}
}

func TestDecrypt(t *testing.T) {
	key, err := GenSecureKey()
	if err != nil {
		t.Error("Error generating secure key")
	}

	cipherText, err := Encrypt(key, []byte("test"))
	if err != nil {
		t.Error("Error encrypting text")
	}

	_, err = Decrypt(key, []byte(cipherText))
	if err != nil {
		t.Error("Error decrypting text")
	}
}

func TestGenSecureString(t *testing.T) {
	_, err := GenSecureString()
	if err != nil {
		t.Error("Error generating secure password")
	}
}

func TestGenPassword(t *testing.T) {
	_, err := GenPassword()
	if err != nil {
		t.Error("Error generating secure password")
	}
}