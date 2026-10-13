package runtime

import (
	"context"
	"errors"
	"fmt"

	"fortunefinger/backend/internal/config"
	"fortunefinger/backend/internal/pkg/gormc"
	"fortunefinger/backend/internal/pkg/logger"
	"fortunefinger/backend/internal/pkg/redisc"
	"fortunefinger/backend/internal/startup/datasource"

	"go.uber.org/zap"
)

type Environment struct {
	Config      config.Config
	Logger      *zap.Logger
	dataSources *gormc.Registry
	redis       *redisc.Client
}

func New(ctx context.Context, configFile string) (*Environment, error) {
	conf, err := config.Load(configFile)
	if err != nil {
		return nil, err
	}

	log, err := logger.New(conf.Log)
	if err != nil {
		return nil, fmt.Errorf("trun: 初始化日志失败: %w", err)
	}
	env := &Environment{
		Config: conf,
		Logger: log,
	}
	dataSources, err := datasource.Init(ctx, conf.Database, log)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("trun: 初始化数据源失败: %w", err),
			env.Close(),
		)
	}
	env.dataSources = dataSources

	redisClient, err := redisc.Init(ctx, conf.Redis)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("trun: 初始化Redis失败: %w", err),
			env.Close(),
		)
	}
	env.redis = redisClient

	return env, nil
}

func (e *Environment) Close() error {
	if e == nil {
		return nil
	}

	var errs []error
	if e.redis != nil {
		if err := e.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("trun: 关闭Redis失败: %w", err))
		}
	}
	if e.dataSources != nil {
		if err := e.dataSources.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if e.Logger != nil {
		// 控制台输出执行 Sync 在部分平台会返回无意义的 EINVAL，故有意忽略。
		_ = e.Logger.Sync()
	}
	return errors.Join(errs...)
}
