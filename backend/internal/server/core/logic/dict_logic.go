package logic

import (
	"context"
	"strings"
	"time"

	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const maxDictOptionCodes = 50

type DictLogic interface {
	GetDictTypeList(ctx context.Context, req *types.GetDictTypeListReq) (*types.GetDictTypeListResp, error)
	GetDictTypeDetail(ctx context.Context, id uint64) (*types.GetDictTypeDetailResp, error)
	CreateDictType(ctx context.Context, req *types.CreateDictTypeReq) error
	UpdateDictType(ctx context.Context, id uint64, req *types.UpdateDictTypeReq) error
	DeleteDictType(ctx context.Context, id uint64) error

	GetDictDataList(ctx context.Context, req *types.GetDictDataListReq) (*types.GetDictDataListResp, error)
	GetDictDataDetail(ctx context.Context, id uint64) (*types.GetDictDataDetailResp, error)
	CreateDictData(ctx context.Context, req *types.CreateDictDataReq) error
	UpdateDictData(ctx context.Context, id uint64, req *types.UpdateDictDataReq) error
	DeleteDictData(ctx context.Context, id uint64) error

	GetDictOptions(ctx context.Context, req *types.GetDictOptionsReq) (*types.GetDictOptionsResp, error)
}

type dictLogic struct {
	svcCtx       *svc.ServiceContext
	dictTypeRepo repository.DictTypeRepository
	dictDataRepo repository.DictDataRepository
}

func NewDictLogic(svcCtx *svc.ServiceContext) DictLogic {
	return &dictLogic{
		svcCtx:       svcCtx,
		dictTypeRepo: repository.NewDictTypeRepository(svcCtx.GormDB, svcCtx.Snowflake),
		dictDataRepo: repository.NewDictDataRepository(svcCtx.GormDB),
	}
}

func (l *dictLogic) GetDictTypeList(ctx context.Context, req *types.GetDictTypeListReq) (*types.GetDictTypeListResp, error) {
	dictTypes, total, err := l.dictTypeRepo.List(ctx, repository.DictTypeListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Code:     req.Code,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.DictTypeListItem, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		list = append(list, types.DictTypeListItem{
			ID:       dictType.ID,
			Name:     dictType.Name,
			Code:     dictType.Code,
			IsSystem: dictType.IsSystem,
			Status:   dictType.Status,
		})
	}
	resp := types.GetDictTypeListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func (l *dictLogic) GetDictTypeDetail(ctx context.Context, id uint64) (*types.GetDictTypeDetailResp, error) {
	dictType, err := l.dictTypeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dictType == nil {
		return nil, errorx.New(errorx.NotFound, "字典类型不存在")
	}
	resp := types.GetDictTypeDetailResp{
		ID:       dictType.ID,
		Name:     dictType.Name,
		Code:     dictType.Code,
		IsSystem: dictType.IsSystem,
		Status:   dictType.Status,
		Remark:   dictType.Remark,
	}
	return &resp, nil
}

func (l *dictLogic) CreateDictType(ctx context.Context, req *types.CreateDictTypeReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictTypeCreateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置新增字典类型防抖失败", zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	existing, err := l.dictTypeRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errorx.New(errorx.Conflict, "字典编码已存在")
	}
	dictType := &model.FortuneDictType{
		Name:      req.Name,
		Code:      req.Code,
		IsSystem:  2,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if _, err := l.dictTypeRepo.Create(ctx, dictType); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 新增字典类型失败", zap.String("dict_code", req.Code), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) UpdateDictType(ctx context.Context, id uint64, req *types.UpdateDictTypeReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictTypeUpdateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置修改字典类型防抖失败", zap.Uint64("dict_type_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	dictType, err := l.dictTypeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if dictType == nil {
		return errorx.New(errorx.NotFound, "字典类型不存在")
	}
	updates := map[string]any{
		"name":       req.Name,
		"status":     req.Status,
		"remark":     req.Remark,
		"updated_by": operatorID,
	}
	if err := l.dictTypeRepo.Update(ctx, id, updates); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 修改字典类型失败", zap.Uint64("dict_type_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) DeleteDictType(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictTypeDeleteDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置删除字典类型防抖失败", zap.Uint64("dict_type_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	dictType, err := l.dictTypeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if dictType == nil {
		return errorx.New(errorx.NotFound, "字典类型不存在")
	}
	if dictType.IsSystem == 1 {
		return errorx.New(errorx.Forbidden, "系统字典不允许删除")
	}
	hasData, err := l.dictDataRepo.HasByCode(ctx, dictType.Code)
	if err != nil {
		return err
	}
	if hasData {
		return errorx.New(errorx.Conflict, "请先删除当前类型下的字典数据")
	}
	if err := l.dictTypeRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 删除字典类型失败", zap.Uint64("dict_type_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) GetDictDataList(ctx context.Context, req *types.GetDictDataListReq) (*types.GetDictDataListResp, error) {
	dictType, err := l.dictTypeRepo.FindByCode(ctx, req.DictCode)
	if err != nil {
		return nil, err
	}
	if dictType == nil {
		return nil, errorx.New(errorx.NotFound, "字典类型不存在")
	}
	data, total, err := l.dictDataRepo.List(ctx, repository.DictDataListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		DictCode: req.DictCode,
		Label:    req.Label,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.DictDataListItem, 0, len(data))
	for _, item := range data {
		list = append(list, types.DictDataListItem{
			ID:        item.ID,
			DictCode:  item.DictCode,
			Label:     item.Label,
			Value:     item.Value,
			Color:     item.Color,
			IsDefault: item.IsDefault,
			Sort:      item.Sort,
			Status:    item.Status,
			Remark:    item.Remark,
			CreatedAt: item.CreatedAt.Format(time.DateTime),
			UpdatedAt: item.UpdatedAt.Format(time.DateTime),
		})
	}
	resp := types.GetDictDataListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func (l *dictLogic) GetDictDataDetail(ctx context.Context, id uint64) (*types.GetDictDataDetailResp, error) {
	data, err := l.dictDataRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errorx.New(errorx.NotFound, "字典数据不存在")
	}
	resp := types.GetDictDataDetailResp{
		ID:        data.ID,
		DictCode:  data.DictCode,
		Label:     data.Label,
		Value:     data.Value,
		Color:     data.Color,
		IsDefault: data.IsDefault,
		Sort:      data.Sort,
		Status:    data.Status,
		Remark:    data.Remark,
	}
	return &resp, nil
}

func (l *dictLogic) CreateDictData(ctx context.Context, req *types.CreateDictDataReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictDataCreateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置新增字典数据防抖失败", zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	if err := l.validateDictData(ctx, 0, req.DictCode, &req.DictDataReq); err != nil {
		return err
	}
	data := &model.FortuneDictData{
		DictCode:  req.DictCode,
		Label:     req.Label,
		Value:     req.Value,
		Color:     req.Color,
		IsDefault: req.IsDefault,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dictDataRepo := repository.NewDictDataRepository(tx)
		id, createErr := dictDataRepo.Create(ctx, data)
		if createErr != nil {
			return createErr
		}
		if req.IsDefault != 1 {
			return nil
		}
		return dictDataRepo.ClearDefaultExcept(ctx, req.DictCode, id, operatorID)
	}); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 新增字典数据失败", zap.String("dict_code", req.DictCode), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) UpdateDictData(ctx context.Context, id uint64, req *types.UpdateDictDataReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictDataUpdateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置修改字典数据防抖失败", zap.Uint64("dict_data_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	data, err := l.dictDataRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return errorx.New(errorx.NotFound, "字典数据不存在")
	}
	if err := l.validateDictData(ctx, id, data.DictCode, &req.DictDataReq); err != nil {
		return err
	}
	updates := map[string]any{
		"label":      req.Label,
		"value":      req.Value,
		"color":      req.Color,
		"is_default": req.IsDefault,
		"sort":       req.Sort,
		"status":     req.Status,
		"remark":     req.Remark,
		"updated_by": operatorID,
	}
	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dictDataRepo := repository.NewDictDataRepository(tx)
		if updateErr := dictDataRepo.Update(ctx, id, updates); updateErr != nil {
			return updateErr
		}
		if req.IsDefault != 1 {
			return nil
		}
		return dictDataRepo.ClearDefaultExcept(ctx, data.DictCode, id, operatorID)
	}); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 修改字典数据失败", zap.Uint64("dict_data_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) DeleteDictData(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.DictDataDeleteDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[DictLogic] Redis设置删除字典数据防抖失败", zap.Uint64("dict_data_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	data, err := l.dictDataRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return errorx.New(errorx.NotFound, "字典数据不存在")
	}
	if err := l.dictDataRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		l.svcCtx.Logger.Error("[DictLogic] 删除字典数据失败", zap.Uint64("dict_data_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *dictLogic) GetDictOptions(ctx context.Context, req *types.GetDictOptionsReq) (*types.GetDictOptionsResp, error) {
	parts := strings.Split(req.Codes, ",")
	codes := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code == "" {
			continue
		}
		if _, exists := seen[code]; exists {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	if len(codes) > maxDictOptionCodes {
		return nil, errorx.New(errorx.InvalidParameter, "单次最多查询50个字典")
	}

	resp := make(types.GetDictOptionsResp, len(codes))
	for _, code := range codes {
		resp[code] = []types.DictOption{}
	}
	var dictTypes []model.FortuneDictType
	var err error
	if len(codes) == 0 {
		dictTypes, err = l.dictTypeRepo.ListEnabled(ctx)
	} else {
		dictTypes, err = l.dictTypeRepo.FindEnabledByCodes(ctx, codes)
	}
	if err != nil {
		return nil, err
	}
	enabledCodes := make([]string, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		enabledCodes = append(enabledCodes, dictType.Code)
	}
	data, err := l.dictDataRepo.ListEnabledByCodes(ctx, enabledCodes)
	if err != nil {
		return nil, err
	}
	for _, item := range data {
		resp[item.DictCode] = append(resp[item.DictCode], types.DictOption{
			Label:     item.Label,
			Value:     item.Value,
			Color:     item.Color,
			IsDefault: item.IsDefault,
			Sort:      item.Sort,
		})
	}
	return &resp, nil
}

func (l *dictLogic) validateDictData(ctx context.Context, id uint64, code string, req *types.DictDataReq) error {
	dictType, err := l.dictTypeRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}
	if dictType == nil {
		return errorx.New(errorx.InvalidParameter, "选择的字典类型不存在")
	}
	existing, err := l.dictDataRepo.FindByCodeAndValue(ctx, code, req.Value)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != id {
		return errorx.New(errorx.Conflict, "当前字典值已存在")
	}
	return nil
}
