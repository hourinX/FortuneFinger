package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"

	"gorm.io/gorm"
)

type CasbinRuleRepository interface {
	ListMenuIDsByRoleCode(ctx context.Context, roleCode string) ([]string, error)
	HasMenuAccess(ctx context.Context, roleCode, menuID string) (bool, error)
	FindRoleCodeByUserID(ctx context.Context, userID string) (string, error)
	ListUserRolesByUserIDs(ctx context.Context, userIDs []string) ([]UserRoleBinding, error)
	HasUserBindingByRoleCode(ctx context.Context, roleCode string) (bool, error)
	HasUserBindingByRoleCodes(ctx context.Context, roleCodes []string) (bool, error)
	CountUserBindingsByRoleCode(ctx context.Context, roleCode string) (int64, error)
	CreateUserRole(ctx context.Context, userID, roleCode string) error
	CreatePermissions(ctx context.Context, roleCode string, menuIDs []string) error
	DeletePermissions(ctx context.Context, roleCode string, menuIDs []string) error
	DeletePermissionsByRoleCode(ctx context.Context, roleCode string) error
	DeletePermissionsByRoleCodes(ctx context.Context, roleCodes []string) error
	DeleteUserRoleByUserID(ctx context.Context, userID string) error
}

type UserRoleBinding struct {
	UserID   string
	RoleCode string
}

type casbinRuleRepository struct {
	db *gorm.DB
}

func NewCasbinRuleRepository(db *gorm.DB) CasbinRuleRepository {
	return &casbinRuleRepository{db: db}
}

func (r *casbinRuleRepository) ListMenuIDsByRoleCode(ctx context.Context, roleCode string) ([]string, error) {
	var menuIDs []string
	err := r.db.WithContext(ctx).Model(&model.CasbinRule{}).
		Where("ptype = ? AND v0 = ?", "p", roleCode).
		Pluck("v1", &menuIDs).Error
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询角色菜单失败", err)
	}
	return menuIDs, nil
}

func (r *casbinRuleRepository) HasMenuAccess(ctx context.Context, roleCode, menuID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.CasbinRule{}).
		Where("ptype = ? AND v0 = ? AND v1 IN ?", "p", roleCode, []string{"*", menuID}).
		Count(&count).Error; err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询角色权限失败", err)
	}
	return count > 0, nil
}

func (r *casbinRuleRepository) FindRoleCodeByUserID(ctx context.Context, userID string) (string, error) {
	var rule model.CasbinRule
	err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ?", "g", userID).
		First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", errorx.Wrap(errorx.DatabaseError, "查询用户角色失败", err)
	}
	return rule.V1, nil
}

func (r *casbinRuleRepository) ListUserRolesByUserIDs(ctx context.Context, userIDs []string) ([]UserRoleBinding, error) {
	if len(userIDs) == 0 {
		return []UserRoleBinding{}, nil
	}
	var rules []model.CasbinRule
	if err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 IN ?", "g", userIDs).
		Find(&rules).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户角色失败", err)
	}
	bindings := make([]UserRoleBinding, 0, len(rules))
	for _, rule := range rules {
		bindings = append(bindings, UserRoleBinding{
			UserID:   rule.V0,
			RoleCode: rule.V1,
		})
	}
	return bindings, nil
}

func (r *casbinRuleRepository) HasUserBindingByRoleCode(ctx context.Context, roleCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CasbinRule{}).
		Where("ptype = ? AND v1 = ?", "g", roleCode).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询角色用户绑定失败", err)
	}
	return count > 0, nil
}

func (r *casbinRuleRepository) HasUserBindingByRoleCodes(ctx context.Context, roleCodes []string) (bool, error) {
	if len(roleCodes) == 0 {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CasbinRule{}).
		Where("ptype = ? AND v1 IN ?", "g", roleCodes).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询角色用户绑定失败", err)
	}
	return count > 0, nil
}

func (r *casbinRuleRepository) CountUserBindingsByRoleCode(ctx context.Context, roleCode string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.CasbinRule{}).
		Where("ptype = ? AND v1 = ?", "g", roleCode).
		Count(&count).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "查询角色用户数量失败", err)
	}
	return count, nil
}

func (r *casbinRuleRepository) CreateUserRole(ctx context.Context, userID, roleCode string) error {
	rule := &model.CasbinRule{
		PType: "g",
		V0:    userID,
		V1:    roleCode,
	}
	if err := r.db.WithContext(ctx).Create(rule).Error; err != nil {
		return errorx.Wrap(errorx.DatabaseError, "新增用户角色绑定失败", err)
	}
	return nil
}

func (r *casbinRuleRepository) CreatePermissions(ctx context.Context, roleCode string, menuIDs []string) error {
	if len(menuIDs) == 0 {
		return nil
	}
	rules := make([]model.CasbinRule, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		rules = append(rules, model.CasbinRule{
			PType: "p",
			V0:    roleCode,
			V1:    menuID,
		})
	}
	if err := r.db.WithContext(ctx).CreateInBatches(rules, 200).Error; err != nil {
		return errorx.Wrap(errorx.DatabaseError, "新增角色权限失败", err)
	}
	return nil
}

func (r *casbinRuleRepository) DeletePermissions(ctx context.Context, roleCode string, menuIDs []string) error {
	if len(menuIDs) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ? AND v1 IN ?", "p", roleCode, menuIDs).
		Delete(&model.CasbinRule{}).Error; err != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除角色权限失败", err)
	}
	return nil
}

func (r *casbinRuleRepository) DeletePermissionsByRoleCode(ctx context.Context, roleCode string) error {
	err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ?", "p", roleCode).
		Delete(&model.CasbinRule{}).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除角色权限失败", err)
	}
	return nil
}

func (r *casbinRuleRepository) DeletePermissionsByRoleCodes(ctx context.Context, roleCodes []string) error {
	if len(roleCodes) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 IN ?", "p", roleCodes).
		Delete(&model.CasbinRule{}).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除角色权限失败", err)
	}
	return nil
}

func (r *casbinRuleRepository) DeleteUserRoleByUserID(ctx context.Context, userID string) error {
	err := r.db.WithContext(ctx).
		Where("ptype = ? AND v0 = ?", "g", userID).
		Delete(&model.CasbinRule{}).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除用户角色绑定失败", err)
	}
	return nil
}
