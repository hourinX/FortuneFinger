package permissioncheck

import (
	"context"
	"strconv"

	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type Service struct {
	svcCtx   *svc.ServiceContext
	menuRepo repository.MenuRepository
	permRepo repository.CasbinRuleRepository
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		svcCtx:   svcCtx,
		menuRepo: repository.NewMenuRepository(svcCtx.GormDB, svcCtx.Snowflake),
		permRepo: repository.NewCasbinRuleRepository(svcCtx.GormDB),
	}
}

func (s *Service) CheckPermission(ctx context.Context, userID uint64, code string) (bool, error) {
	roleCode, err := s.permRepo.FindRoleCodeByUserID(ctx, strconv.FormatUint(userID, 10))
	if err != nil {
		s.svcCtx.Logger.Error("[PermissionCheck] 查询当前用户角色失败", zap.Uint64("user_id", userID), zap.String("permission_code", code), zap.Error(err))
		return false, err
	}
	if roleCode == "" {
		s.svcCtx.Logger.Warn("[PermissionCheck] 未分配角色的用户访问受限接口", zap.Uint64("user_id", userID), zap.String("permission_code", code))
		return false, nil
	}

	menu, err := s.menuRepo.FindByCode(ctx, code)
	if err != nil {
		s.svcCtx.Logger.Error("[PermissionCheck] 查询权限节点失败", zap.Uint64("user_id", userID), zap.String("role_code", roleCode), zap.String("permission_code", code), zap.Error(err))
		return false, err
	}
	if menu == nil || menu.Status != 1 || (menu.Types != 2 && menu.Types != 3) {
		s.svcCtx.Logger.Error("[PermissionCheck] 权限码未绑定有效节点", zap.String("permission_code", code))
		return false, nil
	}

	allowed, err := s.permRepo.HasMenuAccess(ctx, roleCode, strconv.FormatUint(menu.ID, 10))
	if err != nil {
		s.svcCtx.Logger.Error("[PermissionCheck] 查询角色权限失败", zap.Uint64("user_id", userID), zap.String("role_code", roleCode), zap.String("permission_code", code), zap.Error(err))
		return false, err
	}
	if !allowed {
		s.svcCtx.Logger.Warn("[PermissionCheck] 用户访问权限不足", zap.Uint64("user_id", userID), zap.String("role_code", roleCode), zap.String("permission_code", code))
	}
	return allowed, nil
}
