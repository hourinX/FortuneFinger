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
		user.GET("/detail/:id", h.GetUserDetailHandler)
		user.POST("/create", h.CreateUserHandler)
		user.PUT("/update/:id", h.UpdateUserHandler)
		user.DELETE("/:id", h.DeleteUserHandler)

		user.PUT("/role/:id", h.AssignUserRoleHandler)
		user.GET("/getUserInfo", h.GetUserInfoHandler)
	}
}
