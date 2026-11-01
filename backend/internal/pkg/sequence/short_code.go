package sequence

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	shortCodeLength      = 8
	maxShortCodeAttempts = 16
	uppercase            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase            = "abcdefghijklmnopqrstuvwxyz"
	digits               = "0123456789"
	alphanumeric         = uppercase + lowercase + digits
)

type ShortCodeExistsFunc func(context.Context, string) (bool, error)

func GenerateUniqueShortCode(ctx context.Context, exists ShortCodeExistsFunc) (string, error) {
	if exists == nil {
		return "", errors.New("sequence: 编码唯一性检查函数不能为空")
	}
	for range maxShortCodeAttempts {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		code, err := generateShortCode()
		if err != nil {
			return "", err
		}
		found, err := exists(ctx, code)
		if err != nil {
			return "", err
		}
		if !found {
			return code, nil
		}
	}
	return "", errors.New("sequence: 无法生成唯一编码")
}

func generateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)
	groups := []string{uppercase, lowercase, digits}
	for index, group := range groups {
		value, err := randomCharacter(group)
		if err != nil {
			return "", err
		}
		code[index] = value
	}
	for index := len(groups); index < len(code); index++ {
		value, err := randomCharacter(alphanumeric)
		if err != nil {
			return "", err
		}
		code[index] = value
	}
	for index := len(code) - 1; index > 0; index-- {
		position, err := rand.Int(rand.Reader, big.NewInt(int64(index+1)))
		if err != nil {
			return "", err
		}
		other := int(position.Int64())
		code[index], code[other] = code[other], code[index]
	}
	return string(code), nil
}

func randomCharacter(characters string) (byte, error) {
	position, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
	if err != nil {
		return 0, err
	}
	return characters[position.Int64()], nil
}
