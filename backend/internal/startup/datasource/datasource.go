package datasource

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"fortunefinger/backend/internal/config"
	"fortunefinger/backend/internal/pkg/gormc"

	"go.uber.org/zap"
)

func Init(ctx context.Context, conf config.Database, log *zap.Logger) (*gormc.Registry, error) {
	if log == nil {
		return nil, errors.New("datasource: logger不能为空")
	}
	if len(conf.Sources) == 0 {
		return nil, errors.New("datasource: 数据源配置不能为空")
	}
	mainSource, exists := conf.Sources[gormc.MainSource]
	if !exists {
		return nil, fmt.Errorf("datasource: 主数据源「%s」不存在", gormc.MainSource)
	}
	if !mainSource.Required {
		return nil, fmt.Errorf("datasource: 主数据源「%s」必须设置Required为true", gormc.MainSource)
	}

	registry := gormc.NewRegistry()
	names := make([]string, 0, len(conf.Sources))
	for name := range conf.Sources {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		sourceConf := conf.Sources[name]
		sourceLog := log.With(
			zap.String("database", name),
			zap.String("driver", sourceConf.Driver),
		)

		source, err := gormc.Init(ctx, name, sourceConf.Config, sourceLog)
		if err != nil {
			if name == gormc.MainSource || sourceConf.Required {
				closeErr := registry.Close()
				return nil, joinDataSourceInitError(name, err, closeErr)
			}
			sourceLog.Warn("可选数据源初始化失败", zap.Error(err))
			continue
		}

		if err := registry.Add(name, source); err != nil {
			_ = source.Close()
			closeErr := registry.Close()
			return nil, joinDataSourceInitError(name, err, closeErr)
		}
	}

	return registry, nil
}

func joinDataSourceInitError(name string, initErr, closeErr error) error {
	if closeErr != nil {
		return errors.Join(
			fmt.Errorf("datasource: 初始化必要数据源「%s」失败: %w", name, initErr),
			fmt.Errorf("datasource: 关闭已初始化数据源失败: %w", closeErr),
		)
	}
	return fmt.Errorf("datasource: 初始化必要数据源「%s」失败: %w", name, initErr)
}
