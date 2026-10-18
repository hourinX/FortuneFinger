package types

//region 奖品池出入参

type IDPath struct {
	ID uint64 `uri:"id" binding:"required"`
}

type GetPoolListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name" binding:"max=100"`
	Level    string `form:"level" binding:"max=100"`
	Types    uint8  `form:"types" binding:"omitempty,min=1"`
	Status   uint8  `form:"status" binding:"omitempty,oneof=1 2"`
}

type GetPoolListResp struct {
	List  []PoolListItem `json:"list"`
	Total int64          `json:"total"`
}

type PoolListItem struct {
	ID       uint64 `json:"id,string"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	OrgName  string `json:"org_name"`
	Level    string `json:"level"`
	Types    uint8  `json:"types"`
	Quantity uint   `json:"quantity"`
	Status   uint8  `json:"status"`
}

type GetPoolDetailResp struct {
	ID        uint64 `json:"id,string"`
	Name      string `json:"name"`
	StorageID uint64 `json:"storage_id,string"`
	ImageURL  string `json:"image_url"`
	OrgID     uint64 `json:"org_id,string"`
	Level     string `json:"level"`
	Types     uint8  `json:"types"`
	Quantity  uint   `json:"quantity"`
	Status    uint8  `json:"status"`
	Sort      uint   `json:"sort"`
	Remark    string `json:"remark"`
}

type PoolReq struct {
	Name      string `json:"name" binding:"required,max=100"`
	StorageID uint64 `json:"storage_id,string"`
	OrgID     uint64 `json:"org_id,string" binding:"required"`
	Level     string `json:"level" binding:"required,max=100"`
	Types     uint8  `json:"types" binding:"required,min=1"`
	Quantity  uint   `json:"quantity" binding:"required,min=1"`
	Status    uint8  `json:"status" binding:"required,oneof=1 2"`
	Sort      uint   `json:"sort"`
	Remark    string `json:"remark" binding:"max=255"`
}

type CreatePoolReq struct {
	PoolReq
}

type UpdatePoolReq struct {
	PoolReq
}

//endregion
