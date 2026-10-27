package checkin

import (
	"fortunefinger/backend/internal/server/checkin/handler"
	checkinRouter "fortunefinger/backend/internal/server/checkin/router"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	checkinRouter.RegisterCheckin(r, handler.NewCheckinHandler(svcCtx))
}
