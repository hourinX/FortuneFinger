package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/lottery/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRound(r *gin.Engine, h *handler.RoundHandler, check mw.Checker) {
	round := r.Group("/api/lottery/round")
	{
		round.GET("/list", mw.CheckPermission("lottery:round:list", check), h.GetRoundListHandler)
		round.GET("/detail/:id", mw.CheckPermission("lottery:round:info", check), h.GetRoundDetailHandler)
		round.POST("/create", mw.CheckPermission("lottery:round:add", check), h.CreateRoundHandler)
		round.PUT("/update/:id", mw.CheckPermission("lottery:round:edit", check), h.UpdateRoundHandler)
		round.DELETE("/:id", mw.CheckPermission("lottery:round:del", check), h.DeleteRoundHandler)
	}
}
