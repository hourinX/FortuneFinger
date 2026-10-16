package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/user/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUser(r *gin.Engine, h *handler.UserHandler, check mw.Checker) {
	user := r.Group("/api/user")
	{
		user.GET("/list", mw.CheckPermission("perm:user:list", check), h.GetUserListHandler)
		user.GET("/detail/:id", mw.CheckPermission("perm:user:info", check), h.GetUserDetailHandler)
		user.POST("/create", mw.CheckPermission("perm:user:add", check), h.CreateUserHandler)
		user.PUT("/update/:id", mw.CheckPermission("perm:user:edit", check), h.UpdateUserHandler)
		user.DELETE("/:id", mw.CheckPermission("perm:user:del", check), h.DeleteUserHandler)
		user.PUT("/role/:id", mw.CheckPermission("perm:user:role", check), h.AssignUserRoleHandler)
		user.GET("/getUserInfo", h.GetUserInfoHandler)
		user.PUT("/updateUserInfo", h.UpdateUserInfoHandler)
	}
}
