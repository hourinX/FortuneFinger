package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"

	"gorm.io/gorm"
)

type DictDataListFilter struct {
	Page     int
	PageSize int
	DictCode string
	Label    string
	Status   uint8
}

type DictDataRepository interface {
	List(ctx context.Context, filter DictDataListFilter) ([]model.FortuneDictData, int64, error)
	ListEnabledByCodes(ctx context.Context, codes []string) ([]model.FortuneDictData, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneDictData, error)
	FindByCodeAndValue(ctx context.Context, code, value string) (*model.FortuneDictData, error)
	HasByCode(ctx context.Context, code string) (bool, error)
	Create(ctx context.Context, data *model.FortuneDictData) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	ClearDefaultExcept(ctx context.Context, code string, id, updatedBy uint64) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type dictDataRepository struct {
	db *gorm.DB
}

func NewDictDataRepository(db *gorm.DB) DictDataRepository {
	return &dictDataRepository{db: db}
}

func (r *dictDataRepository) List(ctx context.Context, filter DictDataListFilter) ([]model.FortuneDictData, int64, error) {
	var data []model.FortuneDictData
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FortuneDictData{}).Where("deleted_at = ?", 0)
	if filter.DictCode != "" {
		query = query.Where("dict_code = ?", filter.DictCode)
	}
	if filter.Label != "" {
		query = query.Where("label LIKE ?", "%"+filter.Label+"%")
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询字典数据数量失败", err)
	}
	page, pageSize := safePage(filter.Page, filter.PageSize)
	if err := query.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&data).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询字典数据列表失败", err)
	}
	return data, total, nil
}

func (r *dictDataRepository) ListEnabledByCodes(ctx context.Context, codes []string) ([]model.FortuneDictData, error) {
	if len(codes) == 0 {
		return []model.FortuneDictData{}, nil
	}
	var data []model.FortuneDictData
	if err := r.db.WithContext(ctx).
		Where("dict_code IN ? AND status = ? AND deleted_at = ?", codes, 1, 0).
		Order("dict_code ASC, sort ASC, id ASC").
		Find(&data).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询可用字典数据失败", err)
	}
	return data, nil
}

func (r *dictDataRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneDictData, error) {
	var data model.FortuneDictData
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询字典数据失败", err)
	}
	return &data, nil
}

func (r *dictDataRepository) FindByCodeAndValue(ctx context.Context, code, value string) (*model.FortuneDictData, error) {
	var data model.FortuneDictData
	err := r.db.WithContext(ctx).
		Where("dict_code = ? AND value = ? AND deleted_at = ?", code, value, 0).
		First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询字典数据值失败", err)
	}
	return &data, nil
}

func (r *dictDataRepository) HasByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneDictData{}).
		Where("dict_code = ? AND deleted_at = ?", code, 0).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询字典类型关联数据失败", err)
	}
	return count > 0, nil
}

func (r *dictDataRepository) Create(ctx context.Context, data *model.FortuneDictData) (uint64, error) {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增字典数据失败", err)
	}
	return data.ID, nil
}

func (r *dictDataRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneDictData{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改字典数据失败", err)
	}
	return nil
}

func (r *dictDataRepository) ClearDefaultExcept(ctx context.Context, code string, id, updatedBy uint64) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneDictData{}).
		Where("dict_code = ? AND id <> ? AND is_default = ? AND deleted_at = ?", code, id, 1, 0).
		Updates(map[string]any{
			"is_default": 2,
			"updated_by": updatedBy,
		}).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "更新字典默认选项失败", err)
	}
	return nil
}

func (r *dictDataRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneDictData{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除字典数据失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "字典数据不存在")
	}
	return nil
}
