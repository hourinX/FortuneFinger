package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterConfig(r *gin.Engine, h *handler.ConfigHandler, check mw.Checker) {
	config := r.Group("/api/config")
	{
		config.GET("/getConfig", mw.CheckPermission("system:config:info", check), h.GetConfigHandler)
		config.PUT("/save", mw.CheckPermission("system:config:save", check), h.SaveConfigHandler)
		config.POST("/testMail", mw.CheckPermission("system:config:mail", check), h.SendTestMailHandler)
	}
}
