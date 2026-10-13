package user

import (
	permissioncheck "fortunefinger/backend/internal/biz/sys-check-permission"
	"fortunefinger/backend/internal/server/user/handler"
	userRouter "fortunefinger/backend/internal/server/user/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	authHandler := handler.NewAuthHandler(svcCtx)
	permissionChecker := permissioncheck.NewService(svcCtx)
	userRouter.RegisterAuth(r, authHandler)
	userRouter.RegisterUser(r, handler.NewUserHandler(svcCtx), permissionChecker.CheckPermission)
}
