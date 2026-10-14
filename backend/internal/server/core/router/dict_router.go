package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterDict(r *gin.Engine, h *handler.DictHandler, check mw.Checker) {
	dict := r.Group("/api/dict")
	{
		dict.GET("/type/list", mw.CheckPermission("system:dict:list-type", check), h.GetDictTypeListHandler)
		dict.GET("/type/detail/:id", mw.CheckPermission("system:dict:info-type", check), h.GetDictTypeDetailHandler)
		dict.POST("/type/create", mw.CheckPermission("system:dict:add-type", check), h.CreateDictTypeHandler)
		dict.PUT("/type/update/:id", mw.CheckPermission("system:dict:edit-type", check), h.UpdateDictTypeHandler)
		dict.DELETE("/type/:id", mw.CheckPermission("system:dict:del-type", check), h.DeleteDictTypeHandler)

		dict.GET("/data/list", mw.CheckPermission("system:dict:list-data", check), h.GetDictDataListHandler)
		dict.GET("/data/detail/:id", mw.CheckPermission("system:dict:info-data", check), h.GetDictDataDetailHandler)
		dict.POST("/data/create", mw.CheckPermission("system:dict:add-data", check), h.CreateDictDataHandler)
		dict.PUT("/data/update/:id", mw.CheckPermission("system:dict:edit-data", check), h.UpdateDictDataHandler)
		dict.DELETE("/data/:id", mw.CheckPermission("system:dict:del-data", check), h.DeleteDictDataHandler)

		dict.GET("/options", h.GetDictOptionsHandler)
	}
}
