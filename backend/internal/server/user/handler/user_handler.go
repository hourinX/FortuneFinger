package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/user/logic"
	"fortunefinger/backend/internal/server/user/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	logic logic.UserLogic
}

func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{logic: logic.NewUserLogic(svcCtx)}
}

func (h *UserHandler) GetUserListHandler(c *gin.Context) {
	var req types.GetUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户列表参数错误", err))
		return
	}

	resp, err := h.logic.GetUserList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *UserHandler) GetUserDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户ID参数错误", err))
		return
	}

	resp, err := h.logic.GetUserDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var req types.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户参数错误", err))
		return
	}

	if err := h.logic.CreateUser(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户ID参数错误", err))
		return
	}

	var req types.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户参数错误", err))
		return
	}

	if err := h.logic.UpdateUser(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户ID参数错误", err))
		return
	}

	if err := h.logic.DeleteUser(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *UserHandler) AssignUserRoleHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户ID参数错误", err))
		return
	}

	var req types.AssignUserRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "用户角色参数错误", err))
		return
	}

	if err := h.logic.AssignUserRole(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *UserHandler) GetUserInfoHandler(c *gin.Context) {
	var req types.GetUserInfoReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "当前用户参数错误", err))
		return
	}

	resp, err := h.logic.GetUserInfo(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *UserHandler) UpdateUserInfoHandler(c *gin.Context) {
	var req types.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "个人信息参数错误", err))
		return
	}

	if err := h.logic.UpdateUserInfo(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}
