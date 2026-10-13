package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRole(r *gin.Engine, h *handler.RoleHandler, check mw.Checker) {
	role := r.Group("/api/role")
	{
		role.GET("/list", mw.CheckPermission("perm:role:list", check), h.GetRoleListHandler)
		role.GET("/detail/:id", mw.CheckPermission("perm:role:info", check), h.GetRoleDetailHandler)
		role.POST("/create", mw.CheckPermission("perm:role:add", check), h.CreateRoleHandler)
		role.PUT("/update/:id", mw.CheckPermission("perm:role:edit", check), h.UpdateRoleHandler)
		role.DELETE("/:id", mw.CheckPermission("perm:role:del", check), h.DeleteRoleHandler)

		role.GET("/options", h.GetRoleOptionsHandler)
		role.PUT("/permission/:id", mw.CheckPermission("perm:role:edit", check), h.UpdateRolePermissionHandler)
	}
}
