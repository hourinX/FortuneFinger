package core

import (
	permissioncheck "fortunefinger/backend/internal/biz/sys-check-permission"
	"fortunefinger/backend/internal/server/core/handler"
	coreRouter "fortunefinger/backend/internal/server/core/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	permissionChecker := permissioncheck.NewService(svcCtx)
	coreRouter.RegisterConfig(r, handler.NewConfigHandler(svcCtx), permissionChecker.CheckPermission)
	coreRouter.RegisterMenu(r, handler.NewMenuHandler(svcCtx), permissionChecker.CheckPermission)
	coreRouter.RegisterOrg(r, handler.NewOrgHandler(svcCtx), permissionChecker.CheckPermission)
	coreRouter.RegisterRole(r, handler.NewRoleHandler(svcCtx), permissionChecker.CheckPermission)
	coreRouter.RegisterDict(r, handler.NewDictHandler(svcCtx), permissionChecker.CheckPermission)
}
