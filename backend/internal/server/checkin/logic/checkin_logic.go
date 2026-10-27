package logic

import (
	"context"
	"encoding/json"

	sysconfig "fortunefinger/backend/internal/biz/sys-config"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/server/checkin/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

const (
	checkinConfigCode      = "sign"
	defaultCheckinTemplate = "invitation"
)

type CheckinLogic interface {
	GetCheckinTemplate(ctx context.Context, req *types.GetCheckinTemplateReq) (*types.GetCheckinTemplateResp, error)
}

type checkinLogic struct {
	svcCtx        *svc.ServiceContext
	configService *sysconfig.Service
}

func NewCheckinLogic(svcCtx *svc.ServiceContext) CheckinLogic {
	return &checkinLogic{
		svcCtx:        svcCtx,
		configService: sysconfig.NewService(svcCtx),
	}
}

func (l *checkinLogic) GetCheckinTemplate(ctx context.Context, _ *types.GetCheckinTemplateReq) (*types.GetCheckinTemplateResp, error) {
	config, err := l.configService.GetConfigByCode(ctx, checkinConfigCode)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &types.GetCheckinTemplateResp{Template: defaultCheckinTemplate}, nil
	}

	var value types.CheckinConfig
	if err := json.Unmarshal([]byte(config.Config), &value); err != nil {
		l.svcCtx.Logger.Error("[CheckinLogic] 解析签到配置失败",
			zap.Uint64("config_id", config.ID),
			zap.String("config_code", config.Code),
			zap.Error(err),
		)
		return nil, errorx.Wrap(errorx.InternalServer, "签到配置内容无效", err)
	}

	template := value.Screen.Template
	if template == "" {
		template = defaultCheckinTemplate
	}
	return &types.GetCheckinTemplateResp{Template: template}, nil
}
