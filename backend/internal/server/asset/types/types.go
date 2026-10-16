package types

import "mime/multipart"

//region 资源出入参

type UploadReq struct {
	File  *multipart.FileHeader `form:"file" binding:"required"`
	Scene string                `form:"scene" binding:"required,oneof=avatar prize image video event"`
}

type UploadResp struct {
	StorageID uint64 `json:"storage_id,string"`
}

type GenerateAvatarReq struct{}

type GenerateAvatarResp struct {
	StorageID uint64 `json:"storage_id,string"`
	URL       string `json:"url"`
}

//endregion
