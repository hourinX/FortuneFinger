package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/prize/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPool(r *gin.Engine, h *handler.PoolHandler, check mw.Checker) {
	pool := r.Group("/api/prize/pool")
	{
		pool.GET("/list", mw.CheckPermission("prize:pool:list", check), h.GetPoolListHandler)
		pool.GET("/detail/:id", mw.CheckPermission("prize:pool:info", check), h.GetPoolDetailHandler)
		pool.POST("/create", mw.CheckPermission("prize:pool:add", check), h.CreatePoolHandler)
		pool.PUT("/update/:id", mw.CheckPermission("prize:pool:edit", check), h.UpdatePoolHandler)
		pool.DELETE("/:id", mw.CheckPermission("prize:pool:del", check), h.DeletePoolHandler)
	}
}
