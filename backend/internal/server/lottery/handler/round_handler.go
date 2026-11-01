package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/lottery/logic"
	"fortunefinger/backend/internal/server/lottery/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type RoundHandler struct {
	logic logic.RoundLogic
}

func NewRoundHandler(svcCtx *svc.ServiceContext) *RoundHandler {
	return &RoundHandler{logic: logic.NewRoundLogic(svcCtx)}
}

func (h *RoundHandler) GetRoundListHandler(c *gin.Context) {
	var req types.GetRoundListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次列表参数错误", err))
		return
	}

	resp, err := h.logic.GetRoundList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *RoundHandler) GetRoundDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次ID参数错误", err))
		return
	}

	resp, err := h.logic.GetRoundDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *RoundHandler) CreateRoundHandler(c *gin.Context) {
	var req types.CreateRoundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次参数错误", err))
		return
	}

	if err := h.logic.CreateRound(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *RoundHandler) UpdateRoundHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次ID参数错误", err))
		return
	}

	var req types.UpdateRoundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次参数错误", err))
		return
	}

	if err := h.logic.UpdateRound(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *RoundHandler) DeleteRoundHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖轮次ID参数错误", err))
		return
	}

	if err := h.logic.DeleteRound(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}
