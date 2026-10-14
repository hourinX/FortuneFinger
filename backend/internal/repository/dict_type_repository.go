package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type DictTypeListFilter struct {
	Page     int
	PageSize int
	Name     string
	Code     string
	Status   uint8
}

type DictTypeRepository interface {
	List(ctx context.Context, filter DictTypeListFilter) ([]model.FortuneDictType, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneDictType, error)
	FindByCode(ctx context.Context, code string) (*model.FortuneDictType, error)
	ListEnabled(ctx context.Context) ([]model.FortuneDictType, error)
	FindEnabledByCodes(ctx context.Context, codes []string) ([]model.FortuneDictType, error)
	Create(ctx context.Context, dictType *model.FortuneDictType) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type dictTypeRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewDictTypeRepository(db *gorm.DB, generator *snowflake.Generator) DictTypeRepository {
	return &dictTypeRepository{db: db, snowflake: generator}
}

func (r *dictTypeRepository) List(ctx context.Context, filter DictTypeListFilter) ([]model.FortuneDictType, int64, error) {
	var dictTypes []model.FortuneDictType
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FortuneDictType{}).Where("deleted_at = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询字典类型数量失败", err)
	}
	page, pageSize := safePage(filter.Page, filter.PageSize)
	if err := query.Order("is_system ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&dictTypes).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询字典类型列表失败", err)
	}
	return dictTypes, total, nil
}

func (r *dictTypeRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneDictType, error) {
	var dictType model.FortuneDictType
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&dictType).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询字典类型失败", err)
	}
	return &dictType, nil
}

func (r *dictTypeRepository) FindByCode(ctx context.Context, code string) (*model.FortuneDictType, error) {
	var dictType model.FortuneDictType
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at = ?", code, 0).First(&dictType).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询字典编码失败", err)
	}
	return &dictType, nil
}

func (r *dictTypeRepository) ListEnabled(ctx context.Context) ([]model.FortuneDictType, error) {
	var dictTypes []model.FortuneDictType
	if err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at = ?", 1, 0).
		Order("id ASC").
		Find(&dictTypes).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询可用字典类型失败", err)
	}
	return dictTypes, nil
}

func (r *dictTypeRepository) FindEnabledByCodes(ctx context.Context, codes []string) ([]model.FortuneDictType, error) {
	if len(codes) == 0 {
		return []model.FortuneDictType{}, nil
	}
	var dictTypes []model.FortuneDictType
	if err := r.db.WithContext(ctx).
		Where("code IN ? AND status = ? AND deleted_at = ?", codes, 1, 0).
		Find(&dictTypes).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询可用字典类型失败", err)
	}
	return dictTypes, nil
}

func (r *dictTypeRepository) Create(ctx context.Context, dictType *model.FortuneDictType) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成字典类型ID失败", err)
	}
	dictType.ID = id
	if err := r.db.WithContext(ctx).Create(dictType).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增字典类型失败", err)
	}
	return id, nil
}

func (r *dictTypeRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneDictType{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改字典类型失败", err)
	}
	return nil
}

func (r *dictTypeRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneDictType{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除字典类型失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "字典类型不存在")
	}
	return nil
}
