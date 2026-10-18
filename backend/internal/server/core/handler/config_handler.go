package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/core/logic"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	logic logic.ConfigLogic
}

func NewConfigHandler(svcCtx *svc.ServiceContext) *ConfigHandler {
	return &ConfigHandler{logic: logic.NewConfigLogic(svcCtx)}
}

func (h *ConfigHandler) GetConfigHandler(c *gin.Context) {
	var req types.GetConfigReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "系统配置参数错误", err))
		return
	}

	resp, err := h.logic.GetConfig(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *ConfigHandler) SaveConfigHandler(c *gin.Context) {
	var req types.SaveConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "系统配置参数错误", err))
		return
	}

	resp, err := h.logic.SaveConfig(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess, resp)
}

func (h *ConfigHandler) SendTestMailHandler(c *gin.Context) {
	var req types.SendTestMailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "测试邮件参数错误", err))
		return
	}

	err := h.logic.SendTestMail(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, "发送成功")
}
