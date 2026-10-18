package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/prize/logic"
	"fortunefinger/backend/internal/server/prize/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type PoolHandler struct {
	logic logic.PoolLogic
}

func NewPoolHandler(svcCtx *svc.ServiceContext) *PoolHandler {
	return &PoolHandler{logic: logic.NewPoolLogic(svcCtx)}
}

func (h *PoolHandler) GetPoolListHandler(c *gin.Context) {
	var req types.GetPoolListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品列表参数错误", err))
		return
	}

	resp, err := h.logic.GetPoolList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *PoolHandler) GetPoolDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品ID参数错误", err))
		return
	}

	resp, err := h.logic.GetPoolDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *PoolHandler) CreatePoolHandler(c *gin.Context) {
	var req types.CreatePoolReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品参数错误", err))
		return
	}

	if err := h.logic.CreatePool(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *PoolHandler) UpdatePoolHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品ID参数错误", err))
		return
	}

	var req types.UpdatePoolReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品参数错误", err))
		return
	}

	if err := h.logic.UpdatePool(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *PoolHandler) DeletePoolHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "奖品ID参数错误", err))
		return
	}

	if err := h.logic.DeletePool(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}
