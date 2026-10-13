package requestctx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type headersKey struct{}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), headersKey{}, c.Request.Header.Clone())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetHeader(ctx context.Context, name string) string {
	headers, _ := ctx.Value(headersKey{}).(http.Header)
	return headers.Get(name)
}
