package aesgcm

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	keySize     = 32
	keyFileName = "aes-256-gcm-secret.key"
)

var keyMu sync.Mutex

func defaultKeyPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("aesgcm: 获取用户配置目录失败: %w", err)
	}
	return filepath.Join(configDir, "FortuneFinger", keyFileName), nil
}

func loadOrCreateKey(path string) ([]byte, error) {
	keyMu.Lock()
	defer keyMu.Unlock()

	key, err := loadKey(path)
	if err == nil {
		return key, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	key = make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("aesgcm: 生成加密密钥失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("aesgcm: 创建密钥目录失败: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return loadKey(path)
	}
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 创建密钥文件失败: %w", err)
	}
	if _, err := file.WriteString(encoded); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return nil, fmt.Errorf("aesgcm: 写入密钥文件失败: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return nil, fmt.Errorf("aesgcm: 关闭密钥文件失败: %w", err)
	}
	return key, nil
}

func loadKey(path string) ([]byte, error) {
	encoded, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("aesgcm: 加密密钥文件不存在: %w", err)
		}
		return nil, fmt.Errorf("aesgcm: 读取加密密钥失败: %w", err)
	}
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(encoded)))
	if err != nil {
		return nil, fmt.Errorf("aesgcm: 解码加密密钥失败: %w", err)
	}
	if len(key) != keySize {
		return nil, errors.New("aesgcm: 加密密钥长度必须为32字节")
	}
	return key, nil
}
