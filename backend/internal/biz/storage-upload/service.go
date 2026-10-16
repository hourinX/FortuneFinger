package storageupload

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"path"
	"strings"
	"time"

	"fortunefinger/backend/internal/model"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/repository"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type Input struct {
	Reader      io.Reader
	Size        int64
	Name        string
	Extension   string
	ContentType string
	Scene       string
	OperatorID  uint64
}

type Result struct {
	StorageID uint64
	Path      string
}

type Service struct {
	svcCtx      *svc.ServiceContext
	storageRepo repository.StorageRepository
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{
		svcCtx:      svcCtx,
		storageRepo: repository.NewStorageRepository(svcCtx.GormDB, svcCtx.Snowflake),
	}
}

func (s *Service) Upload(ctx context.Context, input Input) (*Result, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		s.svcCtx.Logger.Error("[Upload] 生成对象存储路径失败", zap.Uint64("user_id", input.OperatorID), zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "生成文件存储路径失败", err)
	}

	fileName := hex.EncodeToString(randomBytes)
	if input.Extension != "" {
		fileName += "." + input.Extension
	}
	objectName := buildObjectName(s.svcCtx.Config.Storage.Prefix, input.Scene, fileName, time.Now())
	hasher := sha256.New()
	if err := s.svcCtx.Storage.PutObject(ctx, objectName, io.TeeReader(input.Reader, hasher), input.Size, input.ContentType); err != nil {
		s.svcCtx.Logger.Error("[Upload] 上传对象存储失败",
			zap.Uint64("user_id", input.OperatorID),
			zap.String("scene", input.Scene),
			zap.String("path", objectName),
			zap.Error(err),
		)
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "上传文件失败，请稍后重试", err)
	}

	storageID, err := s.storageRepo.Create(ctx, &model.FortuneStorage{
		Provider:  s.svcCtx.Config.Storage.Driver,
		Bucket:    s.svcCtx.Config.Storage.Bucket,
		Path:      objectName,
		Name:      input.Name,
		Ext:       input.Extension,
		MIMEType:  input.ContentType,
		Size:      uint64(input.Size),
		Hash:      hex.EncodeToString(hasher.Sum(nil)),
		Scene:     input.Scene,
		CreatedBy: input.OperatorID,
		UpdatedBy: input.OperatorID,
	})
	if err != nil {
		s.svcCtx.Logger.Error("[Upload] 新增对象存储记录失败",
			zap.Uint64("user_id", input.OperatorID),
			zap.String("scene", input.Scene),
			zap.String("path", objectName),
			zap.Error(err),
		)
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if deleteErr := s.svcCtx.Storage.DeleteObject(cleanupCtx, objectName); deleteErr != nil {
			s.svcCtx.Logger.Error("[Upload] 清理未落库的对象失败",
				zap.Uint64("user_id", input.OperatorID),
				zap.String("path", objectName),
				zap.Error(deleteErr),
			)
		}
		return nil, err
	}

	return &Result{
		StorageID: storageID,
		Path:      objectName,
	}, nil
}

func buildObjectName(prefix, scene, fileName string, now time.Time) string {
	value := strings.NewReplacer(
		"{scene}", scene,
		"{yyyyMMdd}", now.Format("20060102"),
		"{yyyyMM}", now.Format("200601"),
		"{yyyy-MM-dd}", now.Format(time.DateOnly),
		"{yyyy-MM}", now.Format("2006-01"),
		"{yyyy}", now.Format("2006"),
		"{MM}", now.Format("01"),
		"{dd}", now.Format("02"),
		"{file}", fileName,
	).Replace(prefix)
	return path.Clean(strings.Trim(value, "/"))
}
