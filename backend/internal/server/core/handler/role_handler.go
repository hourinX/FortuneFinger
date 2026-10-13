package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/core/logic"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	logic logic.RoleLogic
}

func NewRoleHandler(svcCtx *svc.ServiceContext) *RoleHandler {
	return &RoleHandler{logic: logic.NewRoleLogic(svcCtx)}
}

func (h *RoleHandler) GetRoleListHandler(c *gin.Context) {
	var req types.GetRoleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色列表参数错误", err))
		return
	}

	resp, err := h.logic.GetRoleList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *RoleHandler) GetRoleDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色ID参数错误", err))
		return
	}

	resp, err := h.logic.GetRoleDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *RoleHandler) CreateRoleHandler(c *gin.Context) {
	var req types.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色参数错误", err))
		return
	}

	if err := h.logic.CreateRole(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *RoleHandler) UpdateRoleHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色ID参数错误", err))
		return
	}

	var req types.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色参数错误", err))
		return
	}

	if err := h.logic.UpdateRole(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *RoleHandler) DeleteRoleHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色ID参数错误", err))
		return
	}

	if err := h.logic.DeleteRole(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *RoleHandler) GetRoleOptionsHandler(c *gin.Context) {
	var req types.GetRoleOptionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色选项参数错误", err))
		return
	}

	resp, err := h.logic.GetRoleOptions(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *RoleHandler) UpdateRolePermissionHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色ID参数错误", err))
		return
	}

	var req types.UpdateRolePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色权限参数错误", err))
		return
	}

	if err := h.logic.UpdateRolePermission(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}
