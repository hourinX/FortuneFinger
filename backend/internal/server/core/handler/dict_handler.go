package handler

import (
	"fortunefinger/backend/internal/pkg/errorx"
	"fortunefinger/backend/internal/pkg/httpx"
	"fortunefinger/backend/internal/server/core/logic"
	"fortunefinger/backend/internal/server/core/types"
	"fortunefinger/backend/internal/svc"

	"github.com/gin-gonic/gin"
)

type DictHandler struct {
	logic logic.DictLogic
}

func NewDictHandler(svcCtx *svc.ServiceContext) *DictHandler {
	return &DictHandler{logic: logic.NewDictLogic(svcCtx)}
}

func (h *DictHandler) GetDictTypeListHandler(c *gin.Context) {
	var req types.GetDictTypeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型列表参数错误", err))
		return
	}

	resp, err := h.logic.GetDictTypeList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *DictHandler) GetDictTypeDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型ID参数错误", err))
		return
	}

	resp, err := h.logic.GetDictTypeDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *DictHandler) CreateDictTypeHandler(c *gin.Context) {
	var req types.CreateDictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型参数错误", err))
		return
	}

	if err := h.logic.CreateDictType(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *DictHandler) UpdateDictTypeHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型ID参数错误", err))
		return
	}

	var req types.UpdateDictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型参数错误", err))
		return
	}

	if err := h.logic.UpdateDictType(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *DictHandler) DeleteDictTypeHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典类型ID参数错误", err))
		return
	}

	if err := h.logic.DeleteDictType(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *DictHandler) GetDictDataListHandler(c *gin.Context) {
	var req types.GetDictDataListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据列表参数错误", err))
		return
	}

	resp, err := h.logic.GetDictDataList(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *DictHandler) GetDictDataDetailHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据ID参数错误", err))
		return
	}

	resp, err := h.logic.GetDictDataDetail(c.Request.Context(), uri.ID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}

func (h *DictHandler) CreateDictDataHandler(c *gin.Context) {
	var req types.CreateDictDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据参数错误", err))
		return
	}

	if err := h.logic.CreateDictData(c.Request.Context(), &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.CreateSuccess)
}

func (h *DictHandler) UpdateDictDataHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据ID参数错误", err))
		return
	}

	var req types.UpdateDictDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据参数错误", err))
		return
	}

	if err := h.logic.UpdateDictData(c.Request.Context(), uri.ID, &req); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.UpdateSuccess)
}

func (h *DictHandler) DeleteDictDataHandler(c *gin.Context) {
	var uri types.IDPath
	if err := c.ShouldBindUri(&uri); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典数据ID参数错误", err))
		return
	}

	if err := h.logic.DeleteDictData(c.Request.Context(), uri.ID); err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.DeleteSuccess)
}

func (h *DictHandler) GetDictOptionsHandler(c *gin.Context) {
	var req types.GetDictOptionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.Error(c, errorx.Wrap(errorx.InvalidParameter, "字典选项参数错误", err))
		return
	}

	resp, err := h.logic.GetDictOptions(c.Request.Context(), &req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.Ok(c, httpx.QuerySuccess, resp)
}
