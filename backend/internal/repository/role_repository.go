package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type RoleListFilter struct {
	Name   string
	Code   string
	Status uint8
}

type RoleRepository interface {
	List(ctx context.Context, filter RoleListFilter) ([]model.FortuneRole, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneRole, error)
	FindByCodes(ctx context.Context, codes []string) ([]model.FortuneRole, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneRole, error)
	FindByCode(ctx context.Context, code string) (*model.FortuneRole, error)
	HasTopLevel(ctx context.Context) (bool, error)
	HasChildren(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, role *model.FortuneRole) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

func (r *roleRepository) FindByCodes(ctx context.Context, codes []string) ([]model.FortuneRole, error) {
	if len(codes) == 0 {
		return []model.FortuneRole{}, nil
	}
	var roles []model.FortuneRole
	if err := r.db.WithContext(ctx).
		Where("code IN ? AND deleted_at = ?", codes, 0).
		Find(&roles).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色失败", err)
	}
	return roles, nil
}

type roleRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewRoleRepository(db *gorm.DB, generator *snowflake.Generator) RoleRepository {
	return &roleRepository{db: db, snowflake: generator}
}

func (r *roleRepository) List(ctx context.Context, filter RoleListFilter) ([]model.FortuneRole, error) {
	var roles []model.FortuneRole
	query := r.db.WithContext(ctx).Where("deleted_at = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Order("parent_id ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色列表失败", err)
	}
	return roles, nil
}

func (r *roleRepository) FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneRole, error) {
	var roles []model.FortuneRole
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at = ?", ids, 0).Find(&roles).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色失败", err)
	}
	return roles, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneRole, error) {
	var role model.FortuneRole
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色失败", err)
	}
	return &role, nil
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (*model.FortuneRole, error) {
	var role model.FortuneRole
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at = ?", code, 0).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色编码失败", err)
	}
	return &role, nil
}

func (r *roleRepository) HasTopLevel(ctx context.Context) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneRole{}).
		Where("parent_id = ? AND deleted_at = ?", 0, 0).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询顶级角色失败", err)
	}
	return count > 0, nil
}

func (r *roleRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneRole{}).
		Where("parent_id = ? AND deleted_at = ?", id, 0).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询下级角色失败", err)
	}
	return count > 0, nil
}

func (r *roleRepository) Create(ctx context.Context, role *model.FortuneRole) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成角色ID失败", err)
	}
	role.ID = id
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增角色失败", err)
	}
	return id, nil
}

func (r *roleRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneRole{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改角色失败", err)
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneRole{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除角色失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "角色不存在")
	}
	return nil
}
