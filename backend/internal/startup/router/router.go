package router

import (
	"fortunefinger/backend/internal/mw"
	"fortunefinger/backend/internal/pkg/requestctx"
	"fortunefinger/backend/internal/server/asset"
	"fortunefinger/backend/internal/server/checkin"
	"fortunefinger/backend/internal/server/core"
	"fortunefinger/backend/internal/server/prize"
	"fortunefinger/backend/internal/server/user"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func New(svcCtx *svc.ServiceContext) *gin.Engine {
	r := gin.New()
	r.Use(mw.Recovery(svcCtx.Logger))
	r.Use(mw.JWT(svcCtx.JWT, svcCtx.Logger))
	r.Use(requestctx.Middleware())

	asset.Register(r, svcCtx)
	checkin.Register(r, svcCtx)
	core.Register(r, svcCtx)
	prize.Register(r, svcCtx)
	user.Register(r, svcCtx)
	return r
}
