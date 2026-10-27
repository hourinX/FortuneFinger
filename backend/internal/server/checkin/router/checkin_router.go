package router

import (
	"fortunefinger/backend/internal/server/checkin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCheckin(r *gin.Engine, h *handler.CheckinHandler) {
	checkin := r.Group("/api/checkin")
	{
		checkin.GET("/template", h.GetCheckinTemplateHandler)
	}
}
