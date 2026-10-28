package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/checkin/logic"
	"fortunefinger/backend/internal/server/checkin/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type CheckinHandler struct {
	logic logic.CheckinLogic
}

func NewCheckinHandler(svcCtx *svc.ServiceContext) *CheckinHandler {
	return &CheckinHandler{logic: logic.NewCheckinLogic(svcCtx)}
}

func (h *CheckinHandler) GetCheckinTemplateHandler(c *gin.Context) {
	var req types.GetCheckinTemplateReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "签到大屏模板参数错误", err))
		return
	}

	resp, err := h.logic.GetCheckinTemplate(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *CheckinHandler) GetCheckinStateHandler(c *gin.Context) {
	var req types.GetCheckinStateReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "签到状态参数错误", err))
		return
	}

	resp, err := h.logic.GetCheckinState(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *CheckinHandler) GetCheckinQRCodeHandler(c *gin.Context) {
	var req types.GetCheckinQRCodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "签到二维码参数错误", err))
		return
	}

	resp, err := h.logic.GetCheckinQRCode(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}
