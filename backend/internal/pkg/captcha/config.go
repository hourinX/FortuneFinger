package captcha

import (
	"fmt"
	"time"
)

type Mode string

const (
	ModeDigit  Mode = "digit"
	ModeMath   Mode = "math"
	ModeLetter Mode = "letter"
	ModeMix    Mode = "mix"

	defaultLength = 4
	defaultTTL    = 5 * time.Minute
)

type Config struct {
	Mode   Mode          `yaml:"Mode" mapstructure:"Mode"`
	Length int           `yaml:"Length" mapstructure:"Length"`
	TTL    time.Duration `yaml:"TTL" mapstructure:"TTL"`
}

func (c Config) withDefaults() Config {
	if c.Mode == "" {
		c.Mode = ModeMath
	}
	if c.Length == 0 {
		c.Length = defaultLength
	}
	if c.TTL == 0 {
		c.TTL = defaultTTL
	}
	return c
}

func (c Config) validate() error {
	switch c.Mode {
	case ModeDigit, ModeMath, ModeLetter, ModeMix:
	default:
		return fmt.Errorf("captcha: 不支持的验证码模式「%s」", c.Mode)
	}
	if c.Length < 1 || c.Length > 10 {
		return fmt.Errorf("captcha: 验证码长度必须在1到10之间，当前值为「%d」", c.Length)
	}
	if c.TTL <= 0 {
		return fmt.Errorf("captcha: 验证码有效期必须大于0，当前值为「%s」", c.TTL)
	}
	return nil
}
