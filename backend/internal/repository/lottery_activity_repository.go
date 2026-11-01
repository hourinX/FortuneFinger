package repository

import (
	"context"
	"errors"
	"time"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type LotteryActivityListFilter struct {
	Page     int
	PageSize int
	Name     string
	OrgID    uint64
	Date     string
	Phase    string
	Status   uint8
}

type LotteryActivityRepository interface {
	List(ctx context.Context, filter LotteryActivityListFilter) ([]model.FortuneActivity, int64, error)
	ListAll(ctx context.Context) ([]model.FortuneActivity, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneActivity, error)
	FindByCode(ctx context.Context, code string) (*model.FortuneActivity, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Create(ctx context.Context, activity *model.FortuneActivity) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type lotteryActivityRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewLotteryActivityRepository(db *gorm.DB, generator *snowflake.Generator) LotteryActivityRepository {
	return &lotteryActivityRepository{db: db, snowflake: generator}
}

func (r *lotteryActivityRepository) List(ctx context.Context, filter LotteryActivityListFilter) ([]model.FortuneActivity, int64, error) {
	var activities []model.FortuneActivity
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FortuneActivity{}).Where("deleted_at = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.OrgID != 0 {
		query = query.Where("org_id = ?", filter.OrgID)
	}
	if filter.Date != "" {
		date, err := time.Parse(time.DateOnly, filter.Date)
		if err != nil {
			return nil, 0, errorx.Wrap(errorx.InvalidParameter, "活动日期格式错误", err)
		}
		query = query.Where("activity_start_time >= ? AND activity_start_time < ?", date, date.AddDate(0, 0, 1))
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	now, err := time.Parse(time.DateTime, time.Now().Format(time.DateTime))
	if err != nil {
		return nil, 0, errorx.Wrap(errorx.InternalServer, "生成当前活动时间失败", err)
	}
	switch filter.Phase {
	case "pending":
		query = query.Where("activity_start_time > ?", now)
	case "active":
		query = query.Where("activity_start_time <= ? AND activity_end_time >= ?", now, now)
	case "finished":
		query = query.Where("activity_end_time < ?", now)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询抽奖活动数量失败", err)
	}
	page, pageSize := safePage(filter.Page, filter.PageSize)
	if err := query.Order("activity_start_time DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&activities).Error; err != nil {
		return nil, 0, errorx.Wrap(errorx.DatabaseError, "查询抽奖活动列表失败", err)
	}
	return activities, total, nil
}

func (r *lotteryActivityRepository) ListAll(ctx context.Context) ([]model.FortuneActivity, error) {
	var activities []model.FortuneActivity
	if err := r.db.WithContext(ctx).
		Where("deleted_at = ?", 0).
		Order("activity_start_time DESC, id DESC").
		Find(&activities).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询抽奖活动失败", err)
	}
	return activities, nil
}

func (r *lotteryActivityRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneActivity, error) {
	var activity model.FortuneActivity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询抽奖活动失败", err)
	}
	return &activity, nil
}

func (r *lotteryActivityRepository) FindByCode(ctx context.Context, code string) (*model.FortuneActivity, error) {
	var activity model.FortuneActivity
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at = ?", code, 0).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询抽奖活动失败", err)
	}
	return &activity, nil
}

func (r *lotteryActivityRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.FortuneActivity{}).Where("code = ?", code).Count(&total).Error; err != nil {
		return false, errorx.Wrap(errorx.DatabaseError, "检查抽奖活动编码失败", err)
	}
	return total > 0, nil
}

func (r *lotteryActivityRepository) Create(ctx context.Context, activity *model.FortuneActivity) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成抽奖活动ID失败", err)
	}
	activity.ID = id
	if err := r.db.WithContext(ctx).Create(activity).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增抽奖活动失败", err)
	}
	return id, nil
}

func (r *lotteryActivityRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneActivity{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates)
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改抽奖活动失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	return nil
}

func (r *lotteryActivityRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneActivity{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除抽奖活动失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "抽奖活动不存在")
	}
	return nil
}
