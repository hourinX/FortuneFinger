package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type StorageRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.FortuneStorage, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneStorage, error)
	Create(ctx context.Context, storage *model.FortuneStorage) (uint64, error)
}

type storageRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewStorageRepository(db *gorm.DB, generator *snowflake.Generator) StorageRepository {
	return &storageRepository{db: db, snowflake: generator}
}

func (r *storageRepository) FindByID(ctx context.Context, id uint64) (*model.FortuneStorage, error) {
	var storage model.FortuneStorage
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, 0).First(&storage).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询对象存储记录失败", err)
	}
	return &storage, nil
}

func (r *storageRepository) FindByIDs(ctx context.Context, ids []uint64) ([]model.FortuneStorage, error) {
	if len(ids) == 0 {
		return []model.FortuneStorage{}, nil
	}
	var storages []model.FortuneStorage
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at = ?", ids, 0).Find(&storages).Error; err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "批量查询对象存储记录失败", err)
	}
	return storages, nil
}

func (r *storageRepository) Create(ctx context.Context, storage *model.FortuneStorage) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成存储ID失败", err)
	}
	storage.ID = id
	if err := r.db.WithContext(ctx).Create(storage).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增对象存储记录失败", err)
	}
	return id, nil
}
