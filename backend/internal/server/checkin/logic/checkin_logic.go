package logic

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	sysconfig "fortunefinger/backend/internal/biz/sys-config"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/checkin/types"
	"fortunefinger/backend/internal/svc"

	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

const (
	checkinConfigCode      = "sign"
	checkinQRCodeLifetime  = 30 * time.Second
	checkinSessionLifetime = 5 * time.Minute
)

type CheckinLogic interface {
	GetCheckinTemplate(ctx context.Context, req *types.GetCheckinTemplateReq) (*types.GetCheckinTemplateResp, error)
	GetCheckinState(ctx context.Context, req *types.GetCheckinStateReq) (*types.GetCheckinStateResp, error)
	GetCheckinQRCode(ctx context.Context, req *types.GetCheckinQRCodeReq) (*types.GetCheckinQRCodeResp, error)
	CreateCheckinSession(ctx context.Context, req *types.CreateCheckinSessionReq) (*types.CreateCheckinSessionResp, error)
}

type checkinLogic struct {
	svcCtx        *svc.ServiceContext
	configService *sysconfig.Service
	activityRepo  repository.LotteryActivityRepository
}

func NewCheckinLogic(svcCtx *svc.ServiceContext) CheckinLogic {
	return &checkinLogic{
		svcCtx:        svcCtx,
		configService: sysconfig.NewService(svcCtx),
		activityRepo:  repository.NewLotteryActivityRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (l *checkinLogic) GetCheckinTemplate(ctx context.Context, req *types.GetCheckinTemplateReq) (*types.GetCheckinTemplateResp, error) {
	activity, err := l.activityRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	return &types.GetCheckinTemplateResp{Template: activity.CheckinTemplate}, nil
}

func (l *checkinLogic) GetCheckinState(ctx context.Context, req *types.GetCheckinStateReq) (*types.GetCheckinStateResp, error) {
	activity, err := l.activityRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	if !activity.ActivityEndTime.After(activity.ActivityStartTime) {
		return nil, errorx.New(errorx.InvalidParameter, "活动结束时间必须晚于活动开始时间")
	}
	startTime := activity.CheckinStartTime
	endTime := activity.CheckinEndTime
	if !endTime.After(startTime) {
		return nil, errorx.New(errorx.InvalidParameter, "签到截止时间必须晚于允许开始签到时间")
	}
	config, err := l.loadCheckinConfig(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	status := "pending"
	if !now.Before(startTime) && now.Before(endTime) {
		status = "active"
	} else if !now.Before(endTime) {
		status = "ended"
	}
	var total *uint64
	if config.Mode == 2 {
		count := uint64(0)
		total = &count
	}

	return &types.GetCheckinStateResp{
		Status:            status,
		ServerTime:        now.Format(time.RFC3339),
		ActivityStartTime: activity.ActivityStartTime.Format(time.RFC3339),
		ActivityEndTime:   activity.ActivityEndTime.Format(time.RFC3339),
		CheckinStartTime:  startTime.Format(time.RFC3339),
		CheckinEndTime:    endTime.Format(time.RFC3339),
		CheckedIn:         0,
		Total:             total,
	}, nil
}

func (l *checkinLogic) GetCheckinQRCode(ctx context.Context, req *types.GetCheckinQRCodeReq) (*types.GetCheckinQRCodeResp, error) {
	state, err := l.GetCheckinState(ctx, &types.GetCheckinStateReq{Code: req.Code})
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
		"v1.%s.%d.%d.%s",
		req.Code,
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

	participantURL, err := url.Parse(strings.TrimRight(l.svcCtx.Config.Participant.BaseURL, "/") + "/checkin")
	if err != nil || participantURL.Scheme == "" || participantURL.Host == "" {
		l.svcCtx.Logger.Error("[CheckinLogic] 参与者端地址配置无效",
			zap.String("base_url", l.svcCtx.Config.Participant.BaseURL),
			zap.Error(err),
		)
		return nil, errorx.New(errorx.InternalServer, "参与者端地址配置无效")
	}
	query := participantURL.Query()
	query.Set("token", token)
	participantURL.RawQuery = query.Encode()

	image, err := qrcode.Encode(participantURL.String(), qrcode.Medium, 360)
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

func (l *checkinLogic) CreateCheckinSession(ctx context.Context, req *types.CreateCheckinSessionReq) (*types.CreateCheckinSessionResp, error) {
	code, err := l.verifyQRCodeToken(req.Token)
	if err != nil {
		l.svcCtx.Logger.Warn("[CheckinLogic] 拒绝无效的签到二维码凭证", zap.Error(err))
		return nil, errorx.New(errorx.Unauthorized, "签到二维码已失效，请重新扫码")
	}

	state, err := l.GetCheckinState(ctx, &types.GetCheckinStateReq{Code: code})
	if err != nil {
		return nil, err
	}
	if state.Status != "active" {
		l.svcCtx.Logger.Warn("[CheckinLogic] 当前签到状态不允许创建签到会话",
			zap.String("status", state.Status),
		)
		return nil, errorx.New(errorx.Conflict, "当前时间不允许签到")
	}

	config, err := l.loadCheckinConfig(ctx)
	if err != nil {
		return nil, err
	}
	sessionID, expiresAt, err := l.createSessionToken()
	if err != nil {
		l.svcCtx.Logger.Error("[CheckinLogic] 创建签到会话失败", zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "创建签到会话失败", err)
	}
	mode := config.Mode
	if mode == 0 {
		mode = 1
	}

	return &types.CreateCheckinSessionResp{
		SessionID: sessionID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		Activity: types.CheckinActivityResp{
			Name:        "FortuneFinger 年度盛典",
			Company:     "FortuneFinger",
			Description: "欢迎来到活动现场，请完善信息完成签到。",
		},
		Config: types.CheckinSessionConfigResp{
			Mode:           mode,
			RequireName:    config.RequireName || mode == 2,
			RequireCompany: config.RequireCompany,
			AllowAvatar:    config.AllowAvatar,
			AllowMessage:   true,
		},
	}, nil
}

func (l *checkinLogic) loadCheckinConfig(ctx context.Context) (*types.CheckinConfig, error) {
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
	return &value, nil
}

func (l *checkinLogic) verifyQRCodeToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 6 || parts[0] != "v1" {
		return "", fmt.Errorf("签到二维码凭证格式无效")
	}
	payload := strings.Join(parts[:5], ".")
	signature, err := base64.RawURLEncoding.DecodeString(parts[5])
	if err != nil {
		return "", fmt.Errorf("签到二维码签名格式无效: %w", err)
	}
	signer := hmac.New(sha256.New, []byte(l.svcCtx.Config.JWT.Secret))
	_, _ = signer.Write([]byte(payload))
	if !hmac.Equal(signature, signer.Sum(nil)) {
		return "", fmt.Errorf("签到二维码签名无效")
	}
	expiresUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil || time.Now().Unix() > expiresUnix {
		return "", fmt.Errorf("签到二维码凭证已过期")
	}
	return parts[1], nil
}

func (l *checkinLogic) createSessionToken() (string, time.Time, error) {
	expiresAt := time.Now().Add(checkinSessionLifetime)
	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		return "", time.Time{}, err
	}
	payload := fmt.Sprintf("s1.%d.%s", expiresAt.Unix(), base64.RawURLEncoding.EncodeToString(nonce))
	signer := hmac.New(sha256.New, []byte(l.svcCtx.Config.JWT.Secret))
	_, _ = signer.Write([]byte(payload))
	return payload + "." + base64.RawURLEncoding.EncodeToString(signer.Sum(nil)), expiresAt, nil
}
