package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	enhancedverification "fortunefinger/backend/internal/biz/enhanced-verification"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/pattern"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/user/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogic interface {
	Login(ctx context.Context, req *types.LoginReq, ip string) (*types.LoginResp, error)
	RefreshToken(ctx context.Context, refreshToken string) (*types.RefreshTokenResp, error)
	Logout(ctx context.Context, refreshToken string) error
	GetCaptcha(ctx context.Context) (*types.GetCaptchaResp, error)
	CreateEnhancedVerification(ctx context.Context, req *types.CreateEnhancedVerificationReq, ip string) (*types.CreateEnhancedVerificationResp, error)
	VerifyEnhancedVerification(ctx context.Context, req *types.VerifyEnhancedVerificationReq, ip string) (*types.VerifyEnhancedVerificationResp, error)
}

const loginVerificationScene = "login"

type authLogic struct {
	svcCtx   *svc.ServiceContext
	userRepo repository.UserRepository
	permRepo repository.CasbinRuleRepository
}

func NewAuthLogic(svcCtx *svc.ServiceContext) AuthLogic {
	return &authLogic{
		svcCtx:   svcCtx,
		userRepo: repository.NewUserRepository(svcCtx.GormDB, svcCtx.Snowflake),
		permRepo: repository.NewCasbinRuleRepository(svcCtx.GormDB),
	}
}

func (l *authLogic) Login(ctx context.Context, req *types.LoginReq, ip string) (*types.LoginResp, error) {
	now := time.Now()
	lockedUntil, err := l.svcCtx.LoginSecurity.LockedUntil(ctx, req.Username, ip, now)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 读取登录锁定状态失败", zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录服务暂时不可用", err)
	}
	if !lockedUntil.IsZero() {
		return nil, lockedError(lockedUntil, now)
	}

	required, err := l.svcCtx.LoginSecurity.RequiresEnhancedVerification(ctx, req.Username, ip)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 读取增强验证状态失败", zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录服务暂时不可用", err)
	}
	if required {
		verified, verifyErr := l.svcCtx.EnhancedVerification.Consume(
			ctx,
			req.VerificationToken,
			loginVerificationScene,
			req.Username,
			ip,
		)
		if verifyErr != nil {
			l.svcCtx.Logger.Error("[AuthLogic] 消费增强验证凭证失败", zap.String("username", req.Username), zap.String("ip", ip), zap.Error(verifyErr))
			return nil, errorx.Wrap(errorx.ServiceUnavailable, "增强验证服务暂时不可用", verifyErr)
		}
		if !verified {
			return nil, errorx.New(errorx.EnhancedVerificationRequired, "")
		}
	}

	answer, exists, err := l.svcCtx.Redis.GetAndDelete(ctx, constant.CaptchaKeyPrefix+req.CaptchaID)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 读取登录验证码失败", zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "验证码服务暂时不可用", err)
	}
	if !exists || answer != req.CaptchaCode {
		return nil, errorx.New(errorx.InvalidParameter, "验证码错误或已过期")
	}

	findUser := l.userRepo.FindByUsername
	if strings.Contains(req.Username, "@") {
		findUser = l.userRepo.FindByEmail
	} else if pattern.Mobile.MatchString(req.Username) {
		findUser = l.userRepo.FindByMobile
	}

	user, err := findUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		state, riskErr := l.svcCtx.LoginSecurity.RecordFailure(ctx, req.Username, ip, now)
		if riskErr != nil {
			l.svcCtx.Logger.Error("[AuthLogic] 记录登录失败次数失败", zap.String("username", req.Username), zap.String("ip", ip), zap.Error(riskErr))
			return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录服务暂时不可用", riskErr)
		}
		if state.LockedUntil.After(now) {
			return nil, lockedError(state.LockedUntil, now)
		}
		return nil, errorx.New(
			errorx.Unauthorized,
			fmt.Sprintf("账号或密码错误，还可尝试%d次", state.RemainingAttempts),
		)
	}

	if user.Status != 1 {
		return nil, errorx.New(errorx.Forbidden, "账号已被禁用")
	}

	roleCode, err := l.permRepo.FindRoleCodeByUserID(ctx, strconv.FormatUint(user.ID, 10))
	if err != nil {
		return nil, err
	}
	if roleCode == "" {
		return nil, errorx.New(errorx.Forbidden, "当前账号未分配角色，请联系管理员")
	}

	menuIDs, err := l.permRepo.ListMenuIDsByRoleCode(ctx, roleCode)
	if err != nil {
		return nil, err
	}
	if len(menuIDs) == 0 {
		return nil, errorx.New(errorx.Forbidden, "请先联系管理员为当前角色分配权限")
	}
	if err := l.svcCtx.LoginSecurity.Reset(ctx, req.Username, ip); err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 清空登录失败记录失败", zap.Uint64("user_id", user.ID), zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录状态重置失败", err)
	}

	tokens, err := l.svcCtx.JWT.Generate(user.ID, req.Remember)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 生成登录凭证失败", zap.Uint64("user_id", user.ID), zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成登录凭证失败", err)
	}

	refreshExpire := l.svcCtx.Config.JWT.RefreshExpire
	if req.Remember {
		refreshExpire = l.svcCtx.Config.JWT.RememberRefreshExpire
	}
	if err := l.saveRefreshToken(ctx, tokens.RefreshToken, user.ID, refreshExpire); err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 保存登录刷新凭证失败", zap.Uint64("user_id", user.ID), zap.String("username", req.Username), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录状态保存失败", err)
	}

	resp := types.LoginResp{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresAt:        now.Add(l.svcCtx.Config.JWT.AccessExpire).Unix(),
		RefreshExpiresAt: now.Add(refreshExpire).Unix(),
		Remember:         req.Remember,
	}
	return &resp, nil
}

func (l *authLogic) RefreshToken(ctx context.Context, refreshToken string) (*types.RefreshTokenResp, error) {
	identity, err := l.svcCtx.JWT.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, errorx.New(errorx.Unauthorized, "刷新凭证无效或已过期")
	}

	oldRefreshKey := refreshTokenKey(refreshToken)
	userID, exists, err := l.svcCtx.Redis.Get(ctx, oldRefreshKey)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 读取刷新凭证会话失败", zap.Uint64("user_id", identity.UserID), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录状态读取失败", err)
	}
	if !exists || userID != strconv.FormatUint(identity.UserID, 10) {
		return nil, errorx.New(errorx.Unauthorized, "刷新凭证无效或已过期")
	}

	user, err := l.userRepo.FindByID(ctx, identity.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New(errorx.Unauthorized, "登录用户不存在")
	}
	if user.Status != 1 {
		return nil, errorx.New(errorx.Forbidden, "账号已被禁用")
	}

	now := time.Now()
	tokens, err := l.svcCtx.JWT.Generate(user.ID, identity.Remember)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 刷新登录凭证生成失败", zap.Uint64("user_id", user.ID), zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "刷新登录凭证失败", err)
	}

	refreshExpire := l.svcCtx.Config.JWT.RefreshExpire
	if identity.Remember {
		refreshExpire = l.svcCtx.Config.JWT.RememberRefreshExpire
	}

	rotated, err := l.svcCtx.Redis.ReplaceIfValue(
		ctx,
		oldRefreshKey,
		refreshTokenKey(tokens.RefreshToken),
		userID,
		userID,
		refreshExpire,
	)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 轮换登录刷新凭证失败", zap.Uint64("user_id", user.ID), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "登录状态保存失败", err)
	}

	if !rotated {
		return nil, errorx.New(errorx.Unauthorized, "刷新凭证无效或已过期")
	}

	resp := types.RefreshTokenResp{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresAt:        now.Add(l.svcCtx.Config.JWT.AccessExpire).Unix(),
		RefreshExpiresAt: now.Add(refreshExpire).Unix(),
		Remember:         identity.Remember,
	}
	return &resp, nil
}

func (l *authLogic) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	if _, _, err := l.svcCtx.Redis.GetAndDelete(ctx, refreshTokenKey(refreshToken)); err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 删除退出登录凭证失败", zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "退出登录失败", err)
	}
	return nil
}

func (l *authLogic) GetCaptcha(ctx context.Context) (*types.GetCaptchaResp, error) {
	id, image, answer, err := l.svcCtx.Captcha.Generate()
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 生成验证码失败", zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成验证码失败", err)
	}

	if err := l.svcCtx.Redis.Set(ctx, constant.CaptchaKeyPrefix+id, answer, l.svcCtx.Captcha.TTL()); err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 保存验证码失败", zap.String("captcha_id", id), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "验证码服务暂时不可用", err)
	}
	resp := types.GetCaptchaResp{CaptchaID: id, Captcha: image}
	return &resp, nil
}

func (l *authLogic) CreateEnhancedVerification(ctx context.Context, req *types.CreateEnhancedVerificationReq, ip string) (*types.CreateEnhancedVerificationResp, error) {
	if req.Scene == loginVerificationScene {
		required, err := l.svcCtx.LoginSecurity.RequiresEnhancedVerification(ctx, req.Subject, ip)
		if err != nil {
			l.svcCtx.Logger.Error("[AuthLogic] 读取增强验证状态失败", zap.String("username", req.Subject), zap.String("ip", ip), zap.Error(err))
			return nil, errorx.Wrap(errorx.ServiceUnavailable, "增强验证服务暂时不可用", err)
		}
		if !required {
			return nil, errorx.New(errorx.Conflict, "当前操作无需增强验证")
		}
	}

	verification, err := l.svcCtx.EnhancedVerification.Create(ctx, req.Scene, req.Subject, ip)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 创建增强验证失败", zap.String("scene", req.Scene), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "增强验证服务暂时不可用", err)
	}

	resp := types.CreateEnhancedVerificationResp{
		VerificationID: verification.ID,
		Mode:           verification.Mode,
		Image:          verification.Image,
		Prompt:         verification.Prompt,
		Thumb:          verification.Thumb,
		ThumbX:         verification.ThumbX,
		ThumbY:         verification.ThumbY,
		ThumbWidth:     verification.ThumbWidth,
		ThumbHeight:    verification.ThumbHeight,
		RequiredPoints: verification.RequiredPoints,
		ExpiresIn:      verification.ExpiresIn,
	}
	return &resp, nil
}

func (l *authLogic) VerifyEnhancedVerification(ctx context.Context, req *types.VerifyEnhancedVerificationReq, ip string) (*types.VerifyEnhancedVerificationResp, error) {
	points := make([]enhancedverification.Point, 0, len(req.Points))
	for _, point := range req.Points {
		points = append(points, enhancedverification.Point{
			X: point.X,
			Y: point.Y,
			T: point.T,
		})
	}

	token, verified, err := l.svcCtx.EnhancedVerification.Verify(
		ctx,
		req.VerificationID,
		req.Scene,
		req.Subject,
		ip,
		points,
	)
	if err != nil {
		l.svcCtx.Logger.Error("[AuthLogic] 校验增强验证失败", zap.String("scene", req.Scene), zap.String("ip", ip), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "增强验证服务暂时不可用", err)
	}
	if !verified {
		return nil, errorx.New(errorx.InvalidParameter, "增强验证未通过，请重试")
	}

	resp := types.VerifyEnhancedVerificationResp{
		VerificationToken: token,
		ExpiresIn:         120,
	}
	return &resp, nil
}

func lockedError(lockedUntil, now time.Time) error {
	remaining := lockedUntil.Sub(now)
	minutes := int64((remaining + time.Minute - 1) / time.Minute)
	if minutes < 1 {
		minutes = 1
	}
	if minutes < 60 {
		return errorx.New(errorx.TooManyRequests, fmt.Sprintf("登录失败次数过多，请%d分钟后重试", minutes))
	}
	hours := (minutes + 59) / 60
	return errorx.New(errorx.TooManyRequests, fmt.Sprintf("登录失败次数过多，请%d小时后重试", hours))
}

func (l *authLogic) saveRefreshToken(ctx context.Context, refreshToken string, userID uint64, expiration time.Duration) error {
	return l.svcCtx.Redis.Set(
		ctx,
		refreshTokenKey(refreshToken),
		strconv.FormatUint(userID, 10),
		expiration,
	)
}

func refreshTokenKey(refreshToken string) string {
	digest := sha256.Sum256([]byte(refreshToken))
	return constant.RefreshTokenKeyPrefix + hex.EncodeToString(digest[:])
}
