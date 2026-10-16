package config

import (
	"time"

	"fortunefinger/backend/internal/pkg/captcha"
	"fortunefinger/backend/internal/pkg/gormc"
	"fortunefinger/backend/internal/pkg/logger"
	"fortunefinger/backend/internal/pkg/redisc"
	"fortunefinger/backend/internal/pkg/storage"
)

type Config struct {
	Service              Service              `yaml:"Service" mapstructure:"Service"`
	Database             Database             `yaml:"Database" mapstructure:"Database"`
	Redis                redisc.Config        `yaml:"Redis" mapstructure:"Redis"`
	JWT                  JWT                  `yaml:"JWT" mapstructure:"JWT"`
	Captcha              captcha.Config       `yaml:"Captcha" mapstructure:"Captcha"`
	EnhancedVerification EnhancedVerification `yaml:"EnhancedVerification" mapstructure:"EnhancedVerification"`
	Login                Login                `yaml:"Login" mapstructure:"Login"`
	Snowflake            Snowflake            `yaml:"Snowflake" mapstructure:"Snowflake"`
	Storage              storage.Config       `yaml:"Storage" mapstructure:"Storage"`
	Log                  logger.Config        `yaml:"Log" mapstructure:"Log"`
}

type Service struct {
	Name string `yaml:"Name" mapstructure:"Name"`
	Host string `yaml:"Host" mapstructure:"Host"`
	Port int    `yaml:"Port" mapstructure:"Port"`
}

type Snowflake struct {
	MachineID int64 `yaml:"MachineID" mapstructure:"MachineID"`
}

type JWT struct {
	Secret                string        `yaml:"Secret" mapstructure:"Secret"`
	Issuer                string        `yaml:"Issuer" mapstructure:"Issuer"`
	Audience              string        `yaml:"Audience" mapstructure:"Audience"`
	AccessExpire          time.Duration `yaml:"AccessExpire" mapstructure:"AccessExpire"`
	RefreshExpire         time.Duration `yaml:"RefreshExpire" mapstructure:"RefreshExpire"`
	RememberRefreshExpire time.Duration `yaml:"RememberRefreshExpire" mapstructure:"RememberRefreshExpire"`
}

type Login struct {
	MaxFailures       int           `yaml:"MaxFailures" mapstructure:"MaxFailures"`
	FailureResetAfter time.Duration `yaml:"FailureResetAfter" mapstructure:"FailureResetAfter"`
}

type EnhancedVerification struct {
	Mode string        `yaml:"Mode" mapstructure:"Mode"`
	TTL  time.Duration `yaml:"TTL" mapstructure:"TTL"`
}

type Database struct {
	Sources map[string]DatabaseSource `yaml:"Sources" mapstructure:"Sources"`
}

type DatabaseSource struct {
	gormc.Config `yaml:",inline" mapstructure:",squash"`
	Required     bool `yaml:"Required" mapstructure:"Required"`
}
