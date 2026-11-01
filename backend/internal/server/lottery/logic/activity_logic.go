package logic

import (
	"context"
	"strconv"
	"time"

	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/sequence"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/lottery/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type ActivityLogic interface {
	GetActivityList(ctx context.Context, req *types.GetActivityListReq) (*types.GetActivityListResp, error)
	GetActivityDetail(ctx context.Context, id uint64) (*types.GetActivityDetailResp, error)
	CreateActivity(ctx context.Context, req *types.CreateActivityReq) error
	UpdateActivity(ctx context.Context, id uint64, req *types.UpdateActivityReq) error
	DeleteActivity(ctx context.Context, id uint64) error
	UpdateActivityCheckinTemplate(ctx context.Context, id uint64, req *types.UpdateActivityCheckinTemplateReq) error
}

type activityLogic struct {
	svcCtx       *svc.ServiceContext
	activityRepo repository.LotteryActivityRepository
	orgRepo      repository.OrgRepository
}

func NewActivityLogic(svcCtx *svc.ServiceContext) ActivityLogic {
	return &activityLogic{
		svcCtx:       svcCtx,
		activityRepo: repository.NewLotteryActivityRepository(svcCtx.GormDB, svcCtx.Snowflake),
		orgRepo:      repository.NewOrgRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (l *activityLogic) GetActivityList(ctx context.Context, req *types.GetActivityListReq) (*types.GetActivityListResp, error) {
	activities, total, err := l.activityRepo.List(ctx, repository.LotteryActivityListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		OrgID:    req.OrgID,
		Date:     req.Date,
		Phase:    req.Phase,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	orgIDs := make([]uint64, 0, len(activities))
	seenOrgIDs := make(map[uint64]struct{}, len(activities))
	for _, activity := range activities {
		if _, exists := seenOrgIDs[activity.OrgID]; exists {
			continue
		}
		seenOrgIDs[activity.OrgID] = struct{}{}
		orgIDs = append(orgIDs, activity.OrgID)
	}

	orgNameByID := make(map[uint64]string, len(orgIDs))
	if len(orgIDs) > 0 {
		orgs, findErr := l.orgRepo.FindByIDs(ctx, orgIDs)
		if findErr != nil {
			return nil, findErr
		}
		for _, org := range orgs {
			orgNameByID[org.ID] = org.Name
		}
	}

	now, err := time.Parse(time.DateTime, time.Now().Format(time.DateTime))
	if err != nil {
		return nil, errorx.Wrap(errorx.InternalServer, "生成当前活动时间失败", err)
	}
	list := make([]types.ActivityListItem, 0, len(activities))
	for _, activity := range activities {
		list = append(list, types.ActivityListItem{
			ID:                activity.ID,
			Code:              activity.Code,
			Name:              activity.Name,
			OrgID:             activity.OrgID,
			OrgName:           orgNameByID[activity.OrgID],
			Venue:             activity.Venue,
			ActivityStartTime: activity.ActivityStartTime.Format(time.DateTime),
			ActivityEndTime:   activity.ActivityEndTime.Format(time.DateTime),
			CheckinStartTime:  activity.CheckinStartTime.Format(time.DateTime),
			CheckinEndTime:    activity.CheckinEndTime.Format(time.DateTime),
			CheckinTemplate:   activity.CheckinTemplate,
			Phase:             activityPhase(activity, now),
			Status:            activity.Status,
		})
	}
	resp := types.GetActivityListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func (l *activityLogic) GetActivityDetail(ctx context.Context, id uint64) (*types.GetActivityDetailResp, error) {
	activity, err := l.activityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errorx.New(errorx.NotFound, "抽奖活动不存在")
	}

	resp := types.GetActivityDetailResp{
		ID:                activity.ID,
		Code:              activity.Code,
		Name:              activity.Name,
		OrgID:             activity.OrgID,
		Venue:             activity.Venue,
		ActivityStartTime: activity.ActivityStartTime.Format(time.DateTime),
		ActivityEndTime:   activity.ActivityEndTime.Format(time.DateTime),
		CheckinStartTime:  activity.CheckinStartTime.Format(time.DateTime),
		CheckinEndTime:    activity.CheckinEndTime.Format(time.DateTime),
		CheckinTemplate:   activity.CheckinTemplate,
		Status:            activity.Status,
		Remark:            activity.Remark,
	}
	return &resp, nil
}

func (l *activityLogic) CreateActivity(ctx context.Context, req *types.CreateActivityReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryActivityCreateDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[ActivityLogic] Redis设置新增抽奖活动防抖失败", zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	activityStartTime, activityEndTime, checkinStartTime, checkinEndTime, err := l.validateActivity(ctx, &req.ActivityReq)
	if err != nil {
		return err
	}
	code, err := sequence.GenerateUniqueShortCode(ctx, l.activityRepo.ExistsByCode)
	if err != nil {
		l.svcCtx.Logger.Error("[ActivityLogic] 生成抽奖活动编码失败", zap.Error(err))
		return errorx.Wrap(errorx.InternalServer, "生成抽奖活动编码失败", err)
	}

	activity := &model.FortuneActivity{
		Code:              code,
		Name:              req.Name,
		OrgID:             req.OrgID,
		Venue:             req.Venue,
		ActivityStartTime: activityStartTime,
		ActivityEndTime:   activityEndTime,
		CheckinStartTime:  checkinStartTime,
		CheckinEndTime:    checkinEndTime,
		CheckinTemplate:   req.CheckinTemplate,
		Status:            req.Status,
		Remark:            req.Remark,
		CreatedBy:         operatorID,
		UpdatedBy:         operatorID,
	}
	if _, err := l.activityRepo.Create(ctx, activity); err != nil {
		l.svcCtx.Logger.Error("[ActivityLogic] 新增抽奖活动失败", zap.String("activity_name", req.Name), zap.Error(err))
		return err
	}
	return nil
}

func (l *activityLogic) UpdateActivity(ctx context.Context, id uint64, req *types.UpdateActivityReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryActivityUpdateDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[ActivityLogic] Redis设置修改抽奖活动防抖失败", zap.Uint64("activity_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	activity, err := l.activityRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if activity == nil {
		return errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	activityStartTime, activityEndTime, checkinStartTime, checkinEndTime, err := l.validateActivity(ctx, &req.ActivityReq)
	if err != nil {
		return err
	}

	updates := map[string]any{
		"name":                req.Name,
		"org_id":              req.OrgID,
		"venue":               req.Venue,
		"activity_start_time": activityStartTime,
		"activity_end_time":   activityEndTime,
		"checkin_start_time":  checkinStartTime,
		"checkin_end_time":    checkinEndTime,
		"checkin_template":    req.CheckinTemplate,
		"status":              req.Status,
		"updated_by":          operatorID,
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if err := l.activityRepo.Update(ctx, id, updates); err != nil {
		appErr, ok := errorx.As(err)
		if !ok || appErr.Code() == errorx.DatabaseError || appErr.Code() == errorx.InternalServer {
			l.svcCtx.Logger.Error("[ActivityLogic] 修改抽奖活动失败", zap.Uint64("activity_id", id), zap.Error(err))
		}
		return err
	}
	return nil
}

func (l *activityLogic) DeleteActivity(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.LotteryActivityDeleteDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[ActivityLogic] Redis设置删除抽奖活动防抖失败", zap.Uint64("activity_id", id), zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	if err := l.activityRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		appErr, ok := errorx.As(err)
		if !ok || appErr.Code() == errorx.DatabaseError || appErr.Code() == errorx.InternalServer {
			l.svcCtx.Logger.Error("[ActivityLogic] 删除抽奖活动失败", zap.Uint64("activity_id", id), zap.Error(err))
		}
		return err
	}
	return nil
}

func (l *activityLogic) UpdateActivityCheckinTemplate(ctx context.Context, id uint64, req *types.UpdateActivityCheckinTemplateReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	activity, err := l.activityRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if activity == nil {
		return errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	if err := l.activityRepo.Update(ctx, id, map[string]any{
		"checkin_template": req.CheckinTemplate,
		"updated_by":       operatorID,
	}); err != nil {
		appErr, ok := errorx.As(err)
		if !ok || appErr.Code() == errorx.DatabaseError || appErr.Code() == errorx.InternalServer {
			l.svcCtx.Logger.Error("[ActivityLogic] 修改抽奖活动签到模板失败", zap.Uint64("activity_id", id), zap.Error(err))
		}
		return err
	}
	return nil
}

func (l *activityLogic) validateActivity(ctx context.Context, req *types.ActivityReq) (time.Time, time.Time, time.Time, time.Time, error) {
	activityStartTime, err := time.Parse("2006-01-02 15:04", req.ActivityStartTime[:16])
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "活动开始时间无效")
	}
	activityEndTime, err := time.Parse("2006-01-02 15:04", req.ActivityEndTime[:16])
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "活动结束时间无效")
	}
	checkinStartTime, err := time.Parse("2006-01-02 15:04", req.CheckinStartTime[:16])
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "签到开始时间无效")
	}
	checkinEndTime, err := time.Parse("2006-01-02 15:04", req.CheckinEndTime[:16])
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "签到结束时间无效")
	}
	activityStartTimestamp := activityStartTime.UnixMilli()
	activityEndTimestamp := activityEndTime.UnixMilli()
	checkinStartTimestamp := checkinStartTime.UnixMilli()
	checkinEndTimestamp := checkinEndTime.UnixMilli()
	if activityEndTimestamp <= activityStartTimestamp {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "活动结束时间必须晚于开始时间")
	}
	if checkinEndTimestamp <= checkinStartTimestamp {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "签到结束时间必须晚于开始时间")
	}
	if checkinEndTimestamp > activityEndTimestamp {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "签到结束时间不能晚于活动结束时间")
	}
	if checkinEndTimestamp-checkinStartTimestamp > (7 * 24 * time.Hour).Milliseconds() {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "签到时间范围不能超过七天")
	}

	org, err := l.orgRepo.FindByID(ctx, req.OrgID)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, err
	}
	if org == nil || (org.Types != 1 && org.Types != 2) {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.InvalidParameter, "选择的归属组织不存在")
	}
	if org.Status != 1 {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, errorx.New(errorx.Conflict, "选择的归属组织已被禁用")
	}
	return activityStartTime, activityEndTime, checkinStartTime, checkinEndTime, nil
}

func activityPhase(activity model.FortuneActivity, now time.Time) string {
	nowTimestamp := now.UnixMilli()
	if nowTimestamp < activity.ActivityStartTime.UnixMilli() {
		return "pending"
	}
	if nowTimestamp > activity.ActivityEndTime.UnixMilli() {
		return "finished"
	}
	return "active"
}
