package router

import (
	"fortunefinger/backend/internal/server/asset/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAsset(r *gin.Engine, h *handler.AssetHandler) {
	asset := r.Group("/api/asset")
	{
		asset.POST("/upload", h.UploadHandler)
		asset.POST("/generateAvatar", h.GenerateAvatarHandler)
	}
}
