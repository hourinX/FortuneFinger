package prize

import (
	permissioncheck "fortunefinger/backend/internal/biz/sys-check-permission"
	"fortunefinger/backend/internal/server/prize/handler"
	prizeRouter "fortunefinger/backend/internal/server/prize/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	permissionChecker := permissioncheck.NewService(svcCtx)
	prizeRouter.RegisterPool(r, handler.NewPoolHandler(svcCtx), permissionChecker.CheckPermission)
}
