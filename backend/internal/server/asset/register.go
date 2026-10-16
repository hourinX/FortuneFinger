package asset

import (
	"fortunefinger/backend/internal/server/asset/handler"
	assetRouter "fortunefinger/backend/internal/server/asset/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	assetRouter.RegisterAsset(r, handler.NewAssetHandler(svcCtx))
}
