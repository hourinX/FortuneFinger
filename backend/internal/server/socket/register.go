package socket

import (
	"fortunefinger/backend/internal/pkg/ws"
	"fortunefinger/backend/internal/server/socket/handler"
	socketRouter "fortunefinger/backend/internal/server/socket/router"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, manager *ws.Manager) {
	socketRouter.RegisterSocket(r, handler.NewSocketHandler(manager))
}
