package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type OrgUserRepository interface {
	HasOrgBinding(ctx context.Context, orgID uint64) (bool, error)
	FindByUserID(ctx context.Context, userID uint64) (*model.FortuneOrgUser, error)
	ListByUserIDs(ctx context.Context, userIDs []uint64) ([]model.FortuneOrgUser, error)
	ListUserIDsByCompanyIDs(ctx context.Context, companyIDs []uint64) ([]uint64, error)
	ListUserIDsByDepartmentID(ctx context.Context, departmentID uint64) ([]uint64, error)
	Create(ctx context.Context, orgUser *model.FortuneOrgUser) (uint64, error)
	UpdateByUserID(ctx context.Context, userID uint64, updates map[string]any) error
	DeleteByUserID(ctx context.Context, userID, deletedAt, updatedBy uint64) error
}

type orgUserRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewOrgUserRepository(db *gorm.DB, generator *snowflake.Generator) OrgUserRepository {
	return &orgUserRepository{db: db, snowflake: generator}
}

func (r *orgUserRepository) FindByUserID(ctx context.Context, userID uint64) (*model.FortuneOrgUser, error) {
	var orgUser model.FortuneOrgUser
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at = ?", userID, 0).First(&orgUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户组织归属失败", err)
	}
	return &orgUser, nil
}

func (r *orgUserRepository) ListByUserIDs(ctx context.Context, userIDs []uint64) ([]model.FortuneOrgUser, error) {
	if len(userIDs) == 0 {
		return []model.FortuneOrgUser{}, nil
	}
	var orgUsers []model.FortuneOrgUser
	if err := r.db.WithContext(ctx).
		Where("user_id IN ? AND deleted_at = ?", userIDs, 0).
		Find(&orgUsers).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户组织归属失败", err)
	}
	return orgUsers, nil
}

func (r *orgUserRepository) ListUserIDsByCompanyIDs(ctx context.Context, companyIDs []uint64) ([]uint64, error) {
	if len(companyIDs) == 0 {
		return []uint64{}, nil
	}
	var userIDs []uint64
	if err := r.db.WithContext(ctx).Model(&model.FortuneOrgUser{}).
		Where("company_id IN ? AND deleted_at = ?", companyIDs, 0).
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询公司用户失败", err)
	}
	return userIDs, nil
}

func (r *orgUserRepository) ListUserIDsByDepartmentID(ctx context.Context, departmentID uint64) ([]uint64, error) {
	var userIDs []uint64
	if err := r.db.WithContext(ctx).Model(&model.FortuneOrgUser{}).
		Where("department_id = ? AND deleted_at = ?", departmentID, 0).
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询部门用户失败", err)
	}
	return userIDs, nil
}

func (r *orgUserRepository) Create(ctx context.Context, orgUser *model.FortuneOrgUser) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成用户组织归属ID失败", err)
	}
	orgUser.ID = id
	if err := r.db.WithContext(ctx).Create(orgUser).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增用户组织归属失败", err)
	}
	return id, nil
}

func (r *orgUserRepository) UpdateByUserID(ctx context.Context, userID uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneOrgUser{}).
		Where("user_id = ? AND deleted_at = ?", userID, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改用户组织归属失败", err)
	}
	return nil
}

func (r *orgUserRepository) DeleteByUserID(ctx context.Context, userID, deletedAt, updatedBy uint64) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneOrgUser{}).
		Where("user_id = ? AND deleted_at = ?", userID, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		}).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除用户组织归属失败", err)
	}
	return nil
}

func (r *orgUserRepository) HasOrgBinding(ctx context.Context, orgID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneOrgUser{}).
		Where("deleted_at = ? AND (company_id = ? OR department_id = ?)", 0, orgID, orgID).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询组织用户归属失败", err)
	}
	return count > 0, nil
}
