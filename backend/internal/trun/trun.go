// Package trun 仅服务于测试，为测试执行提供 Redis 等中间件及 GORM 运行环境。
package trun

import (
	"context"

	truntime "fortunefinger/backend/internal/trun/runtime"
)

type Environment = truntime.Environment

func New(ctx context.Context, configFile string) (*Environment, error) {
	return truntime.New(ctx, configFile)
}

func Run(ctx context.Context, configFile string, fn func(*Environment) error) error {
	return truntime.Run(ctx, configFile, fn)
}
