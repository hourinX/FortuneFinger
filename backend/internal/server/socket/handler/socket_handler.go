package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/pkg/ws"
	"fortunefinger/backend/internal/server/socket/types"

	"github.com/gin-gonic/gin"
)

type SocketHandler struct {
	manager *ws.Manager
}

func NewSocketHandler(manager *ws.Manager) *SocketHandler {
	return &SocketHandler{manager: manager}
}

func (h *SocketHandler) ConnectHandler(c *gin.Context) {
	var req types.ConnectReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "WebSocket连接参数错误", err))
		return
	}

	if err := h.manager.Serve(c.Request.Context(), c.Writer, c.Request, req.Key); err != nil {
		httpx.Error(c, err)
		return
	}
}
