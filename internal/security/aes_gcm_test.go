package security

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecryptAESGCM(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	plaintext := []byte("secret data")

	ciphertext, err := EncryptAESGCM(plaintext, key)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := DecryptAESGCM(key, ciphertext)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted data does not match original.\nGot:  %s\nWant: %s", decrypted, plaintext)
	}
}
