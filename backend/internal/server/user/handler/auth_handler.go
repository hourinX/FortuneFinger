package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/user/logic"
	"fortunefinger/backend/internal/server/user/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	logic logic.AuthLogic
}

func NewAuthHandler(svcCtx *svc.ServiceContext) *AuthHandler {
	return &AuthHandler{logic: logic.NewAuthLogic(svcCtx)}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "登录参数错误", err))
		return
	}

	resp, err := h.logic.Login(c.Request.Context(), &req, c.ClientIP())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	setRefreshTokenCookie(c, resp.RefreshToken, resp.RefreshExpiresAt, resp.Remember)
	httpx.Ok(c, "登录成功", resp)
}

func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "刷新凭证参数错误", err))
		return
	}
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		httpx.Error(c, errorx.New(errorx.Unauthorized, "刷新凭证无效或已过期"))
		return
	}

	resp, err := h.logic.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	setRefreshTokenCookie(c, resp.RefreshToken, resp.RefreshExpiresAt, resp.Remember)
	httpx.Ok(c, "刷新成功", resp)
}

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	var req types.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "退出登录参数错误", err))
		return
	}

	defer clearRefreshTokenCookie(c)
	refreshToken, _ := c.Cookie(refreshTokenCookieName)
	if err := h.logic.Logout(c.Request.Context(), refreshToken); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, "退出成功")
}

func (h *AuthHandler) GetCaptchaHandler(c *gin.Context) {
	var req types.GetCaptchaReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "验证码参数错误", err))
		return
	}

	resp, err := h.logic.GetCaptcha(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *AuthHandler) CreateEnhancedVerificationHandler(c *gin.Context) {
	var req types.CreateEnhancedVerificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "增强验证参数错误", err))
		return
	}

	resp, err := h.logic.CreateEnhancedVerification(c.Request.Context(), &req, c.ClientIP())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *AuthHandler) VerifyEnhancedVerificationHandler(c *gin.Context) {
	var req types.VerifyEnhancedVerificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "增强验证参数错误", err))
		return
	}

	resp, err := h.logic.VerifyEnhancedVerification(c.Request.Context(), &req, c.ClientIP())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}
