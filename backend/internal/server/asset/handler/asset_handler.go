package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/asset/logic"
	"fortunefinger/backend/internal/server/asset/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	logic logic.AssetLogic
}

func NewAssetHandler(svcCtx *svc.ServiceContext) *AssetHandler {
	return &AssetHandler{logic: logic.NewAssetLogic(svcCtx)}
}

func (h *AssetHandler) UploadHandler(c *gin.Context) {
	var req types.UploadReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "上传参数错误", err))
		return
	}

	resp, err := h.logic.Upload(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, "上传成功", resp)
}

func (h *AssetHandler) GenerateAvatarHandler(c *gin.Context) {
	var req types.GenerateAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "生成头像参数错误", err))
		return
	}

	resp, err := h.logic.GenerateAvatar(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, "头像生成成功", resp)
}
