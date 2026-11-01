package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type RoundRepository interface {
	ListByActivityIDs(ctx context.Context, activityIDs []uint64) ([]model.FortuneRound, error)
	FindByID(ctx context.Context, id uint64) (*model.FortuneRound, error)
	Create(ctx context.Context, round *model.FortuneRound) (uint64, error)
	Update(ctx context.Context, id uint64, updates map[string]any) error
	Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error
}

type roundRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewRoundRepository(db *gorm.DB, generator *snowflake.Generator) RoundRepository {
	return &roundRepository{db: db, snowflake: generator}
}

func (r *roundRepository) ListByActivityIDs(ctx context.Context, activityIDs []uint64) ([]model.FortuneRound, error) {
	if len(activityIDs) == 0 {
		return []model.FortuneRound{}, nil
	}
	var rounds []model.FortuneRound
	if err := r.db.WithContext(ctx).
		Where("activity_id IN ? AND deleted_at = ?", activityIDs, 0).
		Order("activity_id ASC, draw_date ASC, draw_time ASC, id ASC").
		Find(&rounds).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询抽奖轮次列表失败", err)
	}
	return rounds, nil
}

func (r *roundRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneRound, error) {
	var round model.FortuneRound
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&round).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询抽奖轮次失败", err)
	}
	return &round, nil
}

func (r *roundRepository) Create(ctx context.Context, round *model.FortuneRound) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成抽奖轮次ID失败", err)
	}
	round.ID = id
	if err := r.db.WithContext(ctx).Create(round).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增抽奖轮次失败", err)
	}
	return id, nil
}

func (r *roundRepository) Update(ctx context.Context, id uint64, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneRound{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(updates)
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改抽奖轮次失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "抽奖轮次不存在")
	}
	return nil
}

func (r *roundRepository) Delete(ctx context.Context, id, deletedAt, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneRound{}).
		Where("id = ? AND deleted_at = ?", id, 0).
		Updates(map[string]any{
			"deleted_at": deletedAt,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "删除抽奖轮次失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "抽奖轮次不存在")
	}
	return nil
}
