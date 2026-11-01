package logic

import (
	"context"
	"strconv"
	"time"

	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/lottery/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type RoundLogic interface {
	GetRoundList(ctx context.Context, req *types.GetRoundListReq) (*types.GetRoundListResp, error)
	GetRoundDetail(ctx context.Context, id uint64) (*types.GetRoundDetailResp, error)
	CreateRound(ctx context.Context, req *types.CreateRoundReq) error
	UpdateRound(ctx context.Context, id uint64, req *types.UpdateRoundReq) error
	DeleteRound(ctx context.Context, id uint64) error
}

type roundLogic struct {
	svcCtx       *svc.ServiceContext
	roundRepo    repository.RoundRepository
	activityRepo repository.LotteryActivityRepository
	prizeRepo    repository.PrizeRepository
}

func NewRoundLogic(svcCtx *svc.ServiceContext) RoundLogic {
	return &roundLogic{
		svcCtx:       svcCtx,
		roundRepo:    repository.NewRoundRepository(svcCtx.GormDB, svcCtx.Snowflake),
		activityRepo: repository.NewLotteryActivityRepository(svcCtx.GormDB, svcCtx.Snowflake),
		prizeRepo:    repository.NewPrizeRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (l *roundLogic) GetRoundList(ctx context.Context, req *types.GetRoundListReq) (*types.GetRoundListResp, error) {
	activities, err := l.activityRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	activityIDs := make([]uint64, 0, len(activities))
	for _, activity := range activities {
		if req.ActivityID == 0 || activity.ID == req.ActivityID {
			activityIDs = append(activityIDs, activity.ID)
		}
	}
	rounds, err := l.roundRepo.ListByActivityIDs(ctx, activityIDs)
	if err != nil {
		return nil, err
	}
	prizeIDs := make([]uint64, 0, len(rounds))
	seenPrizeIDs := make(map[uint64]struct{}, len(rounds))
	for _, round := range rounds {
		if _, exists := seenPrizeIDs[round.PrizeID]; exists {
			continue
		}
		seenPrizeIDs[round.PrizeID] = struct{}{}
		prizeIDs = append(prizeIDs, round.PrizeID)
	}
	prizes, err := l.prizeRepo.FindByIDs(ctx, prizeIDs)
	if err != nil {
		return nil, err
	}
	prizeByID := make(map[uint64]model.FortunePrize, len(prizes))
	for _, prize := range prizes {
		prizeByID[prize.ID] = prize
	}
	roundsByActivityID := make(map[uint64][]model.FortuneRound, len(activityIDs))
	for _, round := range rounds {
		roundsByActivityID[round.ActivityID] = append(roundsByActivityID[round.ActivityID], round)
	}

	now := time.Now()
	list := make([]types.RoundActivityItem, 0, len(activities))
	for _, activity := range activities {
		if req.ActivityID != 0 && activity.ID != req.ActivityID {
			continue
		}
		activityRounds := roundsByActivityID[activity.ID]
		if len(activityRounds) == 0 {
			continue
		}
		dateItems := make([]types.RoundDateItem, 0)
		for index, round := range activityRounds {
			prize := prizeByID[round.PrizeID]
			drawDate := round.DrawDate.Format(time.DateOnly)
			phase := roundPhase(round, activity, now)
			if req.Phase != "" && phase != req.Phase {
				continue
			}
			if len(dateItems) == 0 || dateItems[len(dateItems)-1].DrawDate != drawDate {
				dateItems = append(dateItems, types.RoundDateItem{
					DrawDate: drawDate,
					Rounds:   make([]types.RoundListItem, 0),
				})
			}
			dateIndex := len(dateItems) - 1
			dateItems[dateIndex].Rounds = append(dateItems[dateIndex].Rounds, types.RoundListItem{
				ID:          round.ID,
				Sequence:    index + 1,
				PrizeID:     round.PrizeID,
				PrizeName:   prize.Name,
				Level:       prize.Level,
				Name:        round.Name,
				DrawTime:    round.DrawTime,
				WinnerCount: round.WinnerCount,
				Phase:       phase,
			})
		}
		if len(dateItems) == 0 {
			continue
		}
		roundCount := 0
		for _, dateItem := range dateItems {
			roundCount += len(dateItem.Rounds)
		}
		list = append(list, types.RoundActivityItem{
			ActivityID:        activity.ID,
			ActivityName:      activity.Name,
			ActivityStartTime: activity.ActivityStartTime.Format(time.DateTime),
			ActivityEndTime:   activity.ActivityEndTime.Format(time.DateTime),
			RoundCount:        roundCount,
			Dates:             dateItems,
		})
	}

	total := int64(0)
	for _, activity := range list {
		total += int64(activity.RoundCount)
	}
	page, pageSize := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	list = paginateRoundActivities(list, (page-1)*pageSize, pageSize)
	resp := types.GetRoundListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func paginateRoundActivities(list []types.RoundActivityItem, offset, limit int) []types.RoundActivityItem {
	result := make([]types.RoundActivityItem, 0)
	position := 0
	remaining := limit
	for _, activity := range list {
		if remaining == 0 {
			break
		}
		activityItem := activity
		activityItem.Dates = make([]types.RoundDateItem, 0)
		for _, date := range activity.Dates {
			if remaining == 0 {
				break
			}
			dateItem := date
			dateItem.Rounds = make([]types.RoundListItem, 0)
			for _, round := range date.Rounds {
				if position >= offset && remaining > 0 {
					dateItem.Rounds = append(dateItem.Rounds, round)
					remaining--
				}
				position++
			}
			if len(dateItem.Rounds) > 0 {
				activityItem.Dates = append(activityItem.Dates, dateItem)
			}
		}
		if len(activityItem.Dates) > 0 {
			result = append(result, activityItem)
		}
	}
	return result
}

func (l *roundLogic) GetRoundDetail(ctx context.Context, id uint64) (*types.GetRoundDetailResp, error) {
	round, err := l.roundRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if round == nil {
		return nil, errorx.New(errorx.NotFound, "抽奖轮次不存在")
	}

	resp := types.GetRoundDetailResp{
		ID:          round.ID,
		ActivityID:  round.ActivityID,
		PrizeID:     round.PrizeID,
		Name:        round.Name,
		DrawDate:    round.DrawDate.Format(time.DateOnly),
		DrawTime:    round.DrawTime,
		WinnerCount: round.WinnerCount,
	}
	return &resp, nil
}

func (l *roundLogic) CreateRound(ctx context.Context, req *types.CreateRoundReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryRoundCreateDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] Redis设置新增抽奖轮次防抖失败", zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		l.svcCtx.Logger.Warn("[RoundLogic] 新增抽奖轮次操作过于频繁", zap.Uint64("operator_id", operatorID))
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	drawDate, err := l.validateRound(ctx, &req.RoundReq)
	if err != nil {
		return err
	}
	round := &model.FortuneRound{
		ActivityID:  req.ActivityID,
		PrizeID:     req.PrizeID,
		Name:        req.Name,
		DrawDate:    drawDate,
		DrawTime:    req.DrawTime,
		WinnerCount: req.WinnerCount,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}
	if _, err := l.roundRepo.Create(ctx, round); err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] 新增抽奖轮次失败", zap.Uint64("activity_id", req.ActivityID), zap.Error(err))
		return err
	}
	return nil
}

func (l *roundLogic) UpdateRound(ctx context.Context, id uint64, req *types.UpdateRoundReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryRoundUpdateDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] Redis设置修改抽奖轮次防抖失败", zap.Uint64("round_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		l.svcCtx.Logger.Warn("[RoundLogic] 修改抽奖轮次操作过于频繁", zap.Uint64("round_id", id), zap.Uint64("operator_id", operatorID))
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	round, err := l.roundRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if round == nil {
		return errorx.New(errorx.NotFound, "抽奖轮次不存在")
	}
	if !roundPending(*round, time.Now()) {
		return errorx.New(errorx.Conflict, "已到开奖时间的轮次不能修改")
	}
	drawDate, err := l.validateRound(ctx, &req.RoundReq)
	if err != nil {
		return err
	}
	if err := l.roundRepo.Update(ctx, id, map[string]any{
		"activity_id":  req.ActivityID,
		"prize_id":     req.PrizeID,
		"name":         req.Name,
		"draw_date":    drawDate,
		"draw_time":    req.DrawTime,
		"winner_count": req.WinnerCount,
		"updated_by":   operatorID,
	}); err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] 修改抽奖轮次失败", zap.Uint64("round_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *roundLogic) DeleteRound(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryRoundDeleteDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] Redis设置删除抽奖轮次防抖失败", zap.Uint64("round_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		l.svcCtx.Logger.Warn("[RoundLogic] 删除抽奖轮次操作过于频繁", zap.Uint64("round_id", id), zap.Uint64("operator_id", operatorID))
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	round, err := l.roundRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if round == nil {
		return errorx.New(errorx.NotFound, "抽奖轮次不存在")
	}
	if !roundPending(*round, time.Now()) {
		return errorx.New(errorx.Conflict, "已到开奖时间的轮次不能删除")
	}
	if err := l.roundRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		l.svcCtx.Logger.Error("[RoundLogic] 删除抽奖轮次失败", zap.Uint64("round_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *roundLogic) validateRound(ctx context.Context, req *types.RoundReq) (time.Time, error) {
	drawDate, err := time.Parse(time.DateOnly, req.DrawDate)
	if err != nil {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "开奖日期无效")
	}
	drawClock, err := time.Parse(time.TimeOnly, req.DrawTime)
	if err != nil {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "开奖时间无效")
	}
	drawAt := time.Date(
		drawDate.Year(),
		drawDate.Month(),
		drawDate.Day(),
		drawClock.Hour(),
		drawClock.Minute(),
		drawClock.Second(),
		0,
		time.UTC,
	)

	activity, err := l.activityRepo.FindByID(ctx, req.ActivityID)
	if err != nil {
		return time.Time{}, err
	}
	if activity == nil {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "所属活动不存在")
	}
	if activity.Status != 1 {
		return time.Time{}, errorx.New(errorx.Conflict, "所属活动已被禁用")
	}
	if drawAt.Before(activity.ActivityStartTime) || drawAt.After(activity.ActivityEndTime) {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "开奖时间必须处于活动时间范围内")
	}

	prize, err := l.prizeRepo.FindByID(ctx, req.PrizeID)
	if err != nil {
		return time.Time{}, err
	}
	if prize == nil {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "选择的奖品不存在")
	}
	if prize.Status != 1 {
		return time.Time{}, errorx.New(errorx.Conflict, "选择的奖品已被禁用")
	}
	if prize.OrgID != activity.OrgID {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "奖品与活动不属于同一企业")
	}
	if req.WinnerCount > prize.Quantity {
		return time.Time{}, errorx.New(errorx.InvalidParameter, "中奖人数不能超过奖品数量")
	}
	return drawDate, nil
}

func roundPending(round model.FortuneRound, now time.Time) bool {
	drawClock, err := time.Parse(time.TimeOnly, round.DrawTime)
	if err != nil {
		return false
	}
	drawAt := time.Date(
		round.DrawDate.Year(),
		round.DrawDate.Month(),
		round.DrawDate.Day(),
		drawClock.Hour(),
		drawClock.Minute(),
		drawClock.Second(),
		0,
		round.DrawDate.Location(),
	)
	return now.Before(drawAt)
}

func roundPhase(round model.FortuneRound, activity model.FortuneActivity, now time.Time) string {
	drawClock, err := time.Parse(time.TimeOnly, round.DrawTime)
	if err != nil {
		return "pending"
	}
	drawAt := time.Date(
		round.DrawDate.Year(),
		round.DrawDate.Month(),
		round.DrawDate.Day(),
		drawClock.Hour(),
		drawClock.Minute(),
		drawClock.Second(),
		0,
		round.DrawDate.Location(),
	)
	if now.Before(drawAt) {
		return "pending"
	}
	if now.After(activity.ActivityEndTime) {
		return "finished"
	}
	return "active"
}
