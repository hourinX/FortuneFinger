package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type OrgListFilter struct {
	Name   string
	Code   string
	Types  uint8
	Status uint8
}

type OrgRepository interface {
	List(ctx context.Context, filter OrgListFilter) ([]model.FortuneOrg, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneOrg, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneOrg, error)
	ListByParentID(ctx context.Context, parentID uint64) ([]model.FortuneOrg, error)
	ListCompanies(ctx context.Context) ([]model.FortuneOrg, error)
	FindByParentAndName(ctx context.Context, parentID uint64, name string) (*model.FortuneOrg, error)
	FindByParentAndCode(ctx context.Context, parentID uint64, code string) (*model.FortuneOrg, error)
	HasChildren(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, org *model.FortuneOrg) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type orgRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewOrgRepository(db *gorm.DB, generator *snowflake.Generator) OrgRepository {
	return &orgRepository{db: db, snowflake: generator}
}

func (r *orgRepository) List(ctx context.Context, filter OrgListFilter) ([]model.FortuneOrg, error) {
	var orgs []model.FortuneOrg
	query := r.db.WithContext(ctx).Where("deleted_at = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}
	if filter.Types != 0 {
		query = query.Where("types = ?", filter.Types)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Order("parent_id ASC, sort ASC, id ASC").Find(&orgs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询组织列表失败", err)
	}
	return orgs, nil
}

func (r *orgRepository) FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneOrg, error) {
	var orgs []model.FortuneOrg
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at = ?", ids, 0).Find(&orgs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询组织失败", err)
	}
	return orgs, nil
}

func (r *orgRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneOrg, error) {
	var org model.FortuneOrg
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询组织失败", err)
	}
	return &org, nil
}

func (r *orgRepository) ListByParentID(ctx context.Context, parentID uint64) ([]model.FortuneOrg, error) {
	var orgs []model.FortuneOrg
	if err := r.db.WithContext(ctx).
		Where("parent_id = ? AND deleted_at = ?", parentID, 0).
		Order("sort ASC, id ASC").
		Find(&orgs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询下级组织失败", err)
	}
	return orgs, nil
}

func (r *orgRepository) ListCompanies(ctx context.Context) ([]model.FortuneOrg, error) {
	var orgs []model.FortuneOrg
	if err := r.db.WithContext(ctx).
		Where("parent_id <> ? AND types = ? AND status = ? AND deleted_at = ?", 0, 2, 1, 0).
		Order("sort ASC, id ASC").
		Find(&orgs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询企业选项失败", err)
	}
	return orgs, nil
}

func (r *orgRepository) FindByParentAndName(ctx context.Context, parentID uint64, name string) (*model.FortuneOrg, error) {
	return r.findByParentField(ctx, parentID, "name", name)
}

func (r *orgRepository) FindByParentAndCode(ctx context.Context, parentID uint64, code string) (*model.FortuneOrg, error) {
	return r.findByParentField(ctx, parentID, "code", code)
}

func (r *orgRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneOrg{}).
		Where("parent_id = ? AND deleted_at = ?", id, 0).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询下级组织失败", err)
	}
	return count > 0, nil
}

func (r *orgRepository) Create(ctx context.Context, org *model.FortuneOrg) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成组织ID失败", err)
	}
	org.ID = id
	if err := r.db.WithContext(ctx).Create(org).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增组织失败", err)
	}
	return id, nil
}

func (r *orgRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneOrg{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改组织失败", err)
	}
	return nil
}

func (r *orgRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneOrg{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除组织失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "组织不存在")
	}
	return nil
}

func (r *orgRepository) findByParentField(ctx context.Context, parentID uint64, field, value string) (*model.FortuneOrg, error) {
	var org model.FortuneOrg
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND "+field+" = ? AND deleted_at = ?", parentID, value, 0).
		First(&org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询组织失败", err)
	}
	return &org, nil
}
