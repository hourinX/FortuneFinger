package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMenu(r *gin.Engine, h *handler.MenuHandler, check mw.Checker) {
	menu := r.Group("/api/menu")
	{
		menu.GET("/list", mw.CheckPermission("system:menu:list", check), h.GetMenuListHandler)
		menu.GET("/detail/:id", mw.CheckPermission("system:menu:info", check), h.GetMenuDetailHandler)
		menu.POST("/create", mw.CheckPermission("system:menu:add", check), h.CreateMenuHandler)
		menu.PUT("/update/:id", mw.CheckPermission("system:menu:edit", check), h.UpdateMenuHandler)
		menu.DELETE("/:id", mw.CheckPermission("system:menu:del", check), h.DeleteMenuHandler)
		menu.GET("/getRouters", h.GetRoutersHandler)
		menu.GET("/permission/:id", h.GetPermissionTreeHandler)
	}
}
