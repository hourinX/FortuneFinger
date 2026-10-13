package mw

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"syscall"

	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}
			if recovered == http.ErrAbortHandler {
				panic(recovered)
			}
			if err, broken := brokenConnectionError(recovered); broken {
				log.Warn(
					"HTTP客户端连接已断开",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Error(err),
				)
				c.Abort()
				return
			}

			log.Error(
				"HTTP请求发生panic",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("route", c.FullPath()),
				zap.Any("panic", recovered),
				zap.Stack("stack"),
			)
			c.Abort()
			if !c.Writer.Written() {
				httpx.Error(c, errorx.New(errorx.InternalServer, ""))
			}
		}()
		c.Next()
	}
}

func brokenConnectionError(recovered any) (error, bool) {
	err, ok := recovered.(error)
	if !ok {
		return nil, false
	}

	if errors.Is(err, net.ErrClosed) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.ECONNABORTED) {
		return err, true
	}

	message := strings.ToLower(err.Error())
	broken := strings.Contains(message, "broken pipe") ||
		strings.Contains(message, "connection reset by peer") ||
		strings.Contains(message, "forcibly closed by the remote host")
	return err, broken
}
