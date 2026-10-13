package captcha

import (
	"errors"
	"strings"
	"time"

	base64captcha "github.com/mojocn/base64Captcha"
)

const (
	imageBase64Prefix = "data:image/png;base64,"
	imageHeight       = 80
	imageWidth        = 240
)

type Generator struct {
	driver base64captcha.Driver
	ttl    time.Duration
}

func New(conf Config) (*Generator, error) {
	conf = conf.withDefaults()
	if err := conf.validate(); err != nil {
		return nil, err
	}

	return &Generator{
		driver: newDriver(conf),
		ttl:    conf.TTL,
	}, nil
}

func (g *Generator) Generate() (id, image, answer string, err error) {
	id, question, answer := g.driver.GenerateIdQuestionAnswer()
	item, err := g.driver.DrawCaptcha(question)
	if err != nil {
		return "", "", "", err
	}

	encoded := item.EncodeB64string()
	if !strings.HasPrefix(encoded, imageBase64Prefix) {
		return "", "", "", errors.New("captcha: 验证码图片编码失败")
	}
	image = strings.TrimPrefix(encoded, imageBase64Prefix)
	return id, image, answer, nil
}

func (g *Generator) TTL() time.Duration {
	return g.ttl
}

func newDriver(conf Config) base64captcha.Driver {
	switch conf.Mode {
	case ModeDigit:
		return base64captcha.NewDriverDigit(imageHeight, imageWidth, conf.Length, 0.7, 50)
	case ModeLetter:
		return newCharDriver(imageHeight, imageWidth, conf.Length, base64captcha.TxtAlphabet, false)
	case ModeMix:
		return newCharDriver(imageHeight, imageWidth, conf.Length, base64captcha.TxtSimpleCharaters, false)
	default:
		return newCharDriver(imageHeight, imageWidth, 0, "", true)
	}
}
