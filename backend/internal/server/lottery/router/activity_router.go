package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/lottery/handler"

	"github.com/gin-gonic/gin"
)

func RegisterActivity(r *gin.Engine, h *handler.ActivityHandler, check mw.Checker) {
	activity := r.Group("/api/lottery/activity")
	{
		activity.GET("/list", mw.CheckPermission("lottery:activities:list", check), h.GetActivityListHandler)
		activity.GET("/detail/:id", mw.CheckPermission("lottery:activities:info", check), h.GetActivityDetailHandler)
		activity.POST("/create", mw.CheckPermission("lottery:activities:add", check), h.CreateActivityHandler)
		activity.PUT("/update/:id", mw.CheckPermission("lottery:activities:edit", check), h.UpdateActivityHandler)
		activity.DELETE("/:id", mw.CheckPermission("lottery:activities:del", check), h.DeleteActivityHandler)
		activity.PUT("/template/update/:id", mw.CheckPermission("lottery:activities:template", check), h.UpdateActivityCheckinTemplateHandler)
	}
}
