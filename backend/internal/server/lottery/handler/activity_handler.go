package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/lottery/logic"
	"fortunefinger/backend/internal/server/lottery/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	logic logic.ActivityLogic
}

func NewActivityHandler(svcCtx *svc.ServiceContext) *ActivityHandler {
	return &ActivityHandler{logic: logic.NewActivityLogic(svcCtx)}
}

func (h *ActivityHandler) GetActivityListHandler(c *gin.Context) {
	var req types.GetActivityListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动列表参数错误", err))
		return
	}

	resp, err := h.logic.GetActivityList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *ActivityHandler) GetActivityDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动ID参数错误", err))
		return
	}

	resp, err := h.logic.GetActivityDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *ActivityHandler) CreateActivityHandler(c *gin.Context) {
	var req types.CreateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动参数错误", err))
		return
	}

	if err := h.logic.CreateActivity(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *ActivityHandler) UpdateActivityHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动ID参数错误", err))
		return
	}

	var req types.UpdateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动参数错误", err))
		return
	}

	if err := h.logic.UpdateActivity(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *ActivityHandler) DeleteActivityHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动ID参数错误", err))
		return
	}

	if err := h.logic.DeleteActivity(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *ActivityHandler) UpdateActivityCheckinTemplateHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "抽奖活动ID参数错误", err))
		return
	}

	var req types.UpdateActivityCheckinTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "签到模板参数错误", err))
		return
	}

	if err := h.logic.UpdateActivityCheckinTemplate(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}
