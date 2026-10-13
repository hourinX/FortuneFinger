package logic

import (
	"context"
	"sort"
	"time"

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

type OrgLogic interface {
	GetOrgList(ctx context.Context, req *types.GetOrgListReq) (*types.GetOrgListResp, error)
	GetOrgDetail(ctx context.Context, id uint64) (*types.GetOrgDetailResp, error)
	CreateOrg(ctx context.Context, req *types.CreateOrgReq) error
	UpdateOrg(ctx context.Context, id uint64, req *types.UpdateOrgReq) error
	DeleteOrg(ctx context.Context, id uint64) error

	GetOrgOptions(ctx context.Context) (*types.GetOrgOptionsResp, error)
}

type orgLogic struct {
	svcCtx      *svc.ServiceContext
	orgRepo     repository.OrgRepository
	orgUserRepo repository.OrgUserRepository
}

func NewOrgLogic(svcCtx *svc.ServiceContext) OrgLogic {
	return &orgLogic{
		svcCtx:      svcCtx,
		orgRepo:     repository.NewOrgRepository(svcCtx.GormDB, svcCtx.Snowflake),
		orgUserRepo: repository.NewOrgUserRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (l *orgLogic) GetOrgList(ctx context.Context, req *types.GetOrgListReq) (*types.GetOrgListResp, error) {
	orgs, err := l.orgRepo.List(ctx, repository.OrgListFilter{
		Name:   req.Name,
		Code:   req.Code,
		Types:  req.Types,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	orgIDs := make(map[uint64]struct{}, len(orgs))
	for _, org := range orgs {
		orgIDs[org.ID] = struct{}{}
	}
	collectParentIDs := func(items []model.FortuneOrg) []uint64 {
		parentIDs := make([]uint64, 0)
		seen := make(map[uint64]struct{})
		for _, org := range items {
			if org.ParentID == 0 {
				continue
			}
			if _, exists := orgIDs[org.ParentID]; exists {
				continue
			}
			if _, exists := seen[org.ParentID]; exists {
				continue
			}
			seen[org.ParentID] = struct{}{}
			parentIDs = append(parentIDs, org.ParentID)
		}
		return parentIDs
	}

	parentIDs := collectParentIDs(orgs)
	for len(parentIDs) > 0 {
		parents, findErr := l.orgRepo.FindByIDs(ctx, parentIDs)
		if findErr != nil {
			return nil, findErr
		}
		if len(parents) == 0 {
			break
		}
		orgs = append(orgs, parents...)
		for _, parent := range parents {
			orgIDs[parent.ID] = struct{}{}
		}
		parentIDs = collectParentIDs(parents)
	}

	sort.Slice(orgs, func(i, j int) bool {
		if orgs[i].ParentID != orgs[j].ParentID {
			return orgs[i].ParentID < orgs[j].ParentID
		}
		if orgs[i].Sort != orgs[j].Sort {
			return orgs[i].Sort < orgs[j].Sort
		}
		return orgs[i].ID < orgs[j].ID
	})
	childrenByParentID := make(map[uint64][]model.FortuneOrg)
	for _, org := range orgs {
		childrenByParentID[org.ParentID] = append(childrenByParentID[org.ParentID], org)
	}

	var buildItem func(model.FortuneOrg) types.OrgListItem
	buildItem = func(org model.FortuneOrg) types.OrgListItem {
		item := types.OrgListItem{
			ID:        org.ID,
			ParentID:  org.ParentID,
			Name:      org.Name,
			Code:      org.Code,
			Types:     org.Types,
			Sort:      org.Sort,
			Status:    org.Status,
			Remark:    org.Remark,
			CreatedAt: org.CreatedAt.Format(time.DateTime),
			UpdatedAt: org.UpdatedAt.Format(time.DateTime),
			Children:  make([]types.OrgListItem, 0, len(childrenByParentID[org.ID])),
		}
		for _, child := range childrenByParentID[org.ID] {
			item.Children = append(item.Children, buildItem(child))
		}
		return item
	}

	list := make([]types.OrgListItem, 0)
	for _, org := range orgs {
		if _, parentExists := orgIDs[org.ParentID]; parentExists {
			continue
		}
		list = append(list, buildItem(org))
	}
	resp := types.GetOrgListResp{List: list}
	return &resp, nil
}

func (l *orgLogic) GetOrgDetail(ctx context.Context, id uint64) (*types.GetOrgDetailResp, error) {
	org, err := l.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errorx.New(errorx.NotFound, "组织不存在")
	}
	resp := types.GetOrgDetailResp{
		ID:       org.ID,
		ParentID: org.ParentID,
		Name:     org.Name,
		Code:     org.Code,
		Types:    org.Types,
		Sort:     org.Sort,
		Status:   org.Status,
		Remark:   org.Remark,
	}
	return &resp, nil
}

func (l *orgLogic) CreateOrg(ctx context.Context, req *types.CreateOrgReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.OrgCreateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] Redis设置新增组织防抖失败", zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	if err := l.validateOrg(ctx, 0, &req.OrgReq); err != nil {
		return err
	}
	org := &model.FortuneOrg{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Code:      req.Code,
		Types:     req.Types,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if _, err := l.orgRepo.Create(ctx, org); err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] 新增组织失败", zap.String("org_code", req.Code), zap.Error(err))
		return err
	}
	return nil
}

func (l *orgLogic) UpdateOrg(ctx context.Context, id uint64, req *types.UpdateOrgReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.OrgUpdateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] Redis设置修改组织防抖失败", zap.Uint64("org_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	org, err := l.orgRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if org == nil {
		return errorx.New(errorx.NotFound, "组织不存在")
	}
	if err := l.validateOrg(ctx, id, &req.OrgReq); err != nil {
		return err
	}
	hasChildren, err := l.orgRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren && org.Types != req.Types {
		return errorx.New(errorx.Conflict, "存在下级组织时不能修改组织类型")
	}
	updates := map[string]any{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"code":       req.Code,
		"types":      req.Types,
		"sort":       req.Sort,
		"status":     req.Status,
		"remark":     req.Remark,
		"updated_by": operatorID,
	}
	if err := l.orgRepo.Update(ctx, id, updates); err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] 修改组织失败", zap.Uint64("org_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *orgLogic) DeleteOrg(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.OrgDeleteDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] Redis设置删除组织防抖失败", zap.Uint64("org_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	org, err := l.orgRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if org == nil {
		return errorx.New(errorx.NotFound, "组织不存在")
	}
	hasChildren, err := l.orgRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errorx.New(errorx.Conflict, "请先删除当前组织下的子组织")
	}
	hasBinding, err := l.orgUserRepo.HasOrgBinding(ctx, id)
	if err != nil {
		return err
	}
	if hasBinding {
		return errorx.New(errorx.Conflict, "当前组织仍有关联用户，不能删除")
	}
	if err := l.orgRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
		l.svcCtx.Logger.Error("[OrgLogic] 删除组织失败", zap.Uint64("org_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *orgLogic) GetOrgOptions(ctx context.Context) (*types.GetOrgOptionsResp, error) {
	orgs, err := l.orgRepo.List(ctx, repository.OrgListFilter{Status: 1})
	if err != nil {
		return nil, err
	}
	childrenByParentID := make(map[uint64][]model.FortuneOrg)
	for _, org := range orgs {
		childrenByParentID[org.ParentID] = append(childrenByParentID[org.ParentID], org)
	}
	var buildTree func(uint64) []types.OrgOption
	buildTree = func(parentID uint64) []types.OrgOption {
		children := childrenByParentID[parentID]
		options := make([]types.OrgOption, 0, len(children))
		for _, org := range children {
			options = append(options, types.OrgOption{
				ID:       org.ID,
				Name:     org.Name,
				Children: buildTree(org.ID),
			})
		}
		return options
	}
	resp := types.GetOrgOptionsResp{List: buildTree(0)}
	return &resp, nil
}

func (l *orgLogic) validateOrg(ctx context.Context, id uint64, req *types.OrgReq) error {
	if req.ParentID == 0 {
		if req.Types != 1 {
			return errorx.New(errorx.InvalidParameter, "顶级组织只能是集团")
		}
	} else {
		parent, err := l.orgRepo.FindByID(ctx, req.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errorx.New(errorx.InvalidParameter, "选择的所属组织不存在，请重新选择")
		}
		expectedParentType := req.Types - 1
		if req.Types == 1 || parent.Types != expectedParentType {
			return errorx.New(errorx.InvalidParameter, "组织层级必须按照集团、公司、部门依次建立")
		}
		for id != 0 && parent != nil {
			if parent.ID == id {
				return errorx.New(errorx.InvalidParameter, "不能选择自己或自己的下级作为所属组织")
			}
			if parent.ParentID == 0 {
				break
			}
			parent, err = l.orgRepo.FindByID(ctx, parent.ParentID)
			if err != nil {
				return err
			}
		}
	}

	var nameOrg, codeOrg *model.FortuneOrg
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		var findErr error
		nameOrg, findErr = l.orgRepo.FindByParentAndName(groupCtx, req.ParentID, req.Name)
		return findErr
	})
	group.Go(func() error {
		var findErr error
		codeOrg, findErr = l.orgRepo.FindByParentAndCode(groupCtx, req.ParentID, req.Code)
		return findErr
	})
	if err := group.Wait(); err != nil {
		return err
	}
	if nameOrg != nil && nameOrg.ID != id {
		return errorx.New(errorx.Conflict, "同级组织名称已存在")
	}
	if codeOrg != nil && codeOrg.ID != id {
		return errorx.New(errorx.Conflict, "同级组织编码已存在")
	}
	return nil
}
