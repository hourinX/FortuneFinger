package logic

import (
	"context"
	"strconv"
	"time"

	accesscontrol "fortunefinger/backend/internal/biz/access-control"
	"fortunefinger/backend/internal/biz/operator"
	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/constant"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/mask"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/server/user/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const bcryptMaxPasswordBytes = 72

type UserLogic interface {
	GetUserList(ctx context.Context, req *types.GetUserListReq) (*types.GetUserListResp, error)
	GetUserDetail(ctx context.Context, id uint64) (*types.GetUserDetailResp, error)
	CreateUser(ctx context.Context, req *types.CreateUserReq) error
	UpdateUser(ctx context.Context, id uint64, req *types.UpdateUserReq) error
	DeleteUser(ctx context.Context, id uint64) error

	AssignUserRole(ctx context.Context, id uint64, req *types.AssignUserRoleReq) error
	GetUserInfo(ctx context.Context) (*types.GetUserInfoResp, error)
}

type userLogic struct {
	svcCtx      *svc.ServiceContext
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	orgRepo     repository.OrgRepository
	orgUserRepo repository.OrgUserRepository
	access      *accesscontrol.Service
}

func NewUserLogic(svcCtx *svc.ServiceContext) UserLogic {
	return &userLogic{
		svcCtx:      svcCtx,
		userRepo:    repository.NewUserRepository(svcCtx.GormDB, svcCtx.Snowflake),
		roleRepo:    repository.NewRoleRepository(svcCtx.GormDB, svcCtx.Snowflake),
		orgRepo:     repository.NewOrgRepository(svcCtx.GormDB, svcCtx.Snowflake),
		orgUserRepo: repository.NewOrgUserRepository(svcCtx.GormDB, svcCtx.Snowflake),
		access:      accesscontrol.NewService(svcCtx),
	}
}

func (l *userLogic) GetUserList(ctx context.Context, req *types.GetUserListReq) (*types.GetUserListResp, error) {
	var scopedUserIDs []uint64
	limitByOrg := req.OrgID != 0
	if limitByOrg {
		org, err := l.orgRepo.FindByID(ctx, req.OrgID)
		if err != nil {
			return nil, err
		}
		if org == nil || org.Status != 1 {
			return nil, errorx.New(errorx.InvalidParameter, "选择的组织不存在或已被禁用")
		}
		switch org.Types {
		case 1:
			companies, findErr := l.orgRepo.ListByParentID(ctx, org.ID)
			if findErr != nil {
				return nil, findErr
			}
			companyIDs := make([]uint64, 0, len(companies))
			for _, company := range companies {
				if company.Types == 2 && company.Status == 1 {
					companyIDs = append(companyIDs, company.ID)
				}
			}
			scopedUserIDs, err = l.orgUserRepo.ListUserIDsByCompanyIDs(ctx, companyIDs)
			if err != nil {
				return nil, err
			}
		case 2:
			scopedUserIDs, err = l.orgUserRepo.ListUserIDsByCompanyIDs(ctx, []uint64{org.ID})
			if err != nil {
				return nil, err
			}
		case 3:
			scopedUserIDs, err = l.orgUserRepo.ListUserIDsByDepartmentID(ctx, org.ID)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errorx.New(errorx.InvalidParameter, "组织类型无效")
		}
	}
	users, total, err := l.userRepo.List(ctx, repository.UserListFilter{
		Page:           req.Page,
		PageSize:       req.PageSize,
		Username:       req.Username,
		Nickname:       req.Nickname,
		Status:         req.Status,
		UserIDs:        scopedUserIDs,
		LimitByUserIDs: limitByOrg,
	})
	if err != nil {
		return nil, err
	}
	userIDs := make([]uint64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	roles, err := l.access.ListUserRoles(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	orgUsers, err := l.orgUserRepo.ListByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	orgUserByUserID := make(map[uint64]model.FortuneOrgUser, len(orgUsers))
	departmentIDs := make([]uint64, 0, len(orgUsers))
	seenDepartmentIDs := make(map[uint64]struct{}, len(orgUsers))
	for _, orgUser := range orgUsers {
		orgUserByUserID[orgUser.UserID] = orgUser
		if orgUser.DepartmentID == 0 {
			continue
		}
		if _, exists := seenDepartmentIDs[orgUser.DepartmentID]; exists {
			continue
		}
		seenDepartmentIDs[orgUser.DepartmentID] = struct{}{}
		departmentIDs = append(departmentIDs, orgUser.DepartmentID)
	}
	departments, err := l.orgRepo.FindByIDs(ctx, departmentIDs)
	if err != nil {
		return nil, err
	}
	departmentNameByID := make(map[uint64]string, len(departments))
	for _, department := range departments {
		departmentNameByID[department.ID] = department.Name
	}

	list := make([]types.UserListItem, 0, len(users))
	for _, user := range users {
		role := roles[user.ID]
		orgUser := orgUserByUserID[user.ID]
		list = append(list, types.UserListItem{
			ID:             user.ID,
			Username:       user.Username,
			Nickname:       user.Nickname,
			Avatar:         user.Avatar,
			Email:          user.Email,
			Mobile:         mask.Mobile(user.Mobile),
			Status:         user.Status,
			RoleID:         role.ID,
			RoleName:       role.Name,
			RoleCode:       role.Code,
			DepartmentName: departmentNameByID[orgUser.DepartmentID],
			CreatedAt:      user.CreatedAt.Format(time.DateTime),
			UpdatedAt:      user.UpdatedAt.Format(time.DateTime),
		})
	}

	resp := types.GetUserListResp{
		List:  list,
		Total: total,
	}
	return &resp, nil
}

func (l *userLogic) GetUserDetail(ctx context.Context, id uint64) (*types.GetUserDetailResp, error) {
	user, err := l.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New(errorx.NotFound, "用户不存在")
	}
	roles, err := l.access.ListUserRoles(ctx, []uint64{user.ID})
	if err != nil {
		return nil, err
	}
	role := roles[user.ID]
	orgUser, err := l.orgUserRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	departmentID := uint64(0)
	if orgUser != nil {
		departmentID = orgUser.DepartmentID
	}

	resp := types.GetUserDetailResp{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Email:        user.Email,
		Mobile:       user.Mobile,
		Status:       user.Status,
		RoleID:       role.ID,
		RoleName:     role.Name,
		RoleCode:     role.Code,
		DepartmentID: departmentID,
	}
	return &resp, nil
}

func (l *userLogic) CreateUser(ctx context.Context, req *types.CreateUserReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.UserCreateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] Redis设置新增用户防抖失败", zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}
	if err := l.validateUser(ctx, 0, &req.UserReq); err != nil {
		return err
	}
	companyID, err := l.getDepartmentCompanyID(ctx, req.DepartmentID)
	if err != nil {
		return err
	}
	if len([]byte(req.Password)) > bcryptMaxPasswordBytes {
		return errorx.New(errorx.InvalidParameter, "登录密码不能超过72字节")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] 新增用户密码哈希失败", zap.String("username", req.Username), zap.Error(err))
		return errorx.Wrap(errorx.InternalServer, "用户密码处理失败", err)
	}
	user := &model.FortuneUser{
		Username:  req.Username,
		Nickname:  req.Nickname,
		Password:  string(password),
		Avatar:    req.Avatar,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Status:    req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}
	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx, l.svcCtx.Snowflake)
		userID, createErr := userRepo.Create(ctx, user)
		if createErr != nil {
			return createErr
		}
		orgUserRepo := repository.NewOrgUserRepository(tx, l.svcCtx.Snowflake)
		_, createErr = orgUserRepo.Create(ctx, &model.FortuneOrgUser{
			UserID:       userID,
			CompanyID:    companyID,
			DepartmentID: req.DepartmentID,
			RelationType: 1,
			CreatedBy:    operatorID,
			UpdatedBy:    operatorID,
		})
		return createErr
	}); err != nil {
		l.svcCtx.Logger.Error("[UserLogic] 新增用户失败", zap.String("username", req.Username), zap.Error(err))
		return err
	}
	return nil
}

func (l *userLogic) UpdateUser(ctx context.Context, id uint64, req *types.UpdateUserReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.UserUpdateDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] Redis设置修改用户防抖失败", zap.Uint64("user_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	user, err := l.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errorx.New(errorx.NotFound, "用户不存在")
	}
	if err := l.access.EnsureCanOperateUser(ctx, operatorID, id); err != nil {
		return err
	}
	roleCode, err := l.access.RoleCode(ctx, id)
	if err != nil {
		return err
	}
	orgUser, err := l.orgUserRepo.FindByUserID(ctx, id)
	if err != nil {
		return err
	}
	if roleCode == "super_admin" && req.DepartmentID != 0 {
		return errorx.New(errorx.Conflict, "超级管理员不能分配所属部门")
	}
	if roleCode != "super_admin" && req.DepartmentID == 0 && (orgUser == nil || orgUser.RelationType == 1) {
		return errorx.New(errorx.InvalidParameter, "普通成员必须选择所属部门")
	}
	if operatorID == id && req.Status != 1 {
		l.svcCtx.Logger.Warn("[UserLogic] 禁用当前登录用户被拒绝", zap.Uint64("user_id", id))
		return errorx.New(errorx.Conflict, "不能禁用当前登录用户")
	}
	if err := l.validateUser(ctx, id, &req.UserReq); err != nil {
		return err
	}
	if req.Password != "" && len([]byte(req.Password)) > bcryptMaxPasswordBytes {
		return errorx.New(errorx.InvalidParameter, "登录密码不能超过72字节")
	}

	updates := map[string]any{
		"username":   req.Username,
		"nickname":   req.Nickname,
		"avatar":     req.Avatar,
		"email":      req.Email,
		"mobile":     req.Mobile,
		"status":     req.Status,
		"updated_by": operatorID,
	}
	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			l.svcCtx.Logger.Error("[UserLogic] 修改用户密码哈希失败", zap.Uint64("user_id", id), zap.Error(err))
			return errorx.Wrap(errorx.InternalServer, "用户密码处理失败", err)
		}
		updates["password"] = string(password)
	}
	departmentChanged := orgUser == nil && req.DepartmentID != 0
	if orgUser != nil {
		departmentChanged = orgUser.DepartmentID != req.DepartmentID
	}
	companyID := uint64(0)
	if departmentChanged && req.DepartmentID != 0 {
		companyID, err = l.getDepartmentCompanyID(ctx, req.DepartmentID)
		if err != nil {
			return err
		}
	}
	update := func(db *gorm.DB) error {
		userRepo := repository.NewUserRepository(db, l.svcCtx.Snowflake)
		if updateErr := userRepo.Update(ctx, id, updates); updateErr != nil {
			return updateErr
		}
		if !departmentChanged {
			return nil
		}
		orgUserRepo := repository.NewOrgUserRepository(db, l.svcCtx.Snowflake)
		if orgUser == nil {
			_, updateErr := orgUserRepo.Create(ctx, &model.FortuneOrgUser{
				UserID:       id,
				CompanyID:    companyID,
				DepartmentID: req.DepartmentID,
				RelationType: 1,
				CreatedBy:    operatorID,
				UpdatedBy:    operatorID,
			})
			return updateErr
		}
		orgUpdates := map[string]any{
			"department_id": req.DepartmentID,
			"updated_by":    operatorID,
		}
		if req.DepartmentID != 0 {
			orgUpdates["company_id"] = companyID
		}
		return orgUserRepo.UpdateByUserID(ctx, id, orgUpdates)
	}
	if departmentChanged {
		err = l.svcCtx.GormDB.WithContext(ctx).Transaction(update)
	} else {
		err = update(l.svcCtx.GormDB.WithContext(ctx))
	}
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] 修改用户失败", zap.Uint64("user_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *userLogic) DeleteUser(ctx context.Context, id uint64) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.UserDeleteDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] Redis设置删除用户防抖失败", zap.Uint64("user_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	user, err := l.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errorx.New(errorx.NotFound, "用户不存在")
	}
	if err := l.access.EnsureCanOperateUser(ctx, operatorID, id); err != nil {
		return err
	}
	if operatorID == id {
		l.svcCtx.Logger.Warn("[UserLogic] 删除当前登录用户被拒绝", zap.Uint64("user_id", id))
		return errorx.New(errorx.Conflict, "不能删除当前登录用户")
	}

	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx, l.svcCtx.Snowflake)
		if err := userRepo.Delete(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
			return err
		}
		orgUserRepo := repository.NewOrgUserRepository(tx, l.svcCtx.Snowflake)
		if err := orgUserRepo.DeleteByUserID(ctx, id, uint64(time.Now().Unix()), operatorID); err != nil {
			return err
		}
		permRepo := repository.NewCasbinRuleRepository(tx)
		return permRepo.DeleteUserRoleByUserID(ctx, strconv.FormatUint(id, 10))
	}); err != nil {
		l.svcCtx.Logger.Error("[UserLogic] 删除用户失败", zap.Uint64("user_id", id), zap.Error(err))
		return err
	}
	return nil
}

func (l *userLogic) AssignUserRole(ctx context.Context, id uint64, req *types.AssignUserRoleReq) error {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return err
	}
	acquired, err := l.svcCtx.Redis.SetNX(ctx, constant.UserRoleDebounceKey, "", constant.CommonDebounceTTL)
	if err != nil {
		l.svcCtx.Logger.Error("[UserLogic] Redis设置分配用户角色防抖失败", zap.Uint64("user_id", id), zap.Error(err))
		return errorx.Wrap(errorx.ServiceUnavailable, "服务暂时不可用，请稍后重试", err)
	}
	if !acquired {
		return errorx.New(errorx.TooManyRequests, "操作过于频繁，请勿重复提交")
	}

	user, err := l.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errorx.New(errorx.NotFound, "用户不存在")
	}
	role, err := l.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errorx.New(errorx.InvalidParameter, "选择的角色不存在")
	}
	if role.Status != 1 {
		return errorx.New(errorx.Conflict, "选择的角色已被禁用")
	}
	userID := strconv.FormatUint(id, 10)
	shouldUpdate, err := l.access.ValidateRoleAssignment(ctx, operatorID, id, role.Code)
	if err != nil {
		return err
	}
	if !shouldUpdate {
		return nil
	}

	if err := l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permRepo := repository.NewCasbinRuleRepository(tx)
		if err := permRepo.DeleteUserRoleByUserID(ctx, userID); err != nil {
			return err
		}
		return permRepo.CreateUserRole(ctx, userID, role.Code)
	}); err != nil {
		l.svcCtx.Logger.Error("[UserLogic] 分配用户角色失败", zap.Uint64("user_id", id), zap.Uint64("role_id", req.RoleID), zap.Error(err))
		return err
	}
	return nil
}

func (l *userLogic) GetUserInfo(ctx context.Context) (*types.GetUserInfoResp, error) {
	userID, err := operator.UserID(ctx)
	if err != nil {
		l.svcCtx.Logger.Warn("[UserLogic] 未获取到当前用户信息", zap.Error(err))
		return nil, errorx.New(errorx.Unauthorized, "未登录或登录已过期")
	}

	user, err := l.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		l.svcCtx.Logger.Warn("[UserLogic] 当前用户不存在", zap.Uint64("user_id", userID))
		return nil, errorx.New(errorx.Unauthorized, "登录用户不存在")
	}
	if user.Status != 1 {
		l.svcCtx.Logger.Warn("[UserLogic] 禁用用户获取当前用户信息被拒绝", zap.Uint64("user_id", userID))
		return nil, errorx.New(errorx.Forbidden, "账号已被禁用")
	}
	roleCode, err := l.access.RoleCode(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := types.GetUserInfoResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Mobile:   user.Mobile,
		RoleCode: roleCode,
	}
	return &resp, nil
}

func (l *userLogic) validateUser(ctx context.Context, id uint64, req *types.UserReq) error {
	user, err := l.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if user != nil && user.ID != id {
		return errorx.New(errorx.Conflict, "用户名已存在")
	}
	return nil
}

func (l *userLogic) getDepartmentCompanyID(ctx context.Context, departmentID uint64) (uint64, error) {
	department, err := l.orgRepo.FindByID(ctx, departmentID)
	if err != nil {
		return 0, err
	}
	if department == nil || department.Types != 3 {
		return 0, errorx.New(errorx.InvalidParameter, "选择的所属部门不存在")
	}
	if department.Status != 1 {
		return 0, errorx.New(errorx.Conflict, "选择的所属部门已被禁用")
	}
	company, err := l.orgRepo.FindByID(ctx, department.ParentID)
	if err != nil {
		return 0, err
	}
	if company == nil || company.Types != 2 {
		return 0, errorx.New(errorx.Conflict, "所属部门的公司不存在")
	}
	if company.Status != 1 {
		return 0, errorx.New(errorx.Conflict, "所属部门的公司已被禁用")
	}
	return company.ID, nil
}
