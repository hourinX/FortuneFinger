package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type PrizeListFilter struct {
	Page     int
	PageSize int
	Name     string
	Level    string
	Types    uint8
	Status   uint8
}

type PrizeRepository interface {
	List(ctx context.Context, filter PrizeListFilter) ([]model.FortunePrize, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.FortunePrize, error)
	FindByOrgIDAndLevel(ctx context.Context, orgID uint64, level string) (*model.FortunePrize, error)
	Create(ctx context.Context, prize *model.FortunePrize) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type prizeRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewPrizeRepository(db *gorm.DB, generator *snowflake.Generator) PrizeRepository {
	return &prizeRepository{db: db, snowflake: generator}
}

func (r *prizeRepository) List(ctx context.Context, filter PrizeListFilter) ([]model.FortunePrize, int64, error) {
	var prizes []model.FortunePrize
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FortunePrize{}).Where("deleted_at = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Level != "" {
		query = query.Where("level = ?", filter.Level)
	}
	if filter.Types != 0 {
		query = query.Where("types = ?", filter.Types)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询奖品数量失败", err)
	}
	page, pageSize := safePage(filter.Page, filter.PageSize)
	if err := query.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prizes).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询奖品列表失败", err)
	}
	return prizes, total, nil
}

func (r *prizeRepository) FindByID(ctx context.Context, id uint64) (*model.FortunePrize, error) {
	var prize model.FortunePrize
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&prize).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询奖品失败", err)
	}
	return &prize, nil
}

func (r *prizeRepository) FindByOrgIDAndLevel(ctx context.Context, orgID uint64, level string) (*model.FortunePrize, error) {
	var prize model.FortunePrize
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND level = ? AND deleted_at = ?", orgID, level, 0).
		First(&prize).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询企业奖项失败", err)
	}
	return &prize, nil
}

func (r *prizeRepository) Create(ctx context.Context, prize *model.FortunePrize) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成奖品ID失败", err)
	}
	prize.ID = id
	if err := r.db.WithContext(ctx).Create(prize).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增奖品失败", err)
	}
	return id, nil
}

func (r *prizeRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.FortunePrize{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates)
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改奖品失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "奖品不存在")
	}
	return nil
}

func (r *prizeRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortunePrize{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除奖品失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "奖品不存在")
	}
	return nil
}
