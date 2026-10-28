package logic

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	sysconfig "fortunefinger/backend/internal/biz/sys-config"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/server/checkin/types"
	"fortunefinger/backend/internal/svc"

	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

const (
	checkinConfigCode      = "sign"
	defaultCheckinTemplate = "invitation"
	checkinTimeLayout      = "2006-01-02 15:04:05"
	checkinQRCodeLifetime  = 30 * time.Second
)

type CheckinLogic interface {
	GetCheckinTemplate(ctx context.Context, req *types.GetCheckinTemplateReq) (*types.GetCheckinTemplateResp, error)
	GetCheckinState(ctx context.Context, req *types.GetCheckinStateReq) (*types.GetCheckinStateResp, error)
	GetCheckinQRCode(ctx context.Context, req *types.GetCheckinQRCodeReq) (*types.GetCheckinQRCodeResp, error)
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

func (l *checkinLogic) GetCheckinState(ctx context.Context, _ *types.GetCheckinStateReq) (*types.GetCheckinStateResp, error) {
	config, err := l.configService.GetConfigByCode(ctx, checkinConfigCode)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errorx.New(errorx.NotFound, "签到配置不存在")
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

	startTime, err := time.ParseInLocation(checkinTimeLayout, value.CheckinStartTime, time.Local)
	if err != nil {
		return nil, errorx.New(errorx.InvalidParameter, "允许开始签到时间无效")
	}
	endTime, err := time.ParseInLocation(checkinTimeLayout, value.CheckinEndTime, time.Local)
	if err != nil {
		return nil, errorx.New(errorx.InvalidParameter, "签到截止时间无效")
	}
	if !endTime.After(startTime) {
		return nil, errorx.New(errorx.InvalidParameter, "签到截止时间必须晚于允许开始签到时间")
	}

	now := time.Now()
	status := "pending"
	if !now.Before(startTime) && now.Before(endTime) {
		status = "active"
	} else if !now.Before(endTime) {
		status = "ended"
	}
	var total *uint64
	if value.Mode == 2 {
		count := uint64(0)
		total = &count
	}

	return &types.GetCheckinStateResp{
		Status:           status,
		ServerTime:       now.Format(time.RFC3339),
		CheckinStartTime: startTime.Format(time.RFC3339),
		CheckinEndTime:   endTime.Format(time.RFC3339),
		CheckedIn:        0,
		Total:            total,
	}, nil
}

func (l *checkinLogic) GetCheckinQRCode(ctx context.Context, _ *types.GetCheckinQRCodeReq) (*types.GetCheckinQRCodeResp, error) {
	state, err := l.GetCheckinState(ctx, &types.GetCheckinStateReq{})
	if err != nil {
		return nil, err
	}
	if state.Status != "active" {
		l.svcCtx.Logger.Warn("[CheckinLogic] 当前时间不允许签发签到二维码",
			zap.String("status", state.Status),
			zap.String("checkin_start_time", state.CheckinStartTime),
			zap.String("checkin_end_time", state.CheckinEndTime),
		)
		return nil, errorx.New(errorx.Conflict, "当前时间不允许签到")
	}

	now := time.Now()
	expiresAt := now.Add(checkinQRCodeLifetime)
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		l.svcCtx.Logger.Error("[CheckinLogic] 生成签到二维码随机值失败", zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成签到二维码失败", err)
	}

	payload := fmt.Sprintf(
		"v1.%d.%d.%s",
		now.Unix(),
		expiresAt.Unix(),
		base64.RawURLEncoding.EncodeToString(nonce),
	)
	signer := hmac.New(sha256.New, []byte(l.svcCtx.Config.JWT.Secret))
	if _, err := signer.Write([]byte(payload)); err != nil {
		l.svcCtx.Logger.Error("[CheckinLogic] 签名签到二维码失败", zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成签到二维码失败", err)
	}
	signature := base64.RawURLEncoding.EncodeToString(signer.Sum(nil))
	token := payload + "." + signature

	// TODO(checkin-h5): H5 签到入口确定后，增加可配置的签到页面基础地址，并将二维码内容改为“完整 H5 地址 + token 查询参数”。
	// 正式公网部署使用域名和 HTTPS；本地测试或现场内网模式允许配置同一局域网可访问的 IP 地址，使系统相机、微信等扫码工具
	// 能够直接打开签到页面。当前阶段先编码短期凭证本身。
	image, err := qrcode.Encode(token, qrcode.Medium, 360)
	if err != nil {
		l.svcCtx.Logger.Error("[CheckinLogic] 生成签到二维码图片失败", zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成签到二维码失败", err)
	}

	// TODO(checkin-identity): H5 提交签到时，将 FingerprintJS 等浏览器指纹库生成的标识与 UA、屏幕参数、Canvas 特征组合后哈希，
	// 再结合所属企业、真实姓名及手机号或邮箱写入签到记录。同一设备或浏览器默认只允许代入一名参与者；命中重复或高风险特征时，
	// 升级为邮箱验证码校验。该方案配合 30 秒动态二维码，目标是拦截约八成常见的随手代签，但不能作为强身份认证或保证绝对准确率。
	return &types.GetCheckinQRCodeResp{
		Token:        token,
		QRCode:       "data:image/png;base64," + base64.StdEncoding.EncodeToString(image),
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		RefreshAfter: uint8(checkinQRCodeLifetime / time.Second),
	}, nil
}
