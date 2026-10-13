package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

func Load(path string) (Config, error) {
	loader := viper.New()
	loader.SetConfigFile(path)
	if err := loader.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("config: 读取配置文件「%s」失败: %w", path, err)
	}

	var conf Config
	if err := loader.UnmarshalExact(
		&conf,
		viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc()),
	); err != nil {
		return Config{}, fmt.Errorf("config: 解析配置文件「%s」失败: %w", path, err)
	}
	return conf, nil
}
