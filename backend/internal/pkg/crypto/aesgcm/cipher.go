package aesgcm

import (
	"crypto/aes"
	stdcipher "crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const ciphertextVersion = "v1"

type Cipher struct {
	keyPath string
}

func New() (*Cipher, error) {
	keyPath, err := defaultKeyPath()
	if err != nil {
		return nil, err
	}
	return &Cipher{keyPath: keyPath}, nil
}

func (c *Cipher) Encrypt(plaintext []byte) (string, error) {
	key, err := loadOrCreateKey(c.keyPath)
	if err != nil {
		return "", err
	}

	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("aesgcm: 生成Nonce失败: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, []byte(ciphertextVersion))
	payload := append(nonce, ciphertext...)
	return ciphertextVersion + ":" + base64.RawStdEncoding.EncodeToString(payload), nil
}

func (c *Cipher) Decrypt(value string) ([]byte, error) {
	version, encoded, found := strings.Cut(value, ":")
	if !found || version != ciphertextVersion || encoded == "" {
		return nil, errors.New("aesgcm: 密文格式无效")
	}

	payload, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 解码密文失败: %w", err)
	}
	key, err := loadKey(c.keyPath)
	if err != nil {
		return nil, err
	}
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	if len(payload) < gcm.NonceSize()+gcm.Overhead() {
		return nil, errors.New("aesgcm: 密文长度无效")
	}

	nonce := payload[:gcm.NonceSize()]
	ciphertext := payload[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, []byte(version))
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 解密失败: %w", err)
	}
	return plaintext, nil
}

func (c *Cipher) KeyPath() string {
	return c.keyPath
}

func newGCM(key []byte) (stdcipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 创建AES密码块失败: %w", err)
	}
	gcm, err := stdcipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 创建GCM失败: %w", err)
	}
	return gcm, nil
}
