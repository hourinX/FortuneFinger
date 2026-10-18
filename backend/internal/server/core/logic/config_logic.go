package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	mailsend "fortunefinger/backend/internal/biz/mail-send"
	"fortunefinger/backend/internal/biz/operator"
	sysconfig "fortunefinger/backend/internal/biz/sys-config"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/crypto/aesgcm"
	"fortunefinger/backend/internal/pkg/errorx"
	mailhtml "fortunefinger/backend/internal/pkg/html"
	"fortunefinger/backend/internal/pkg/mailer"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type ConfigLogic interface {
	GetConfig(ctx context.Context, req *types.GetConfigReq) (*types.GetConfigResp, error)
	SaveConfig(ctx context.Context, req *types.SaveConfigReq) (*types.SaveConfigResp, error)
	SendTestMail(ctx context.Context, req *types.SendTestMailReq) error
}

type configLogic struct {
	svcCtx        *svc.ServiceContext
	configService *sysconfig.Service
	mailService   *mailsend.Service
	userRepo      repository.UserRepository
}

func NewConfigLogic(svcCtx *svc.ServiceContext) ConfigLogic {
	return &configLogic{
		svcCtx:        svcCtx,
		configService: sysconfig.NewService(svcCtx),
		mailService:   mailsend.NewService(svcCtx),
		userRepo:      repository.NewUserRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (l *configLogic) GetConfig(ctx context.Context, req *types.GetConfigReq) (*types.GetConfigResp, error) {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New(errorx.Unauthorized, "非法的操作用户")
	}

	config, err := l.configService.GetConfigByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &types.GetConfigResp{
			Code:   req.Code,
			Config: json.RawMessage(`{}`),
		}, nil
	}

	configVal := json.RawMessage(config.Config)
	if !json.Valid(configVal) {
		err := errorx.New(errorx.InternalServer, "系统配置内容无效")
		l.svcCtx.Logger.Error("[ConfigLogic] 解析系统配置失败",
			zap.Uint64("config_id", config.ID),
			zap.String("config_code", config.Code),
			zap.Error(err),
		)
		return nil, err
	}

	resp := types.GetConfigResp{
		ID:     config.ID,
		Code:   config.Code,
		Config: configVal,
	}
	return &resp, nil
}

func (l *configLogic) SaveConfig(ctx context.Context, req *types.SaveConfigReq) (*types.SaveConfigResp, error) {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New(errorx.Unauthorized, "非法的操作用户")
	}

	debounceKey := constant.ConfigSaveDebounceKey + ":" + req.Code
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[ConfigLogic] Redis设置保存系统配置防抖失败",
			zap.String("config_code", req.Code),
			zap.Uint64("operator_id", operatorID),
			zap.Error(err),
		)
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		err := errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
		return nil, err
	}

	var config map[string]json.RawMessage
	if err := json.Unmarshal(req.Config, &config); err != nil || config == nil {
		return nil, errorx.New(errorx.InvalidParameter, "系统配置内容必须是JSON对象")
	}
	if req.Code == "mail" {
		delete(config, "testEmail")
		cipher, err := aesgcm.New()
		if err != nil {
			l.svcCtx.Logger.Error("[ConfigLogic] 初始化邮件认证信息加密器失败",
				zap.String("config_code", req.Code),
				zap.Uint64("operator_id", operatorID),
				zap.Error(err),
			)
			return nil, errorx.Wrap(errorx.InternalServer, "邮件认证信息加密失败", err)
		}
		for _, field := range []string{"username", "password"} {
			rawValue, exists := config[field]
			if !exists {
				continue
			}
			var value string
			if err := json.Unmarshal(rawValue, &value); err != nil {
				return nil, errorx.New(errorx.InvalidParameter, "邮件认证信息格式错误")
			}
			if value == "" {
				continue
			}
			if strings.HasPrefix(value, "v1:") {
				if _, err := cipher.Decrypt(value); err != nil {
					return nil, errorx.New(errorx.InvalidParameter, "邮件认证信息密文无效")
				}
				continue
			}
			encrypted, err := cipher.Encrypt([]byte(value))
			if err != nil {
				l.svcCtx.Logger.Error("[ConfigLogic] 加密邮件认证信息失败",
					zap.String("config_code", req.Code),
					zap.String("config_field", field),
					zap.Uint64("operator_id", operatorID),
					zap.Error(err),
				)
				return nil, errorx.Wrap(errorx.InternalServer, "邮件认证信息加密失败", err)
			}
			encoded, err := json.Marshal(encrypted)
			if err != nil {
				l.svcCtx.Logger.Error("[ConfigLogic] 序列化邮件认证信息密文失败",
					zap.String("config_code", req.Code),
					zap.String("config_field", field),
					zap.Uint64("operator_id", operatorID),
					zap.Error(err),
				)
				return nil, errorx.Wrap(errorx.InternalServer, "邮件认证信息加密失败", err)
			}
			config[field] = encoded
		}
	}
	configValue, err := json.Marshal(config)
	if err != nil {
		return nil, errorx.Wrap(errorx.InternalServer, "系统配置内容无效", err)
	}

	configID, err := l.configService.SaveConfig(ctx, sysconfig.SaveInput{
		ID:         req.ID,
		Code:       req.Code,
		Config:     string(configValue),
		OperatorID: operatorID,
	})
	if err == nil {
		return &types.SaveConfigResp{
			ID:   configID,
			Code: req.Code,
		}, nil
	}

	appErr, ok := errorx.As(err)
	if ok && (appErr.Code() == errorx.Conflict || appErr.Code() == errorx.NotFound) {
		return nil, err
	}
	l.svcCtx.Logger.Error("[ConfigLogic] 保存系统配置失败",
		zap.Uint64("config_id", req.ID),
		zap.String("config_code", req.Code),
		zap.Uint64("operator_id", operatorID),
		zap.Error(err),
	)
	return nil, err
}

func (l *configLogic) SendTestMail(ctx context.Context, req *types.SendTestMailReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	user, err := l.userRepo.FindByID(ctx, operatorID)
	if err != nil {
		return err
	}
	if user == nil {
		return errorx.New(errorx.Unauthorized, "非法的操作用户")
	}

	debounceKey := constant.ConfigTestMailDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.ConfigTestMailDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[ConfigLogic] Redis设置测试邮件防抖失败",
			zap.Uint64("operator_id", operatorID),
			zap.Error(err),
		)
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "测试邮件发送过于频繁，请稍后重试")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	return l.mailService.Send(ctx, mailsend.Input{
		To:       []string{req.Email},
		Subject:  "FortuneFinger 测试邮件",
		TextBody: "FortuneFinger 邮件服务配置测试成功，发送时间：" + now,
		HTMLBuilder: func(info mailsend.ConnectionInfo) string {
			return mailhtml.TestTemplate(mailhtml.TestData{
				SentAt:  now,
				Channel: info.Channel,
				Status:  info.Status,
			})
		},
		Embeds: []mailer.Embed{
			{
				Name:      "fortune-finger-wordmark.png",
				ContentID: mailhtml.WordmarkContentID,
				Content:   mailhtml.Wordmark(),
			},
		},
	})
}
