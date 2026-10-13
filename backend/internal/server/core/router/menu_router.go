package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMenu(r *gin.Engine, h *handler.MenuHandler, check mw.Checker) {
	menu := r.Group("/api/menu")
	{
		menu.GET("/list", h.GetMenuListHandler)
		menu.GET("/detail/:id", h.GetMenuDetailHandler)
		menu.POST("/create", h.CreateMenuHandler)
		menu.PUT("/update/:id", h.UpdateMenuHandler)
		menu.DELETE("/:id", h.DeleteMenuHandler)

		menu.GET("/getRouters", h.GetRoutersHandler)
		menu.GET("/permission/:id", mw.CheckPermission("perm:role:info", check), h.GetPermissionTreeHandler)
	}
}
