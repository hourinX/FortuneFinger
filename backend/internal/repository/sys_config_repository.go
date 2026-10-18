package repository

import (
	"context"
	"errors"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/snowflake"

	"gorm.io/gorm"
)

type SysConfigRepository interface {
	GetConfigByCode(ctx context.Context, code string) (*model.FortuneSysConfig, error)
	Create(ctx context.Context, config *model.FortuneSysConfig) (uint64, error)
	UpdateConfig(ctx context.Context, id uint64, code, config string, updatedBy uint64) error
}

type sysConfigRepository struct {
	db        *gorm.DB
	snowflake *snowflake.Generator
}

func NewSysConfigRepository(db *gorm.DB, generator *snowflake.Generator) SysConfigRepository {
	return &sysConfigRepository{db: db, snowflake: generator}
}

func (r *sysConfigRepository) GetConfigByCode(ctx context.Context, code string) (*model.FortuneSysConfig, error) {
	var config model.FortuneSysConfig
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.DatabaseError, "查询系统配置失败", err)
	}
	return &config, nil
}

func (r *sysConfigRepository) Create(ctx context.Context, config *model.FortuneSysConfig) (uint64, error) {
	id, err := r.snowflake.Next()
	if err != nil {
		return 0, errorx.Wrap(errorx.InternalServer, "生成系统配置ID失败", err)
	}
	config.ID = id
	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return 0, errorx.Wrap(errorx.DatabaseError, "新增系统配置失败", err)
	}
	return id, nil
}

func (r *sysConfigRepository) UpdateConfig(ctx context.Context, id uint64, code, config string, updatedBy uint64) error {
	result := r.db.WithContext(ctx).Model(&model.FortuneSysConfig{}).
		Where("id = ? AND code = ?", id, code).
		Updates(map[string]any{
			"config":     config,
			"updated_by": updatedBy,
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.DatabaseError, "修改系统配置失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errorx.NotFound, "系统配置不存在")
	}
	return nil
}
