package router

import (
	"fortunefinger/backend/internal/server/user/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(r *gin.Engine, h *handler.AuthHandler) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", h.LoginHandler)
		auth.POST("/refresh", h.RefreshTokenHandler)
		auth.POST("/logout", h.LogoutHandler)
		auth.GET("/getCaptcha", h.GetCaptchaHandler)
		auth.POST("/verification/create", h.CreateEnhancedVerificationHandler)
		auth.POST("/verification/verify", h.VerifyEnhancedVerificationHandler)
	}
}
