package logic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	genavatar "fortunefinger/backend/internal/biz/gen-avatar"
	"fortunefinger/backend/internal/biz/operator"
	storageupload "fortunefinger/backend/internal/biz/storage-upload"
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/server/asset/types"
	"fortunefinger/backend/internal/svc"

	"go.uber.org/zap"
)

type AssetLogic interface {
	Upload(ctx context.Context, req *types.UploadReq) (*types.UploadResp, error)
	GenerateAvatar(ctx context.Context, req *types.GenerateAvatarReq) (*types.GenerateAvatarResp, error)
}

type assetLogic struct {
	svcCtx          *svc.ServiceContext
	uploader        *storageupload.Service
	avatarGenerator *genavatar.Service
}

func NewAssetLogic(svcCtx *svc.ServiceContext) AssetLogic {
	return &assetLogic{
		svcCtx:          svcCtx,
		uploader:        storageupload.NewService(svcCtx),
		avatarGenerator: genavatar.NewService(svcCtx),
	}
}

func (l *assetLogic) Upload(ctx context.Context, req *types.UploadReq) (*types.UploadResp, error) {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.File.Size <= 0 {
		err = errorx.New(errorx.InvalidParameter, "上传文件不能为空")
		return nil, err
	}
	if req.Scene == "avatar" && req.File.Size > 5<<20 {
		err = errorx.New(errorx.InvalidParameter, "头像文件不能超过5MB")
		return nil, err
	}

	if len([]rune(req.File.Filename)) > 255 {
		err = errorx.New(errorx.InvalidParameter, "文件名不能超过255个字符")
		return nil, err
	}

	file, err := req.File.Open()
	if err != nil {
		l.svcCtx.Logger.Error("[AssetLogic] 打开上传文件失败", zap.Uint64("user_id", operatorID), zap.Error(err))
		return nil, errorx.Wrap(errorx.InternalServer, "读取上传文件失败", err)
	}
	defer file.Close()

	header := make([]byte, 512)
	n, readErr := file.Read(header)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		l.svcCtx.Logger.Error("[AssetLogic] 读取上传文件失败", zap.Uint64("user_id", operatorID), zap.Error(readErr))
		return nil, errorx.Wrap(errorx.InternalServer, "读取上传文件失败", readErr)
	}

	if n == 0 {
		err = errorx.New(errorx.InvalidParameter, "上传文件不能为空")
		return nil, err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, errorx.Wrap(errorx.InternalServer, "处理上传文件失败", err)
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(req.File.Filename)), ".")
	if len(ext) > 20 {
		err = errorx.New(errorx.InvalidParameter, "文件扩展名不能超过20个字符")
		return nil, err
	}

	mimeType := http.DetectContentType(header[:n])
	if req.Scene == "avatar" && !strings.HasPrefix(mimeType, "image/") {
		err = errorx.New(errorx.InvalidParameter, "头像只能上传图片文件")
		return nil, err
	}

	result, err := l.uploader.Upload(ctx, storageupload.Input{
		Reader:      file,
		Size:        req.File.Size,
		Name:        req.File.Filename,
		Extension:   ext,
		ContentType: mimeType,
		Scene:       req.Scene,
		OperatorID:  operatorID,
	})
	if err != nil {
		return nil, err
	}

	return &types.UploadResp{StorageID: result.StorageID}, nil
}

func (l *assetLogic) GenerateAvatar(ctx context.Context, _ *types.GenerateAvatarReq) (*types.GenerateAvatarResp, error) {
	operatorID, err := operator.UserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.avatarGenerator.GeneratePNG(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	url, err := l.svcCtx.Storage.PresignGetObject(ctx, result.Path, time.Hour)
	if err != nil {
		l.svcCtx.Logger.Error("[AssetLogic] 生成头像访问地址失败", zap.Uint64("storage_id", result.StorageID), zap.Error(err))
		return nil, errorx.Wrap(errorx.ServiceUnavailable, "加载生成头像失败", err)
	}

	return &types.GenerateAvatarResp{
		StorageID: result.StorageID,
		URL:       url,
	}, nil
}
