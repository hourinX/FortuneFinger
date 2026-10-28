package router

import (
	"fortunefinger/backend/internal/server/socket/handler"

	"github.com/gin-gonic/gin"
)

func RegisterSocket(r *gin.Engine, h *handler.SocketHandler) {
	r.GET("/api/ws", h.ConnectHandler)
}
