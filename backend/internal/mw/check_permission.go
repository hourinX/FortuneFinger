package mw

import (
	"context"

	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Checker func(ctx context.Context, userID uint64, code string) (bool, error)

func CheckPermission(code string, check Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if check == nil {
			httpx.Error(c, errorx.New(errorx.ServiceUnavailable, "权限服务不可用"))
			c.Abort()
			return
		}

		userID, err := operator.UserID(c.Request.Context())
		if err != nil {
			httpx.Error(c, err)
			c.Abort()
			return
		}

		allowed, err := check(c.Request.Context(), userID, code)
		if err != nil {
			httpx.Error(c, err)
			c.Abort()
			return
		}
		if !allowed {
			httpx.Error(c, errorx.New(errorx.Forbidden, "当前操作用户不具备该权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}
