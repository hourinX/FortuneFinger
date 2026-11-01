package lottery

import (
	permissioncheck "fortunefinger/backend/internal/biz/sys-check-permission"
	"fortunefinger/backend/internal/server/lottery/handler"
	lotteryRouter "fortunefinger/backend/internal/server/lottery/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	permissionChecker := permissioncheck.NewService(svcCtx)
	lotteryRouter.RegisterActivity(r, handler.NewActivityHandler(svcCtx), permissionChecker.CheckPermission)
	lotteryRouter.RegisterRound(r, handler.NewRoundHandler(svcCtx), permissionChecker.CheckPermission)
}
