package logic

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type MenuLogic interface {
	GetMenuList(ctx context.Context, req *types.GetMenuListReq) (*types.GetMenuListResp, error)
	GetMenuDetail(ctx context.Context, id uint64) (*types.GetMenuDetailResp, error)
	CreateMenu(ctx context.Context, req *types.CreateMenuReq) error
	UpdateMenu(ctx context.Context, id uint64, req *types.UpdateMenuReq) error
	DeleteMenu(ctx context.Context, id uint64) error
	GetRouters(ctx context.Context) (*types.GetRoutersResp, error)
	GetPermission(ctx context.Context) (*types.GetPermissionResp, error)
	GetPermissionTree(ctx context.Context, roleID uint64) (*types.GetPermissionTreeResp, error)
}

type menuLogic struct {
	svcCtx   *svc.ServiceContext
	menuRepo repository.MenuRepository
	roleRepo repository.RoleRepository
	permRepo repository.CasbinRuleRepository
}

func NewMenuLogic(svcCtx *svc.ServiceContext) MenuLogic {
	return &menuLogic{
		svcCtx:   svcCtx,
		menuRepo: repository.NewMenuRepository(svcCtx.GormDB, svcCtx.Snowflake),
		roleRepo: repository.NewRoleRepository(svcCtx.GormDB, svcCtx.Snowflake),
		permRepo: repository.NewCasbinRuleRepository(svcCtx.GormDB),
	}
}

func (l *menuLogic) GetMenuList(ctx context.Context, req *types.GetMenuListReq) (*types.GetMenuListResp, error) {
	menus, err := l.menuRepo.List(ctx, repository.MenuListFilter{
		Title:  req.Title,
		Types:  req.Types,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	menuIDs := make(map[uint64]struct{}, len(menus))
	for _, menu := range menus {
		menuIDs[menu.ID] = struct{}{}
	}

	collectParentIDs := func(items []model.FortuneMenu) []uint64 {
		parentIDs := make([]uint64, 0)
		seen := make(map[uint64]struct{})

		for _, menu := range items {
			if menu.PID == 0 {
				continue
			}
			if _, exists := menuIDs[menu.PID]; exists {
				continue
			}
			if _, exists := seen[menu.PID]; exists {
				continue
			}
			seen[menu.PID] = struct{}{}
			parentIDs = append(parentIDs, menu.PID)
		}
		return parentIDs
	}

	parentIDs := collectParentIDs(menus)
	for len(parentIDs) > 0 {
		parents, findErr := l.menuRepo.FindByIDs(ctx, parentIDs)
		if findErr != nil {
			return nil, findErr
		}
		if len(parents) == 0 {
			break
		}
		menus = append(menus, parents...)
		for _, parent := range parents {
			menuIDs[parent.ID] = struct{}{}
		}
		parentIDs = collectParentIDs(parents)
	}

	sort.Slice(menus, func(i, j int) bool {
		if menus[i].PID != menus[j].PID {
			return menus[i].PID < menus[j].PID
		}
		if menus[i].Sort != menus[j].Sort {
			return menus[i].Sort < menus[j].Sort
		}
		return menus[i].ID < menus[j].ID
	})

	childrenByPID := make(map[uint64][]model.FortuneMenu)
	for _, menu := range menus {
		childrenByPID[menu.PID] = append(childrenByPID[menu.PID], menu)
	}

	var buildItem func(model.FortuneMenu) types.MenuListItem
	buildItem = func(menu model.FortuneMenu) types.MenuListItem {
		item := types.MenuListItem{
			ID:        menu.ID,
			PID:       menu.PID,
			Types:     menu.Types,
			Title:     menu.Title,
			Node:      menu.Node,
			Route:     menu.Route,
			Code:      menu.Code,
			Icon:      menu.Icon,
			Remark:    menu.Remark,
			Sort:      menu.Sort,
			KeepAlive: menu.KeepAlive,
			Status:    menu.Status,
			CreatedAt: menu.CreatedAt.Format(time.DateTime),
			UpdatedAt: menu.UpdatedAt.Format(time.DateTime),
			Children:  make([]types.MenuListItem, 0, len(childrenByPID[menu.ID])),
		}
		for _, child := range childrenByPID[menu.ID] {
			item.Children = append(item.Children, buildItem(child))
		}
		return item
	}

	list := make([]types.MenuListItem, 0)
	for _, menu := range menus {
		if _, parentExists := menuIDs[menu.PID]; parentExists {
			continue
		}
		list = append(list, buildItem(menu))
	}

	resp := types.GetMenuListResp{List: list}
	return &resp, nil
}

func (l *menuLogic) GetMenuDetail(ctx context.Context, id uint64) (*types.GetMenuDetailResp, error) {
	menu, err := l.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errorx.New(errorx.NotFound, "菜单不存在")
	}

	resp := types.GetMenuDetailResp{
		ID:        menu.ID,
		PID:       menu.PID,
		Types:     menu.Types,
		Title:     menu.Title,
		Node:      menu.Node,
		Route:     menu.Route,
		Code:      menu.Code,
		Icon:      menu.Icon,
		Remark:    menu.Remark,
		Sort:      menu.Sort,
		KeepAlive: menu.KeepAlive,
		Status:    menu.Status,
	}
	return &resp, nil
}

func (l *menuLogic) CreateMenu(ctx context.Context, req *types.CreateMenuReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.MenuCreateDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] Redis设置新增菜单防抖失败", zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	if err := l.validateMenu(ctx, 0, &req.MenuReq); err != nil {
		return err
	}

	menu := &model.FortuneMenu{
		PID:       req.PID,
		Types:     req.Types,
		Title:     req.Title,
		Node:      req.Node,
		Icon:      req.Icon,
		Remark:    req.Remark,
		Sort:      req.Sort,
		KeepAlive: req.KeepAlive,
		Status:    req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	switch req.Types {
	case 1:
		menu.KeepAlive = 2
	case 2:
		menu.Route = &req.Route
		menu.Code = &req.Code
	case 3:
		menu.Code = &req.Code
		menu.Icon = ""
		menu.KeepAlive = 2
	}

	if _, err := l.menuRepo.Create(ctx, menu); err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] 新增菜单失败", zap.String("menu_node", req.Node), zap.Error(err))
		return err
	}
	return nil
}

func (l *menuLogic) UpdateMenu(ctx context.Context, id uint64, req *types.UpdateMenuReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.MenuUpdateDebounceKey + ":" + strconv.FormatUint(id, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] Redis设置修改菜单防抖失败", zap.Uint64("menu_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	menu, err := l.menuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errorx.New(errorx.NotFound, "菜单不存在")
	}
	if req.Types == 2 {
		if menu.Types == 2 && menu.Node == req.Node && menu.Code != nil {
			req.Code = *menu.Code
		} else {
			req.Code = strings.ReplaceAll(req.Node, "/", ":") + ":view"
		}
	}

	if err := l.validateMenu(ctx, id, &req.MenuReq); err != nil {
		return err
	}
	hasChildren, err := l.menuRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren && menu.Types != req.Types {
		return errorx.New(errorx.Conflict, "存在下级节点时不能修改节点类型")
	}

	updates := map[string]any{
		"pid":        req.PID,
		"types":      req.Types,
		"title":      req.Title,
		"node":       req.Node,
		"route":      nil,
		"code":       nil,
		"icon":       req.Icon,
		"remark":     req.Remark,
		"sort":       req.Sort,
		"keep_alive": req.KeepAlive,
		"status":     req.Status,
		"updated_by": operatorID,
	}

	switch req.Types {
	case 1:
		updates["keep_alive"] = uint8(2)
	case 2:
		updates["route"] = req.Route
		updates["code"] = req.Code
	case 3:
		updates["code"] = req.Code
		updates["icon"] = ""
		updates["keep_alive"] = uint8(2)
	}
	if err := l.menuRepo.Update(ctx, id, updates); err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] 修改菜单失败", zap.Uint64("menu_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *menuLogic) DeleteMenu(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}

	debounceKey := constant.MenuDeleteDebounceKey + ":" + strconv.FormatUint(operatorID, 10)
	acquired, err := l.svcCtx.Redis.SetNX(ctx, debounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] Redis设置删除菜单防抖失败", zap.Uint64("menu_id", id), zap.Uint64("operator_id", operatorID), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	menu, err := l.menuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errorx.New(errorx.NotFound, "菜单不存在")
	}

	hasChildren, err := l.menuRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errorx.New(errorx.Conflict, "请先删除当前菜单下的子节点")
	}

	if err := l.menuRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		l.svcCtx.Logger.Error("[MenuLogic] 删除菜单失败", zap.Uint64("menu_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *menuLogic) GetRouters(ctx context.Context) (*types.GetRoutersResp, error) {
	userID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}

	roleCode, err := l.permRepo.FindRoleCodeByUserID(ctx, strconv.FormatUint(userID, 10))
	if err != nil {
		return nil, err
	}

	menuIDs := make([]string, 0)
	if roleCode != "" {
		menuIDs, err = l.permRepo.ListMenuIDsByRoleCode(ctx, roleCode)
		if err != nil {
			return nil, err
		}
	}

	menus, err := l.menuRepo.ListEnabledRouterMenus(ctx)
	if err != nil {
		return nil, err
	}

	hasAllMenuAccess := false
	authorizedMenuIDs := make(map[string]struct{}, len(menuIDs))
	for _, menuID := range menuIDs {
		if menuID == "*" {
			hasAllMenuAccess = true
			break
		}
		authorizedMenuIDs[menuID] = struct{}{}
	}

	menuMap := make(map[uint64]model.FortuneMenu, len(menus))
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	visibleMenuIDs := make(map[uint64]struct{})
	for _, menu := range menus {
		if hasAllMenuAccess {
			visibleMenuIDs[menu.ID] = struct{}{}
			continue
		}
		if menu.Types != 2 {
			continue
		}
		if _, exists := authorizedMenuIDs[strconv.FormatUint(menu.ID, 10)]; !exists {
			continue
		}

		visibleMenuIDs[menu.ID] = struct{}{}
		parentID := menu.PID
		for parentID != 0 {
			parent, exists := menuMap[parentID]
			if !exists || parent.Types != 1 {
				break
			}
			visibleMenuIDs[parent.ID] = struct{}{}
			parentID = parent.PID
		}
	}

	menusByPID := make(map[uint64][]model.FortuneMenu)
	for _, menu := range menus {
		if _, exists := visibleMenuIDs[menu.ID]; exists {
			menusByPID[menu.PID] = append(menusByPID[menu.PID], menu)
		}
	}

	var buildRouterTree func(uint64) []types.RouterItem
	buildRouterTree = func(pid uint64) []types.RouterItem {
		childMenus := menusByPID[pid]
		routers := make([]types.RouterItem, 0, len(childMenus))
		for _, menu := range childMenus {
			routers = append(routers, types.RouterItem{
				ID:        menu.ID,
				Types:     menu.Types,
				Title:     menu.Title,
				Node:      menu.Node,
				Route:     menu.Route,
				Icon:      menu.Icon,
				KeepAlive: menu.KeepAlive,
				Children:  buildRouterTree(menu.ID),
			})
		}
		return routers
	}

	resp := types.GetRoutersResp{List: buildRouterTree(0)}
	return &resp, nil
}

func (l *menuLogic) GetPermission(ctx context.Context) (*types.GetPermissionResp, error) {
	userID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}
	roleCode, err := l.permRepo.FindRoleCodeByUserID(ctx, strconv.FormatUint(userID, 10))
	if err != nil {
		return nil, err
	}
	if roleCode == "" {
		return &types.GetPermissionResp{List: []string{}}, nil
	}
	menuIDs, err := l.permRepo.ListMenuIDsByRoleCode(ctx, roleCode)
	if err != nil {
		return nil, err
	}
	hasAllPermission := false
	authorizedMenuIDs := make(map[string]struct{}, len(menuIDs))
	for _, menuID := range menuIDs {
		if menuID == "*" {
			hasAllPermission = true
			break
		}
		authorizedMenuIDs[menuID] = struct{}{}
	}
	menus, err := l.menuRepo.List(ctx, repository.MenuListFilter{Types: 3, Status: 1})
	if err != nil {
		return nil, err
	}
	permissions := make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu.Code == nil || *menu.Code == "" {
			continue
		}
		if !hasAllPermission {
			if _, exists := authorizedMenuIDs[strconv.FormatUint(menu.ID, 10)]; !exists {
				continue
			}
		}
		permissions = append(permissions, *menu.Code)
	}
	return &types.GetPermissionResp{List: permissions}, nil
}

func (l *menuLogic) GetPermissionTree(ctx context.Context, roleID uint64) (*types.GetPermissionTreeResp, error) {
	role, err := l.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errorx.New(errorx.NotFound, "角色不存在")
	}
	menus, err := l.menuRepo.List(ctx, repository.MenuListFilter{Status: 1})
	if err != nil {
		return nil, err
	}

	allPermissionIDs := make(map[string]struct{})
	for _, menu := range menus {
		if menu.Types == 2 || menu.Types == 3 {
			allPermissionIDs[strconv.FormatUint(menu.ID, 10)] = struct{}{}
		}
	}

	allowedIDs := make(map[string]struct{})
	if role.ParentID == 0 {
		allowedIDs = allPermissionIDs
	} else {
		parent, findErr := l.roleRepo.FindByID(ctx, role.ParentID)
		if findErr != nil {
			return nil, findErr
		}
		if parent == nil {
			return nil, errorx.New(errorx.Conflict, "角色的父级不存在")
		}
		if parent.ParentID == 0 {
			allowedIDs = allPermissionIDs
		} else {
			parentMenuIDs, listErr := l.permRepo.ListMenuIDsByRoleCode(ctx, parent.Code)
			if listErr != nil {
				return nil, listErr
			}

			for _, menuID := range parentMenuIDs {
				if menuID == "*" {
					allowedIDs = allPermissionIDs
					break
				}
				if _, exists := allPermissionIDs[menuID]; exists {
					allowedIDs[menuID] = struct{}{}
				}
			}
		}
	}

	selectedIDs := make(map[string]struct{})
	if role.ParentID == 0 {
		selectedIDs = allPermissionIDs
	} else {
		roleMenuIDs, listErr := l.permRepo.ListMenuIDsByRoleCode(ctx, role.Code)
		if listErr != nil {
			return nil, listErr
		}

		for _, menuID := range roleMenuIDs {
			if menuID == "*" {
				selectedIDs = allPermissionIDs
				break
			}

			if _, exists := allowedIDs[menuID]; exists {
				selectedIDs[menuID] = struct{}{}
			}
		}
	}

	childrenByPID := make(map[uint64][]model.FortuneMenu)
	for _, menu := range menus {
		childrenByPID[menu.PID] = append(childrenByPID[menu.PID], menu)
	}

	var buildTree func(uint64) []types.PermissionTreeItem
	buildTree = func(pid uint64) []types.PermissionTreeItem {
		items := make([]types.PermissionTreeItem, 0)
		for _, menu := range childrenByPID[pid] {
			children := buildTree(menu.ID)
			menuID := strconv.FormatUint(menu.ID, 10)
			if menu.Types != 1 {
				if _, exists := allowedIDs[menuID]; !exists {
					continue
				}
			} else if len(children) == 0 && role.ParentID != 0 {
				continue
			}
			code := ""
			if menu.Code != nil {
				code = *menu.Code
			}
			items = append(items, types.PermissionTreeItem{
				ID:       menu.ID,
				Name:     menu.Title,
				Code:     code,
				Children: children,
			})
		}
		return items
	}

	selected := make([]string, 0, len(selectedIDs))
	for _, menu := range menus {
		menuID := strconv.FormatUint(menu.ID, 10)
		if _, exists := selectedIDs[menuID]; exists {
			selected = append(selected, menuID)
		}
	}
	resp := types.GetPermissionTreeResp{
		List:        buildTree(0),
		SelectedIDs: selected,
	}
	return &resp, nil
}

func (l *menuLogic) validateMenu(ctx context.Context, id uint64, req *types.MenuReq) error {
	if req.PID == 0 && req.Types == 3 {
		return errorx.New(errorx.InvalidParameter, "请选择权限所属的菜单")
	}
	if req.PID != 0 {
		parent, err := l.menuRepo.FindByID(ctx, req.PID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errorx.New(errorx.InvalidParameter, "选择的所属父级不存在，请重新选择")
		}

		parentType := uint8(1)
		message := "目录和菜单只能添加在目录下"
		if req.Types == 3 {
			parentType = 2
			message = "权限只能添加在菜单下"
		}
		if parent.Types != parentType {
			return errorx.New(errorx.InvalidParameter, message)
		}
		for id != 0 && parent != nil {
			if parent.ID == id {
				return errorx.New(errorx.InvalidParameter, "不能选择自己或自己的下级作为所属父级")
			}
			if parent.PID == 0 {
				break
			}
			parent, err = l.menuRepo.FindByID(ctx, parent.PID)
			if err != nil {
				return err
			}
		}
	}

	if req.Types == 2 {
		if req.Route == "" {
			return errorx.New(errorx.InvalidParameter, "菜单路由不能为空")
		}
		if id == 0 {
			req.Code = strings.ReplaceAll(req.Node, "/", ":") + ":view"
		}
	}
	if (req.Types == 2 || req.Types == 3) && req.Code == "" {
		return errorx.New(errorx.InvalidParameter, "权限编码不能为空")
	}
	if (req.Types == 2 || req.Types == 3) && utf8.RuneCountInString(req.Code) > 100 {
		return errorx.New(errorx.InvalidParameter, "权限编码长度不能超过100个字符")
	}

	var nodeMenu, routeMenu, codeMenu *model.FortuneMenu
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		var findErr error
		nodeMenu, findErr = l.menuRepo.FindByNode(groupCtx, req.Node)
		return findErr
	})
	if req.Types == 2 {
		group.Go(func() error {
			var findErr error
			routeMenu, findErr = l.menuRepo.FindByRoute(groupCtx, req.Route)
			return findErr
		})
	}
	if req.Types == 2 || req.Types == 3 {
		group.Go(func() error {
			var findErr error
			codeMenu, findErr = l.menuRepo.FindByCode(groupCtx, req.Code)
			return findErr
		})
	}
	if err := group.Wait(); err != nil {
		return err
	}
	if nodeMenu != nil && nodeMenu.ID != id {
		return errorx.New(errorx.Conflict, "菜单节点已存在")
	}
	if routeMenu != nil && routeMenu.ID != id {
		return errorx.New(errorx.Conflict, "菜单路由已存在")
	}
	if codeMenu != nil && codeMenu.ID != id {
		return errorx.New(errorx.Conflict, "权限编码已存在")
	}
	return nil
}
