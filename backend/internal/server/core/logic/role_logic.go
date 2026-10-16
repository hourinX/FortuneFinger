package logic

import (
	"context"
	"sort"
	"strconv"
	"time"

	accesscontrol "fortunefinger/backend/internal/biz/access-control"
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

type RoleLogic interface {
	GetRoleList(ctx context.Context, req *types.GetRoleListReq) (*types.GetRoleListResp, error)
	GetRoleDetail(ctx context.Context, id uint64) (*types.GetRoleDetailResp, error)
	CreateRole(ctx context.Context, req *types.CreateRoleReq) error
	UpdateRole(ctx context.Context, id uint64, req *types.UpdateRoleReq) error
	DeleteRole(ctx context.Context, id uint64) error
	GetRoleOptions(ctx context.Context) (*types.GetRoleOptionsResp, error)
	UpdateRolePermission(ctx context.Context, id uint64, req *types.UpdateRolePermissionReq) error
}

type roleLogic struct {
	svcCtx   *svc.ServiceContext
	roleRepo repository.RoleRepository
	menuRepo repository.MenuRepository
	permRepo repository.CasbinRuleRepository
	access   *accesscontrol.Service
}

func NewRoleLogic(svcCtx *svc.ServiceContext) RoleLogic {
	return &roleLogic{
		svcCtx:   svcCtx,
		roleRepo: repository.NewRoleRepository(svcCtx.GormDB, svcCtx.Snowflake),
		menuRepo: repository.NewMenuRepository(svcCtx.GormDB, svcCtx.Snowflake),
		permRepo: repository.NewCasbinRuleRepository(svcCtx.GormDB),
		access:   accesscontrol.NewService(svcCtx),
	}
}

func (l *roleLogic) GetRoleList(ctx context.Context, req *types.GetRoleListReq) (*types.GetRoleListResp, error) {
	roles, err := l.roleRepo.List(ctx, repository.RoleListFilter{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	roleIDs := make(map[uint64]struct{}, len(roles))
	for _, role := range roles {
		roleIDs[role.ID] = struct{}{}
	}

	collectParentIDs := func(items []model.FortuneRole) []uint64 {
		parentIDs := make([]uint64, 0)
		seen := make(map[uint64]struct{})
		for _, role := range items {
			if role.ParentID == 0 {
				continue
			}
			if _, exists := roleIDs[role.ParentID]; exists {
				continue
			}
			if _, exists := seen[role.ParentID]; exists {
				continue
			}
			seen[role.ParentID] = struct{}{}
			parentIDs = append(parentIDs, role.ParentID)
		}
		return parentIDs
	}

	parentIDs := collectParentIDs(roles)
	for len(parentIDs) > 0 {
		parents, findErr := l.roleRepo.FindByIDs(ctx, parentIDs)
		if findErr != nil {
			return nil, findErr
		}
		if len(parents) == 0 {
			break
		}
		roles = append(roles, parents...)
		for _, parent := range parents {
			roleIDs[parent.ID] = struct{}{}
		}
		parentIDs = collectParentIDs(parents)
	}

	sort.Slice(roles, func(i, j int) bool {
		if roles[i].ParentID != roles[j].ParentID {
			return roles[i].ParentID < roles[j].ParentID
		}
		return roles[i].ID < roles[j].ID
	})

	childrenByParentID := make(map[uint64][]model.FortuneRole)
	for _, role := range roles {
		childrenByParentID[role.ParentID] = append(childrenByParentID[role.ParentID], role)
	}

	var buildItem func(model.FortuneRole) types.RoleListItem
	buildItem = func(role model.FortuneRole) types.RoleListItem {
		item := types.RoleListItem{
			ID:        role.ID,
			ParentID:  role.ParentID,
			Name:      role.Name,
			Code:      role.Code,
			Status:    role.Status,
			CreatedAt: role.CreatedAt.Format(time.DateTime),
			UpdatedAt: role.UpdatedAt.Format(time.DateTime),
			Children:  make([]types.RoleListItem, 0, len(childrenByParentID[role.ID])),
		}
		for _, child := range childrenByParentID[role.ID] {
			item.Children = append(item.Children, buildItem(child))
		}
		return item
	}

	list := make([]types.RoleListItem, 0)
	for _, role := range roles {
		if _, parentExists := roleIDs[role.ParentID]; parentExists {
			continue
		}
		list = append(list, buildItem(role))
	}

	resp := types.GetRoleListResp{List: list}
	return &resp, nil
}

func (l *roleLogic) GetRoleDetail(ctx context.Context, id uint64) (*types.GetRoleDetailResp, error) {
	role, err := l.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errorx.New(errorx.NotFound, "角色不存在")
	}

	resp := types.GetRoleDetailResp{
		ID:       role.ID,
		ParentID: role.ParentID,
		Name:     role.Name,
		Code:     role.Code,
		Status:   role.Status,
	}
	return &resp, nil
}

func (l *roleLogic) CreateRole(ctx context.Context, req *types.CreateRoleReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.RoleCreateDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] Redis设置新增角色防抖失败", zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	if err := l.validateRole(ctx, 0, &req.RoleReq); err != nil {
		return err
	}

	if req.ParentID == 0 {
		if req.Code != "super_admin" {
			return errorx.New(errorx.InvalidParameter, "顶级角色编码必须为super_admin")
		}
		hasTopLevel, err := l.roleRepo.HasTopLevel(ctx)
		if err != nil {
			return err
		}
		if hasTopLevel {
			return errorx.New(errorx.Conflict, "系统只允许存在一个顶级角色")
		}
	}

	roleByCode, err := l.roleRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if roleByCode != nil {
		return errorx.New(errorx.Conflict, "角色编码已存在")
	}

	role := &model.FortuneRole{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if _, err := l.roleRepo.Create(ctx, role); err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] 新增角色失败", zap.String("role_code", req.Code), zap.Error(err))
		return err
	}
	return nil
}

func (l *roleLogic) UpdateRole(ctx context.Context, id uint64, req *types.UpdateRoleReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.RoleUpdateDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] Redis设置修改角色防抖失败", zap.Uint64("role_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	role, err := l.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return errorx.New(errorx.NotFound, "角色不存在")
	}

	if err := l.access.EnsureCanOperateRole(ctx, operatorID, role.Code); err != nil {
		return err
	}

	if role.ParentID == 0 {
		if req.ParentID != 0 {
			return errorx.New(errorx.Forbidden, "顶级角色不允许更改所属父级")
		}
		if req.Status != 1 {
			return errorx.New(errorx.Forbidden, "顶级角色不允许禁用")
		}
	} else if req.ParentID == 0 {
		return errorx.New(errorx.Forbidden, "普通角色不允许设为顶级角色")
	}

	if err := l.validateRole(ctx, id, &req.RoleReq); err != nil {
		return err
	}

	parentChanged := role.ParentID != req.ParentID
	var subtreeRoleCodes []string
	if parentChanged {
		subtreeRoleCodes, err = l.listSubtreeRoleCodes(ctx, id)
		if err != nil {
			return err
		}
		hasUsers, err := l.permRepo.HasUserBindingByRoleCodes(ctx, subtreeRoleCodes)
		if err != nil {
			return err
		}
		if hasUsers {
			return errorx.New(errorx.Conflict, "当前角色或其下级角色已绑定用户，请先更换用户角色")
		}
	}

	updates := map[string]any{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"status":     req.Status,
		"updated_by": operatorID,
	}

	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx, l.svcCtx.Snowflake)
		if err := roleRepo.Update(ctx, id, updates); err != nil {
			return err
		}
		if !parentChanged {
			return nil
		}
		permRepo := repository.NewCasbinRuleRepository(tx)
		return permRepo.DeletePermissionsByRoleCodes(ctx, subtreeRoleCodes)
	}); err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] 修改角色失败", zap.Uint64("role_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *roleLogic) DeleteRole(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.RoleDeleteDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] Redis设置删除角色防抖失败", zap.Uint64("role_id", id), zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	role, err := l.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return errorx.New(errorx.NotFound, "角色不存在")
	}

	if err := l.access.EnsureCanOperateRole(ctx, operatorID, role.Code); err != nil {
		return err
	}
	if role.ParentID == 0 {
		return errorx.New(errorx.Forbidden, "顶级角色不允许删除")
	}

	hasChildren, err := l.roleRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errorx.New(errorx.Conflict, "请先删除当前角色下的子角色")
	}

	hasUsers, err := l.permRepo.HasUserBindingByRoleCode(ctx, role.Code)
	if err != nil {
		return err
	}
	if hasUsers {
		return errorx.New(errorx.Conflict, "当前角色已绑定用户，请先更换用户角色")
	}

	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx, l.svcCtx.Snowflake)
		if err := roleRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
			return err
		}
		permRepo := repository.NewCasbinRuleRepository(tx)
		return permRepo.DeletePermissionsByRoleCode(ctx, role.Code)
	}); err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] 删除角色失败", zap.Uint64("role_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *roleLogic) GetRoleOptions(ctx context.Context) (*types.GetRoleOptionsResp, error) {
	roles, err := l.roleRepo.List(ctx, repository.RoleListFilter{})
	if err != nil {
		return nil, err
	}

	childrenByParentID := make(map[uint64][]model.FortuneRole)
	for _, role := range roles {
		childrenByParentID[role.ParentID] = append(childrenByParentID[role.ParentID], role)
	}

	var buildOptions func(uint64) []types.RoleOption
	buildOptions = func(parentID uint64) []types.RoleOption {
		items := childrenByParentID[parentID]
		options := make([]types.RoleOption, 0, len(items))
		for _, role := range items {
			options = append(options, types.RoleOption{
				Key:      role.ID,
				Label:    role.Name,
				Code:     role.Code,
				Status:   role.Status,
				Children: buildOptions(role.ID),
			})
		}
		return options
	}

	resp := types.GetRoleOptionsResp(buildOptions(0))
	return &resp, nil
}

func (l *roleLogic) UpdateRolePermission(ctx context.Context, id uint64, req *types.UpdateRolePermissionReq) error {
	debounceKey := constant.RolePermissionDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] Redis设置角色权限防抖失败", zap.Uint64("role_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	role, err := l.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return errorx.New(errorx.NotFound, "角色不存在")
	}
	if role.ParentID == 0 {
		return errorx.New(errorx.Forbidden, "顶级角色的权限不允许修改")
	}

	parent, err := l.roleRepo.FindByID(ctx, role.ParentID)
	if err != nil {
		return err
	}
	if parent == nil {
		return errorx.New(errorx.Conflict, "角色的父级不存在")
	}

	menus, err := l.menuRepo.List(ctx, repository.MenuListFilter{Status: 1})
	if err != nil {
		return err
	}

	validIDs := make(map[string]struct{}, len(menus))
	for _, menu := range menus {
		if menu.Types == 2 || menu.Types == 3 {
			validIDs[strconv.FormatUint(menu.ID, 10)] = struct{}{}
		}
	}

	allowedIDs := validIDs
	if parent.ParentID != 0 {
		allowedIDs = make(map[string]struct{})
		parentMenuIDs, listErr := l.permRepo.ListMenuIDsByRoleCode(ctx, parent.Code)
		if listErr != nil {
			return listErr
		}

		for _, menuID := range parentMenuIDs {
			if menuID == "*" {
				allowedIDs = validIDs
				break
			}
			if _, exists := validIDs[menuID]; exists {
				allowedIDs[menuID] = struct{}{}
			}
		}
	}

	wantedIDs := make(map[string]struct{}, len(req.MenuIDs))
	for _, menuID := range req.MenuIDs {
		if _, exists := allowedIDs[menuID]; !exists {
			return errorx.New(errorx.Forbidden, "选择的权限超出当前角色可分配范围")
		}
		wantedIDs[menuID] = struct{}{}
	}

	currentMenuIDs, err := l.permRepo.ListMenuIDsByRoleCode(ctx, role.Code)
	if err != nil {
		return err
	}
	currentIDs := make(map[string]struct{}, len(currentMenuIDs))
	for _, menuID := range currentMenuIDs {
		currentIDs[menuID] = struct{}{}
	}

	createIDs := make([]string, 0)
	for menuID := range wantedIDs {
		if _, exists := currentIDs[menuID]; !exists {
			createIDs = append(createIDs, menuID)
		}
	}

	deleteIDs := make([]string, 0)
	for menuID := range currentIDs {
		if _, exists := wantedIDs[menuID]; !exists {
			deleteIDs = append(deleteIDs, menuID)
		}
	}

	if len(createIDs) == 0 && len(deleteIDs) == 0 {
		return nil
	}
	sort.Strings(createIDs)
	sort.Strings(deleteIDs)

	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permRepo := repository.NewCasbinRuleRepository(tx)
		if err := permRepo.DeletePermissions(ctx, role.Code, deleteIDs); err != nil {
			return err
		}
		return permRepo.CreatePermissions(ctx, role.Code, createIDs)
	}); err != nil {
		l.svcCtx.Logger.Error("[RoleLogic] 修改角色权限失败", zap.Uint64("role_id", id), zap.String("role_code", role.Code), zap.Error(err))
		return err
	}
	return nil
}

func (l *roleLogic) validateRole(ctx context.Context, id uint64, req *types.RoleReq) error {
	if req.ParentID != 0 {
		parent, err := l.roleRepo.FindByID(ctx, req.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errorx.New(errorx.InvalidParameter, "选择的父角色不存在，请重新选择")
		}
		for id != 0 && parent != nil {
			if parent.ID == id {
				return errorx.New(errorx.InvalidParameter, "不能选择自己或自己的下级作为父角色")
			}
			if parent.ParentID == 0 {
				break
			}
			parent, err = l.roleRepo.FindByID(ctx, parent.ParentID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *roleLogic) listSubtreeRoleCodes(ctx context.Context, rootID uint64) ([]string, error) {
	roles, err := l.roleRepo.List(ctx, repository.RoleListFilter{})
	if err != nil {
		return nil, err
	}

	childrenByParentID := make(map[uint64][]model.FortuneRole)
	roleByID := make(map[uint64]model.FortuneRole, len(roles))
	for _, role := range roles {
		roleByID[role.ID] = role
		childrenByParentID[role.ParentID] = append(childrenByParentID[role.ParentID], role)
	}
	root, exists := roleByID[rootID]
	if !exists {
		return nil, errorx.New(errorx.NotFound, "角色不存在")
	}

	codes := make([]string, 0)
	queue := []model.FortuneRole{root}
	for len(queue) > 0 {
		role := queue[0]
		queue = queue[1:]
		codes = append(codes, role.Code)
		queue = append(queue, childrenByParentID[role.ID]...)
	}
	return codes, nil
}
