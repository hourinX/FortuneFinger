package storage

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DriverMinIO = "minio"
	DriverOSS   = "oss"
)

type Config struct {
	Driver    string `yaml:"Driver" mapstructure:"Driver"`
	Endpoint  string `yaml:"Endpoint" mapstructure:"Endpoint"`
	Region    string `yaml:"Region" mapstructure:"Region"`
	AccessKey string `yaml:"AccessKey" mapstructure:"AccessKey"`
	SecretKey string `yaml:"SecretKey" mapstructure:"SecretKey"`
	Bucket    string `yaml:"Bucket" mapstructure:"Bucket"`
	Prefix    string `yaml:"Prefix" mapstructure:"Prefix"`
	Secure    bool   `yaml:"Secure" mapstructure:"Secure"`
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Driver) == "" {
		return errors.New("storage: Storage.Driver不能为空")
	}
	if strings.TrimSpace(c.Endpoint) == "" {
		return errors.New("storage: Storage.Endpoint不能为空")
	}
	if strings.TrimSpace(c.AccessKey) == "" {
		return errors.New("storage: Storage.AccessKey不能为空")
	}
	if strings.TrimSpace(c.SecretKey) == "" {
		return errors.New("storage: Storage.SecretKey不能为空")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return errors.New("storage: Storage.Bucket不能为空")
	}
	if !strings.Contains(c.Prefix, "{file}") {
		return errors.New("storage: Storage.Prefix必须包含{file}变量")
	}
	remaining := strings.NewReplacer(
		"{scene}", "",
		"{yyyyMMdd}", "",
		"{yyyyMM}", "",
		"{yyyy-MM-dd}", "",
		"{yyyy-MM}", "",
		"{yyyy}", "",
		"{MM}", "",
		"{dd}", "",
		"{file}", "",
	).Replace(c.Prefix)
	if strings.ContainsAny(remaining, "{}") {
		return fmt.Errorf("storage: Storage.Prefix包含不支持的变量「%s」", c.Prefix)
	}
	if c.Driver == DriverOSS && strings.TrimSpace(c.Region) == "" {
		return errors.New("storage: OSS驱动的Storage.Region不能为空")
	}
	if c.Driver != DriverMinIO && c.Driver != DriverOSS {
		return fmt.Errorf("storage: 不支持的存储驱动「%s」", c.Driver)
	}
	return nil
}
