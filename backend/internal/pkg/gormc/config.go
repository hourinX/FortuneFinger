package gormc

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	DriverMySQL     = "mysql"
	DriverOracle    = "oracle"
	DriverPostgres  = "pgsql"
	DriverSQLServer = "sqlserver"
)

type Config struct {
	Driver      string `yaml:"Driver" mapstructure:"Driver"`
	Host        string `yaml:"Host" mapstructure:"Host"`
	Port        int    `yaml:"Port" mapstructure:"Port"`
	User        string `yaml:"User" mapstructure:"User"`
	Password    string `yaml:"Password" mapstructure:"Password"`
	DBName      string `yaml:"DBName" mapstructure:"DBName"`
	ServiceName string `yaml:"ServiceName" mapstructure:"ServiceName"`
	// SSLMode 会按驱动分别映射为 MySQL TLSConfig、PostgreSQL sslmode 或 SQL Server encrypt。
	SSLMode         string        `yaml:"SSLMode" mapstructure:"SSLMode"`
	MaxOpenConns    int           `yaml:"MaxOpenConns" mapstructure:"MaxOpenConns"`
	MaxIdleConns    int           `yaml:"MaxIdleConns" mapstructure:"MaxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime" mapstructure:"ConnMaxLifetime"`
	ConnMaxIdleTime time.Duration `yaml:"ConnMaxIdleTime" mapstructure:"ConnMaxIdleTime"`
	SlowThreshold   time.Duration `yaml:"SlowThreshold" mapstructure:"SlowThreshold"`
}

func (c Config) withDefaults() Config {
	if c.Port <= 0 {
		switch c.Driver {
		case DriverMySQL:
			c.Port = 3306
		case DriverOracle:
			c.Port = 1521
		case DriverPostgres:
			c.Port = 5432
		case DriverSQLServer:
			c.Port = 1433
		}
	}
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = 50
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = 20
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		c.MaxIdleConns = c.MaxOpenConns
	}
	if c.ConnMaxLifetime <= 0 {
		c.ConnMaxLifetime = 4 * time.Hour
	}
	if c.ConnMaxIdleTime <= 0 {
		c.ConnMaxIdleTime = 10 * time.Minute
	}
	if c.SlowThreshold <= 0 {
		c.SlowThreshold = 500 * time.Millisecond
	}
	return c
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Driver) == "" {
		return errors.New("gormc: 数据库驱动不能为空")
	}
	switch c.Driver {
	case DriverMySQL, DriverOracle, DriverPostgres, DriverSQLServer:
	default:
		return fmt.Errorf("gormc: 不支持的数据库驱动「%s」", c.Driver)
	}
	if strings.TrimSpace(c.Host) == "" {
		return errors.New("gormc: 数据库主机不能为空")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("gormc: 数据库端口必须在1到65535之间")
	}
	if strings.TrimSpace(c.User) == "" {
		return errors.New("gormc: 数据库用户不能为空")
	}
	if c.Driver == DriverOracle {
		if strings.TrimSpace(c.ServiceName) == "" {
			return errors.New("gormc: Oracle服务名不能为空")
		}
	} else if strings.TrimSpace(c.DBName) == "" {
		return errors.New("gormc: 数据库名不能为空")
	}
	return nil
}
