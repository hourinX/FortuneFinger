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
	"fortunefinger/backend/internal/server/prize/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type PoolLogic interface {
	GetPoolList(ctx context.Context, req *types.GetPoolListReq) (*types.GetPoolListResp, error)
	GetPoolDetail(ctx context.Context, id uint64) (*types.GetPoolDetailResp, error)
	CreatePool(ctx context.Context, req *types.CreatePoolReq) error
	UpdatePool(ctx context.Context, id uint64, req *types.UpdatePoolReq) error
	DeletePool(ctx context.Context, id uint64) error
}

type poolLogic struct {
	svcCtx       *svc.ServiceContext
	prizeRepo    repository.PrizeRepository
	storageRepo  repository.StorageRepository
	orgRepo      repository.OrgRepository
	dictDataRepo repository.DictDataRepository
}

func NewPoolLogic(svcCtx *svc.ServiceContext) PoolLogic {
	return &poolLogic{
		svcCtx:       svcCtx,
		prizeRepo:    repository.NewPrizeRepository(svcCtx.GormDB, svcCtx.Snowflake),
		storageRepo:  repository.NewStorageRepository(svcCtx.GormDB, svcCtx.Snowflake),
		orgRepo:      repository.NewOrgRepository(svcCtx.GormDB, svcCtx.Snowflake),
		dictDataRepo: repository.NewDictDataRepository(svcCtx.GormDB),
	}
}

func (l *poolLogic) GetPoolList(ctx context.Context, req *types.GetPoolListReq) (*types.GetPoolListResp, error) {
	prizes, total, err := l.prizeRepo.List(ctx, repository.PrizeListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Level:    req.Level,
		Types:    req.Types,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	imageURLs, err := l.resolvePrizeImageURLs(ctx, prizes)
	if err != nil {
		return nil, err
	}

	orgIDs := make([]uint64, 0, len(prizes))
	seenOrgIDs := make(map[uint64]struct{}, len(prizes))
	for _, prize := range prizes {
		if prize.OrgID == 0 {
			continue
		}
		if _, exists := seenOrgIDs[prize.OrgID]; exists {
			continue
		}
		seenOrgIDs[prize.OrgID] = struct{}{}
		orgIDs = append(orgIDs, prize.OrgID)
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

	list := make([]types.PoolListItem, 0, len(prizes))
	for _, prize := range prizes {
		list = append(list, types.PoolListItem{
			ID:       prize.ID,
			Name:     prize.Name,
			ImageURL: imageURLs[prize.ID],
			OrgName:  orgNameByID[prize.OrgID],
			Level:    prize.Level,
			Types:    prize.Types,
			Quantity: prize.Quantity,
			Status:   prize.Status,
		})
	}
	resp := types.GetPoolListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func (l *poolLogic) GetPoolDetail(ctx context.Context, id uint64) (*types.GetPoolDetailResp, error) {
	prize, err := l.prizeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if prize == nil {
		return nil, errorx.New(errorx.NotFound, "奖品不存在")
	}

	imageURLs, err := l.resolvePrizeImageURLs(ctx, []model.FortunePrize{*prize})
	if err != nil {
		return nil, err
	}

	resp := types.GetPoolDetailResp{
		ID:        prize.ID,
		Name:      prize.Name,
		StorageID: prize.StorageID,
		ImageURL:  imageURLs[prize.ID],
		OrgID:     prize.OrgID,
		Level:     prize.Level,
		Types:     prize.Types,
		Quantity:  prize.Quantity,
		Status:    prize.Status,
		Sort:      prize.Sort,
		Remark:    prize.Remark,
	}
	return &resp, nil
}

func (l *poolLogic) CreatePool(ctx context.Context, req *types.CreatePoolReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.PrizeCreateDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[PoolLogic] Redis设置新增奖品防抖失败", zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	if err := l.validatePool(ctx, 0, &req.PoolReq); err != nil {
		return err
	}

	prize := &model.FortunePrize{
		Name:      req.Name,
		StorageID: req.StorageID,
		OrgID:     req.OrgID,
		Level:     req.Level,
		Types:     req.Types,
		Quantity:  req.Quantity,
		Status:    req.Status,
		Sort:      req.Sort,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if _, err := l.prizeRepo.Create(ctx, prize); err != nil {
		l.svcCtx.Logger.Error("[PoolLogic] 新增奖品失败", zap.String("prize_name", req.Name), zap.Error(err))
		return err
	}
	return nil
}

func (l *poolLogic) UpdatePool(ctx context.Context, id uint64, req *types.UpdatePoolReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.PrizeUpdateDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[PoolLogic] Redis设置修改奖品防抖失败", zap.Uint64("prize_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	prize, err := l.prizeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if prize == nil {
		return errorx.New(errorx.NotFound, "奖品不存在")
	}

	if err := l.validatePool(ctx, id, &req.PoolReq); err != nil {
		return err
	}

	updates := map[string]any{
		"name":       req.Name,
		"storage_id": req.StorageID,
		"org_id":     req.OrgID,
		"level":      req.Level,
		"types":      req.Types,
		"quantity":   req.Quantity,
		"status":     req.Status,
		"sort":       req.Sort,
		"remark":     req.Remark,
		"updated_by": operatorID,
	}
	if err := l.prizeRepo.Update(ctx, id, updates); err != nil {
		appErr, ok := errorx.As(err)
		if !ok || appErr.Code() == errorx.DatabaseError || appErr.Code() == errorx.InternalServer {
			l.svcCtx.Logger.Error("[PoolLogic] 修改奖品失败", zap.Uint64("prize_id", id), zap.Error(err))
		}
		return err
	}
	return nil
}

func (l *poolLogic) DeletePool(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.PrizeDeleteDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[PoolLogic] Redis设置删除奖品防抖失败", zap.Uint64("prize_id", id), zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	if err := l.prizeRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		appErr, ok := errorx.As(err)
		if !ok || appErr.Code() == errorx.DatabaseError || appErr.Code() == errorx.InternalServer {
			l.svcCtx.Logger.Error("[PoolLogic] 删除奖品失败", zap.Uint64("prize_id", id), zap.Error(err))
		}
		return err
	}
	return nil
}

func (l *poolLogic) validatePool(ctx context.Context, id uint64, req *types.PoolReq) error {
	org, err := l.orgRepo.FindByID(ctx, req.OrgID)
	if err != nil {
		return err
	}
	if org == nil || org.ParentID == 0 || org.Types != 2 {
		return errorx.New(errorx.InvalidParameter, "选择的归属企业不存在")
	}
	if org.Status != 1 {
		return errorx.New(errorx.Conflict, "选择的归属企业已被禁用")
	}

	level, err := l.dictDataRepo.FindByCodeAndValue(ctx, "prize_level", req.Level)
	if err != nil {
		return err
	}
	if level == nil || level.Status != 1 {
		return errorx.New(errorx.InvalidParameter, "选择的奖项等级不存在或已被禁用")
	}

	existingPrize, err := l.prizeRepo.FindByOrgIDAndLevel(ctx, req.OrgID, req.Level)
	if err != nil {
		return err
	}
	if existingPrize != nil && existingPrize.ID != id {
		return errorx.New(errorx.Conflict, "当前企业已配置该奖项等级")
	}

	prizeType, err := l.dictDataRepo.FindByCodeAndValue(ctx, "prize_types", strconv.FormatUint(uint64(req.Types), 10))
	if err != nil {
		return err
	}
	if prizeType == nil || prizeType.Status != 1 {
		return errorx.New(errorx.InvalidParameter, "选择的奖品类型不存在或已被禁用")
	}

	if req.StorageID == 0 {
		return nil
	}
	storage, err := l.storageRepo.FindByID(ctx, req.StorageID)
	if err != nil {
		return err
	}
	if storage == nil || storage.Scene != "prize" {
		return errorx.New(errorx.InvalidParameter, "选择的奖品图片不存在")
	}
	return nil
}

func (l *poolLogic) resolvePrizeImageURLs(ctx context.Context, prizes []model.FortunePrize) (map[uint64]string, error) {
	urls := make(map[uint64]string, len(prizes))
	storageIDs := make([]uint64, 0, len(prizes))
	storageIDByPrizeID := make(map[uint64]uint64, len(prizes))

	for _, prize := range prizes {
		if prize.StorageID == 0 {
			continue
		}
		storageIDs = append(storageIDs, prize.StorageID)
		storageIDByPrizeID[prize.ID] = prize.StorageID
	}

	storages, err := l.storageRepo.FindByIDs(ctx, storageIDs)
	if err != nil {
		return nil, err
	}

	storageByID := make(map[uint64]model.FortuneStorage, len(storages))
	for _, storageRecord := range storages {
		storageByID[storageRecord.ID] = storageRecord
	}

	for prizeID, storageID := range storageIDByPrizeID {
		storageRecord, exists := storageByID[storageID]
		if !exists || storageRecord.Scene != "prize" {
			continue
		}

		url, signErr := l.svcCtx.Storage.PresignGetObject(ctx, storageRecord.Path, time.Hour)
		if signErr != nil {
			l.svcCtx.Logger.Error("[PoolLogic] 生成奖品图片访问地址失败",
				zap.Uint64("prize_id", prizeID),
				zap.Uint64("storage_id", storageID),
				zap.Error(signErr),
			)
			return nil, errorx.Wrap(errorx.ServiceUnavailable, "加载奖品图片失败", signErr)
		}
		urls[prizeID] = url
	}
	return urls, nil
}
