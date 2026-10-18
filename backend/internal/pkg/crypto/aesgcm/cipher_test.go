package aesgcm

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCipherRoundTrip(t *testing.T) {
	cipher := &Cipher{keyPath: filepath.Join(t.TempDir(), keyFileName)}
	plaintext := []byte("smtp-secret")

	first, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	second, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() second error = %v", err)
	}
	if first == second {
		t.Fatal("Encrypt() returned identical ciphertext for the same plaintext")
	}
	if !strings.HasPrefix(first, ciphertextVersion+":") {
		t.Fatalf("Encrypt() value = %q, missing version prefix", first)
	}

	decrypted, err := cipher.Decrypt(first)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
	key, err := os.ReadFile(cipher.keyPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if strings.TrimSpace(string(key)) == "" {
		t.Fatal("generated key file is empty")
	}
}

func TestDecryptRejectsTamperedCiphertext(t *testing.T) {
	cipher := &Cipher{keyPath: filepath.Join(t.TempDir(), keyFileName)}
	value, err := cipher.Encrypt([]byte("smtp-secret"))
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	last := value[len(value)-1]
	replacement := byte('A')
	if last == replacement {
		replacement = 'B'
	}
	value = value[:len(value)-1] + string(replacement)

	if _, err := cipher.Decrypt(value); err == nil {
		t.Fatal("Decrypt() accepted tampered ciphertext")
	}
}

func TestDecryptDoesNotCreateMissingKey(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), keyFileName)
	cipher := &Cipher{keyPath: keyPath}

	if _, err := cipher.Decrypt("v1:YWJj"); err == nil {
		t.Fatal("Decrypt() succeeded without a key file")
	}
	if _, err := os.Stat(keyPath); !os.IsNotExist(err) {
		t.Fatalf("Decrypt() created a missing key file, Stat() error = %v", err)
	}
}
