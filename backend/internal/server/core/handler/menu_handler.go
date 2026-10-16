package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/core/logic"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	logic logic.MenuLogic
}

func NewMenuHandler(svcCtx *svc.ServiceContext) *MenuHandler {
	return &MenuHandler{logic: logic.NewMenuLogic(svcCtx)}
}

func (h *MenuHandler) GetMenuListHandler(c *gin.Context) {
	var req types.GetMenuListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单列表参数错误", err))
		return
	}

	resp, err := h.logic.GetMenuList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *MenuHandler) GetMenuDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单ID参数错误", err))
		return
	}

	resp, err := h.logic.GetMenuDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *MenuHandler) CreateMenuHandler(c *gin.Context) {
	var req types.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单参数错误", err))
		return
	}

	if err := h.logic.CreateMenu(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *MenuHandler) UpdateMenuHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单ID参数错误", err))
		return
	}

	var req types.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单参数错误", err))
		return
	}

	if err := h.logic.UpdateMenu(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *MenuHandler) DeleteMenuHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "菜单ID参数错误", err))
		return
	}

	if err := h.logic.DeleteMenu(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *MenuHandler) GetRoutersHandler(c *gin.Context) {
	var req types.GetRoutersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "路由参数错误", err))
		return
	}

	resp, err := h.logic.GetRouters(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *MenuHandler) GetPermissionHandler(c *gin.Context) {
	var req types.GetPermissionReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "权限码参数错误", err))
		return
	}

	resp, err := h.logic.GetPermission(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *MenuHandler) GetPermissionTreeHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "角色ID参数错误", err))
		return
	}

	resp, err := h.logic.GetPermissionTree(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}
