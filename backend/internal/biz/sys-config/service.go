package sysconfig

import (
	"context"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/svc"
)

type Config struct {
	ID     uint64
	Code   string
	Config string
}

type SaveInput struct {
	ID         uint64
	Code       string
	Config     string
	OperatorID uint64
}

type Service struct {
	configRepo repository.SysConfigRepository
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		configRepo: repository.NewSysConfigRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (s *Service) SaveConfig(ctx context.Context, input SaveInput) (uint64, error) {
	if input.ID == 0 {
		existing, err := s.configRepo.GetConfigByCode(ctx, input.Code)
		if err != nil {
			return 0, err
		}
		if existing != nil {
			return 0, errorx.New(errorx.Conflict, "系统配置已存在，请刷新后重试")
		}
		return s.configRepo.Create(ctx, &model.FortuneSysConfig{
			Code:      input.Code,
			Config:    input.Config,
			CreatedBy: input.OperatorID,
			UpdatedBy: input.OperatorID,
		})
	}
	if err := s.configRepo.UpdateConfig(ctx, input.ID, input.Code, input.Config, input.OperatorID); err != nil {
		return 0, err
	}
	return input.ID, nil
}

func (s *Service) GetConfigByCode(ctx context.Context, code string) (*Config, error) {
	config, err := s.configRepo.GetConfigByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &Config{
		ID:     config.ID,
		Code:   config.Code,
		Config: config.Config,
	}, nil
}
