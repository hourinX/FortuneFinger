package redisc

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Config struct {
	Host     string `yaml:"Host" mapstructure:"Host"`
	Password string `yaml:"Password" mapstructure:"Password"`
	DB       int    `yaml:"DB" mapstructure:"DB"`
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return errors.New("redisc: Redis地址不能为空")
	}
	if _, _, err := net.SplitHostPort(c.Host); err != nil {
		return fmt.Errorf("redisc: Redis地址「%s」格式错误，应形如localhost:6379", c.Host)
	}
	if c.DB < 0 {
		return fmt.Errorf("redisc: Redis库号不能为负数「%d」", c.DB)
	}
	return nil
}
