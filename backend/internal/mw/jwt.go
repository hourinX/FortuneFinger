package mw

import (
	"strconv"
	"strings"

	"fortunefinger/backend/internal/pkg/auth"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/pkg/requestctx/userid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWT(manager *auth.Manager, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Del(userid.Header)

		if isPublicRoute(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn(
				"HTTP身份认证请求头无效",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			writeUnauthorized(c)
			return
		}

		userID, err := manager.VerifyAccessToken(parts[1])
		if err != nil {
			log.Warn(
				"HTTP访问凭证校验失败",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.Error(err),
			)
			writeUnauthorized(c)
			return
		}

		c.Request.Header.Set(userid.Header, strconv.FormatUint(userID, 10))
		c.Next()
	}
}

func writeUnauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	c.Abort()
	httpx.Error(c, errorx.New(errorx.Unauthorized, "未登录或登录已过期"))
}
