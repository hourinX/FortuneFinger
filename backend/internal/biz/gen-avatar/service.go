package genavatar

import (
	"bytes"
	"context"
	"strconv"

	storageupload "fortunefinger/backend/internal/biz/storage-upload"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/svc"

	"github.com/hourinX/gdenticon/v2/pkg/gdenticon"
	"go.uber.org/zap"
)

const avatarScene = "avatar"

type Result struct {
	StorageID uint64
	Path      string
}

type Service struct {
	svcCtx   *svc.ServiceContext
	uploader *storageupload.Service
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		svcCtx:   svcCtx,
		uploader: storageupload.NewService(svcCtx),
	}
}

func (s *Service) GeneratePNG(ctx context.Context, userID uint64) (*Result, error) {
	data, err := renderPNG(userID)
	if err != nil {
		s.svcCtx.Logger.Error("[GenerateAvatar] 生成默认头像失败", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成头像失败", err)
	}

	result, err := s.uploader.Upload(ctx, storageupload.Input{
		Reader:      bytes.NewReader(data),
		Size:        int64(len(data)),
		Name:        "identicon.png",
		Extension:   "png",
		ContentType: "image/png",
		Scene:       avatarScene,
		OperatorID:  userID,
	})
	if err != nil {
		return nil, err
	}

	return &Result{
		StorageID: result.StorageID,
		Path:      result.Path,
	}, nil
}

func renderPNG(userID uint64) ([]byte, error) {
	icon, err := gdenticon.Generate(strconv.FormatUint(userID, 10), gdenticon.GenerateOptions{
		GridSize: gdenticon.DefaultGridSize,
		Mirror:   gdenticon.DefaultMirror,
		Format:   gdenticon.FormatPNG,
	})
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	if err := gdenticon.Write(&content, icon, gdenticon.WriteOptions{}); err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}
