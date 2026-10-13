package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type MenuListFilter struct {
	Title  string
	Types  uint8
	Status uint8
}

type MenuRepository interface {
	List(ctx context.Context, filter MenuListFilter) ([]model.FortuneMenu, error)
	ListEnabledRouterMenus(ctx context.Context) ([]model.FortuneMenu, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneMenu, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneMenu, error)
	FindByNode(ctx context.Context, node string) (*model.FortuneMenu, error)
	FindByRoute(ctx context.Context, route string) (*model.FortuneMenu, error)
	FindByCode(ctx context.Context, code string) (*model.FortuneMenu, error)
	HasChildren(ctx context.Context, id uint64) (bool, error)
	Create(ctx context.Context, menu *model.FortuneMenu) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type menuRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewMenuRepository(db *gorm.DB, generator *snowflake.Generator) MenuRepository {
	return &menuRepository{db: db, snowflake: generator}
}

func (r *menuRepository) List(ctx context.Context, filter MenuListFilter) ([]model.FortuneMenu, error) {
	var menus []model.FortuneMenu
	query := r.db.WithContext(ctx).Where("deleted_at = ?", 0)
	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Types != 0 {
		query = query.Where("types = ?", filter.Types)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Order("pid ASC, sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询菜单列表失败", err)
	}
	return menus, nil
}

func (r *menuRepository) ListEnabledRouterMenus(ctx context.Context) ([]model.FortuneMenu, error) {
	var menus []model.FortuneMenu
	err := r.db.WithContext(ctx).
		Where("deleted_at = ? AND status = ? AND types IN ?", 0, 1, []int{1, 2}).
		Order("pid ASC, sort ASC, id ASC").
		Find(&menus).Error
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户路由失败", err)
	}
	return menus, nil
}

func (r *menuRepository) FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneMenu, error) {
	var menus []model.FortuneMenu
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at = ?", ids, 0).Find(&menus).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询菜单失败", err)
	}
	return menus, nil
}

func (r *menuRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneMenu, error) {
	var menu model.FortuneMenu
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询菜单失败", err)
	}
	return &menu, nil
}

func (r *menuRepository) FindByNode(ctx context.Context, node string) (*model.FortuneMenu, error) {
	return r.findByField(ctx, "node", node)
}

func (r *menuRepository) FindByRoute(ctx context.Context, route string) (*model.FortuneMenu, error) {
	return r.findByField(ctx, "route", route)
}

func (r *menuRepository) FindByCode(ctx context.Context, code string) (*model.FortuneMenu, error) {
	return r.findByField(ctx, "code", code)
}

func (r *menuRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FortuneMenu{}).
		Where("pid = ? AND deleted_at = ?", id, 0).
		Count(&count).Error
	if err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "查询下级菜单失败", err)
	}
	return count > 0, nil
}

func (r *menuRepository) Create(ctx context.Context, menu *model.FortuneMenu) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成菜单ID失败", err)
	}
	menu.ID = id
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增菜单失败", err)
	}
	return id, nil
}

func (r *menuRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneMenu{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改菜单失败", err)
	}
	return nil
}

func (r *menuRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneMenu{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除菜单失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "菜单不存在")
	}
	return nil
}

func (r *menuRepository) findByField(ctx context.Context, field, value string) (*model.FortuneMenu, error) {
	var menu model.FortuneMenu
	err := r.db.WithContext(ctx).Where(field+" = ? AND deleted_at = ?", value, 0).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询菜单失败", err)
	}
	return &menu, nil
}
