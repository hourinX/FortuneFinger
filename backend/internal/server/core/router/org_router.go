package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/server/core/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrg(r *gin.Engine, h *handler.OrgHandler, check mw.Checker) {
	org := r.Group("/api/org")
	{
		org.GET("/list", mw.CheckPermission("perm:org:list", check), h.GetOrgListHandler)
		org.GET("/detail/:id", mw.CheckPermission("perm:org:info", check), h.GetOrgDetailHandler)
		org.POST("/create", mw.CheckPermission("perm:org:add", check), h.CreateOrgHandler)
		org.PUT("/update/:id", mw.CheckPermission("perm:org:edit", check), h.UpdateOrgHandler)
		org.DELETE("/:id", mw.CheckPermission("perm:org:del", check), h.DeleteOrgHandler)
		org.GET("/options", h.GetOrgOptionsHandler)
		org.GET("/companies", h.GetOrgCompaniesHandler)
	}
}
