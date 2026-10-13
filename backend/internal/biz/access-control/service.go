package accesscontrol

import (
	"context"
	"strconv"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

const (
	superAdminRoleCode = "super_admin"
	maxSuperAdminUsers = 2
)

type Service struct {
	svcCtx   *svc.ServiceContext
	roleRepo repository.RoleRepository
	permRepo repository.CasbinRuleRepository
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		svcCtx:   svcCtx,
		roleRepo: repository.NewRoleRepository(svcCtx.GormDB, svcCtx.Snowflake),
		permRepo: repository.NewCasbinRuleRepository(svcCtx.GormDB),
	}
}

func (s *Service) RoleCode(ctx context.Context, userID uint64) (string, error) {
	return s.permRepo.FindRoleCodeByUserID(ctx, strconv.FormatUint(userID, 10))
}

func (s *Service) ListUserRoles(ctx context.Context, userIDs []uint64) (map[uint64]UserRole, error) {
	values := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		values = append(values, strconv.FormatUint(userID, 10))
	}
	bindings, err := s.permRepo.ListUserRolesByUserIDs(ctx, values)
	if err != nil {
		return nil, err
	}
	roleCodes := make([]string, 0, len(bindings))
	seenCodes := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		if _, exists := seenCodes[binding.RoleCode]; exists {
			continue
		}
		seenCodes[binding.RoleCode] = struct{}{}
		roleCodes = append(roleCodes, binding.RoleCode)
	}
	roles, err := s.roleRepo.FindByCodes(ctx, roleCodes)
	if err != nil {
		return nil, err
	}
	roleByCode := make(map[string]model.FortuneRole, len(roles))
	for _, role := range roles {
		roleByCode[role.Code] = role
	}
	result := make(map[uint64]UserRole, len(bindings))
	for _, binding := range bindings {
		userID, parseErr := strconv.ParseUint(binding.UserID, 10, 64)
		if parseErr != nil {
			continue
		}
		view := UserRole{Code: binding.RoleCode}
		if role, exists := roleByCode[binding.RoleCode]; exists {
			view.ID = role.ID
			view.Name = role.Name
		}
		result[userID] = view
	}
	return result, nil
}

func (s *Service) EnsureCanOperateUser(ctx context.Context, operatorID, targetUserID uint64) error {
	targetRoleCode, err := s.RoleCode(ctx, targetUserID)
	if err != nil {
		return err
	}
	if targetRoleCode != superAdminRoleCode {
		return nil
	}
	operatorRoleCode, err := s.RoleCode(ctx, operatorID)
	if err != nil {
		return err
	}
	if operatorRoleCode == superAdminRoleCode {
		return nil
	}
	s.svcCtx.Logger.Warn("[AccessControl] 普通角色操作超级管理员被拒绝", zap.Uint64("operator_id", operatorID), zap.Uint64("user_id", targetUserID))
	return errorx.New(errorx.Forbidden, "普通角色不能操作超级管理员")
}

func (s *Service) EnsureCanOperateRole(ctx context.Context, operatorID uint64, targetRoleCode string) error {
	if targetRoleCode != superAdminRoleCode {
		return nil
	}
	operatorRoleCode, err := s.RoleCode(ctx, operatorID)
	if err != nil {
		return err
	}
	if operatorRoleCode == superAdminRoleCode {
		return nil
	}
	s.svcCtx.Logger.Warn("[AccessControl] 普通角色操作超级管理员角色被拒绝", zap.Uint64("operator_id", operatorID))
	return errorx.New(errorx.Forbidden, "普通角色不能操作超级管理员角色")
}

func (s *Service) ValidateRoleAssignment(ctx context.Context, operatorID, targetUserID uint64, targetRoleCode string) (bool, error) {
	currentRoleCode, err := s.RoleCode(ctx, targetUserID)
	if err != nil {
		return false, err
	}
	operatorRoleCode, err := s.RoleCode(ctx, operatorID)
	if err != nil {
		return false, err
	}
	if currentRoleCode == superAdminRoleCode && operatorRoleCode != superAdminRoleCode {
		s.svcCtx.Logger.Warn("[AccessControl] 普通角色操作超级管理员被拒绝", zap.Uint64("operator_id", operatorID), zap.Uint64("user_id", targetUserID))
		return false, errorx.New(errorx.Forbidden, "普通角色不能操作超级管理员")
	}
	if currentRoleCode == targetRoleCode {
		return false, nil
	}
	if currentRoleCode == superAdminRoleCode {
		s.svcCtx.Logger.Warn("[AccessControl] 修改超级管理员角色被拒绝", zap.Uint64("operator_id", operatorID), zap.Uint64("user_id", targetUserID))
		return false, errorx.New(errorx.Forbidden, "超级管理员角色不允许变更")
	}
	if targetRoleCode != superAdminRoleCode {
		return true, nil
	}
	if operatorRoleCode != superAdminRoleCode {
		s.svcCtx.Logger.Warn("[AccessControl] 普通角色分配超级管理员被拒绝", zap.Uint64("operator_id", operatorID), zap.Uint64("user_id", targetUserID))
		return false, errorx.New(errorx.Forbidden, "只有超级管理员可以分配超级管理员角色")
	}
	count, err := s.permRepo.CountUserBindingsByRoleCode(ctx, superAdminRoleCode)
	if err != nil {
		return false, err
	}
	if count >= maxSuperAdminUsers {
		s.svcCtx.Logger.Warn("[AccessControl] 超级管理员数量达到上限", zap.Uint64("operator_id", operatorID), zap.Uint64("user_id", targetUserID), zap.Int64("current_count", count))
		return false, errorx.New(errorx.Conflict, "系统最多允许存在两个超级管理员")
	}
	return true, nil
}
