package dependency

import (
	"context"
	"errors"
	"fmt"

	enhancedverification "fortunefinger/backend/internal/biz/enhanced-verification"
	loginsecurity "fortunefinger/backend/internal/biz/login-security"
	"fortunefinger/backend/internal/config"
	"fortunefinger/backend/internal/pkg/auth"
	"fortunefinger/backend/internal/pkg/captcha"
	"fortunefinger/backend/internal/pkg/gormc"
	"fortunefinger/backend/internal/pkg/logger"
	"fortunefinger/backend/internal/pkg/redisc"
	"fortunefinger/backend/internal/pkg/snowflake"
	"fortunefinger/backend/internal/pkg/storage"
	"fortunefinger/backend/internal/startup/datasource"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

func Init(ctx context.Context, conf config.Config) (*svc.ServiceContext, error) {
	log, err := logger.New(conf.Log)
	if err != nil {
		return nil, fmt.Errorf("dependency: 初始化日志失败: %w", err)
	}

	dataSources, err := datasource.Init(ctx, conf.Database, log)
	if err != nil {
		return nil, initializationError(log, nil, nil, "初始化数据源失败", err)
	}

	redisClient, err := redisc.Init(ctx, conf.Redis)
	if err != nil {
		return nil, initializationError(log, dataSources, nil, "初始化Redis失败", err)
	}

	storageClient, err := storage.Init(ctx, conf.Storage)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化对象存储失败", err)
	}

	captchaGenerator, err := captcha.New(conf.Captcha)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化验证码生成器失败", err)
	}

	snowflakeGenerator, err := snowflake.New(conf.Snowflake.MachineID)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化雪花ID生成器失败", err)
	}

	jwtManager, err := auth.NewManager(
		conf.JWT.Secret,
		conf.JWT.Issuer,
		conf.JWT.Audience,
		conf.JWT.AccessExpire,
		conf.JWT.RefreshExpire,
		conf.JWT.RememberRefreshExpire,
	)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化JWT管理器失败", err)
	}

	loginSecurity, err := loginsecurity.NewService(redisClient, conf.Login.MaxFailures, conf.Login.FailureResetAfter)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化登录风控失败", err)
	}

	enhancedVerification, err := enhancedverification.NewService(
		redisClient,
		conf.EnhancedVerification.Mode,
		conf.EnhancedVerification.TTL,
	)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "初始化增强验证服务失败", err)
	}

	svcCtx, err := svc.NewServiceContext(
		conf,
		log,
		dataSources,
		redisClient,
		storageClient,
		captchaGenerator,
		snowflakeGenerator,
		jwtManager,
		loginSecurity,
		enhancedVerification,
	)
	if err != nil {
		return nil, initializationError(log, dataSources, redisClient, "装配ServiceContext失败", err)
	}
	return svcCtx, nil
}

func initializationError(
	log *zap.Logger,
	dataSources *gormc.Registry,
	redisClient *redisc.Client,
	message string,
	cause error,
) error {
	initErr := fmt.Errorf("dependency: %s: %w", message, cause)
	log.Error(message, zap.Error(initErr))
	closeErr := closeInfrastructure(dataSources, redisClient)
	_ = log.Sync()
	return errors.Join(initErr, closeErr)
}

func Close(svcCtx *svc.ServiceContext) error {
	if svcCtx == nil {
		return nil
	}
	return closeInfrastructure(svcCtx.DataSources, svcCtx.Redis)
}

func closeInfrastructure(dataSources *gormc.Registry, redisClient *redisc.Client) error {
	var errs []error
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("dependency: 关闭Redis失败: %w", err))
		}
	}
	if dataSources != nil {
		if err := dataSources.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
