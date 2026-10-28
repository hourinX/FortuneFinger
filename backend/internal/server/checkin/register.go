package checkin

import (
	"context"
	"encoding/json"
	"time"

	"fortunefinger/backend/internal/pkg/ws"
	"fortunefinger/backend/internal/server/checkin/handler"
	"fortunefinger/backend/internal/server/checkin/logic"
	checkinRouter "fortunefinger/backend/internal/server/checkin/router"
	"fortunefinger/backend/internal/server/checkin/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svcCtx *svc.ServiceContext) {
	checkinRouter.RegisterCheckin(r, handler.NewCheckinHandler(svcCtx))
}

func RegisterSocket(manager *ws.Manager, svcCtx *svc.ServiceContext) {
	checkinLogic := logic.NewCheckinLogic(svcCtx)
	manager.MustRegister(ws.Topic{
		Key:      "checkin.state",
		Interval: 2 * time.Second,
		Mode:     ws.PushOnChange,
		Load: func(ctx context.Context) (any, error) {
			return checkinLogic.GetCheckinState(ctx, &types.GetCheckinStateReq{})
		},
		Fingerprint: func(data any) ([]byte, error) {
			state, ok := data.(*types.GetCheckinStateResp)
			if !ok {
				return json.Marshal(data)
			}
			return json.Marshal(struct {
				Status           string  `json:"status"`
				CheckinStartTime string  `json:"checkin_start_time"`
				CheckinEndTime   string  `json:"checkin_end_time"`
				CheckedIn        uint64  `json:"checked_in"`
				Total            *uint64 `json:"total"`
			}{
				Status:           state.Status,
				CheckinStartTime: state.CheckinStartTime,
				CheckinEndTime:   state.CheckinEndTime,
				CheckedIn:        state.CheckedIn,
				Total:            state.Total,
			})
		},
	})
}
