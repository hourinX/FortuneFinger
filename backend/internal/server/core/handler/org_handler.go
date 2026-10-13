package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/core/logic"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	logic logic.OrgLogic
}

func NewOrgHandler(svcCtx *svc.ServiceContext) *OrgHandler {
	return &OrgHandler{logic: logic.NewOrgLogic(svcCtx)}
}

func (h *OrgHandler) GetOrgListHandler(c *gin.Context) {
	var req types.GetOrgListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织列表参数错误", err))
		return
	}

	resp, err := h.logic.GetOrgList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *OrgHandler) GetOrgDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织ID参数错误", err))
		return
	}

	resp, err := h.logic.GetOrgDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *OrgHandler) CreateOrgHandler(c *gin.Context) {
	var req types.CreateOrgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织参数错误", err))
		return
	}

	if err := h.logic.CreateOrg(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *OrgHandler) UpdateOrgHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织ID参数错误", err))
		return
	}

	var req types.UpdateOrgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织参数错误", err))
		return
	}

	if err := h.logic.UpdateOrg(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *OrgHandler) DeleteOrgHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织ID参数错误", err))
		return
	}

	if err := h.logic.DeleteOrg(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *OrgHandler) GetOrgOptionsHandler(c *gin.Context) {
	var req types.GetOrgOptionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "组织选项参数错误", err))
		return
	}

	resp, err := h.logic.GetOrgOptions(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}
