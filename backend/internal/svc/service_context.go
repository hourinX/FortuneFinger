package svc

import (
	"errors"

	enhancedverification "fortunefinger/backend/internal/biz/enhanced-verification"
	loginsecurity "fortunefinger/backend/internal/biz/login-security"
	"fortunefinger/backend/internal/config"
	"fortunefinger/backend/internal/pkg/auth"
	"fortunefinger/backend/internal/pkg/captcha"
	"fortunefinger/backend/internal/pkg/gormc"
	"fortunefinger/backend/internal/pkg/redisc"
	"fortunefinger/backend/internal/pkg/snowflake"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config               config.Config
	Logger               *zap.Logger
	GormDB               *gorm.DB
	DataSources          *gormc.Registry
	Redis                *redisc.Client
	Captcha              *captcha.Generator
	Snowflake            *snowflake.Generator
	JWT                  *auth.Manager
	LoginSecurity        *loginsecurity.Service
	EnhancedVerification *enhancedverification.Service
}

func NewServiceContext(
	conf config.Config,
	log *zap.Logger,
	dataSources *gormc.Registry,
	redisClient *redisc.Client,
	captchaGenerator *captcha.Generator,
	snowflakeGenerator *snowflake.Generator,
	jwtManager *auth.Manager,
	loginSecurity *loginsecurity.Service,
	enhancedVerification *enhancedverification.Service,
) (*ServiceContext, error) {
	if log == nil {
		return nil, errors.New("svc: logger不能为空")
	}
	if dataSources == nil {
		return nil, errors.New("svc: 数据源不能为空")
	}
	if redisClient == nil {
		return nil, errors.New("svc: Redis客户端不能为空")
	}
	if captchaGenerator == nil {
		return nil, errors.New("svc: 验证码生成器不能为空")
	}
	if snowflakeGenerator == nil {
		return nil, errors.New("svc: 雪花ID生成器不能为空")
	}
	if jwtManager == nil {
		return nil, errors.New("svc: JWT管理器不能为空")
	}
	if loginSecurity == nil {
		return nil, errors.New("svc: 登录安全服务不能为空")
	}
	if enhancedVerification == nil {
		return nil, errors.New("svc: 增强验证服务不能为空")
	}

	mainDB, err := dataSources.Main()
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:               conf,
		Logger:               log,
		GormDB:               mainDB,
		DataSources:          dataSources,
		Redis:                redisClient,
		Captcha:              captchaGenerator,
		Snowflake:            snowflakeGenerator,
		JWT:                  jwtManager,
		LoginSecurity:        loginSecurity,
		EnhancedVerification: enhancedVerification,
	}, nil
}
