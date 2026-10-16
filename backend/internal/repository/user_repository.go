package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type UserListFilter struct {
	Page           int
	PageSize       int
	Username       string
	Nickname       string
	Status         uint8
	UserIDs        []uint64
	LimitByUserIDs bool
}

type UserRepository interface {
	List(ctx context.Context, filter UserListFilter) ([]model.FortuneUser, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneUser, error)
	FindByUsername(ctx context.Context, username string) (*model.FortuneUser, error)
	FindByEmail(ctx context.Context, email string) (*model.FortuneUser, error)
	FindByMobile(ctx context.Context, mobile string) (*model.FortuneUser, error)
	Create(ctx context.Context, user *model.FortuneUser) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type userRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewUserRepository(db *gorm.DB, generator *snowflake.Generator) UserRepository {
	return &userRepository{db: db, snowflake: generator}
}

func (r *userRepository) List(ctx context.Context, filter UserListFilter) ([]model.FortuneUser, int64, error) {
	var users []model.FortuneUser
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FortuneUser{}).Where("deleted_at = ?", 0)
	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.LimitByUserIDs {
		if len(filter.UserIDs) == 0 {
			return []model.FortuneUser{}, 0, nil
		}
		query = query.Where("id IN ?", filter.UserIDs)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询用户数量失败", err)
	}
	page, pageSize := safePage(filter.Page, filter.PageSize)
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询用户列表失败", err)
	}
	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneUser, error) {
	var user model.FortuneUser
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户失败", err)
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.FortuneUser, error) {
	var user model.FortuneUser
	err := r.db.WithContext(ctx).Where("username = ? AND deleted_at = ?", username, 0).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户失败", err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.FortuneUser, error) {
	var user model.FortuneUser
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at = ?", email, 0).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户邮箱失败", err)
	}
	return &user, nil
}

func (r *userRepository) FindByMobile(ctx context.Context, mobile string) (*model.FortuneUser, error) {
	var user model.FortuneUser
	err := r.db.WithContext(ctx).Where("mobile = ? AND deleted_at = ?", mobile, 0).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询用户手机号失败", err)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.FortuneUser) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成用户ID失败", err)
	}
	user.ID = id
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增用户失败", err)
	}
	return id, nil
}

func (r *userRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	err := r.db.WithContext(ctx).Model(&model.FortuneUser{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates).Error
	if err != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改用户失败", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneUser{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除用户失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "用户不存在")
	}
	return nil
}

func safePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
