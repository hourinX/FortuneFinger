package gormc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Source struct {
	db   *gorm.DB
	pool *sql.DB
}

// 初始化数据库
func Init(ctx context.Context, sourceName string, conf Config, log *zap.Logger) (*Source, error) {
	if strings.TrimSpace(sourceName) == "" {
		return nil, errors.New("gormc: Database.Sources的键名不能为空")
	}
	if log == nil {
		return nil, fmt.Errorf("gormc: 数据源「%s」logger不能为空", sourceName)
	}
	conf = conf.withDefaults()
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("gormc: 数据源「%s」配置无效: %w", sourceName, err)
	}

	pool, dialector, err := openConnection(ctx, sourceName, conf)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(
		dialector,
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			DisableAutomaticPing:                     true,
			Logger:                                   newGormLogger(log, conf.SlowThreshold),
		},
	)
	if err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf(
			"gormc: 打开数据源「%s」的「%s」数据库失败: %w",
			sourceName,
			conf.Driver,
			err,
		)
	}

	return &Source{
		db:   db,
		pool: pool,
	}, nil
}

func (s *Source) DB() *gorm.DB {
	return s.db
}

func (s *Source) Close() error {
	return s.pool.Close()
}
