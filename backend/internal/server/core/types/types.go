package types

//region 公共类型

type IDPath struct {
	ID uint64 `uri:"id" binding:"required"`
}

//endregion

//region 组织管理出入参

// 获取组织列表
type GetOrgListReq struct {
	Name   string `form:"name" binding:"max=100"`
	Code   string `form:"code" binding:"max=100"`
	Types  uint8  `form:"types" binding:"omitempty,oneof=1 2 3"`
	Status uint8  `form:"status" binding:"omitempty,oneof=1 2"`
}

type GetOrgListResp struct {
	List []OrgListItem `json:"list"`
}

type OrgListItem struct {
	ID        uint64        `json:"id,string"`
	ParentID  uint64        `json:"parent_id,string"`
	Name      string        `json:"name"`
	Code      string        `json:"code"`
	Types     uint8         `json:"types"`
	Sort      uint          `json:"sort"`
	Status    uint8         `json:"status"`
	Remark    string        `json:"remark"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Children  []OrgListItem `json:"children,omitempty"`
}

// 获取组织详情
type GetOrgDetailResp struct {
	ID       uint64 `json:"id,string"`
	ParentID uint64 `json:"parent_id,string"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Types    uint8  `json:"types"`
	Sort     uint   `json:"sort"`
	Status   uint8  `json:"status"`
	Remark   string `json:"remark"`
}

// 组织新增和修改共用参数
type OrgReq struct {
	ParentID uint64 `json:"parent_id,string"`
	Name     string `json:"name" binding:"required,max=100"`
	Code     string `json:"code" binding:"required,max=100"`
	Types    uint8  `json:"types" binding:"required,oneof=1 2 3"`
	Sort     uint   `json:"sort"`
	Status   uint8  `json:"status" binding:"required,oneof=1 2"`
	Remark   string `json:"remark" binding:"max=255"`
}

type CreateOrgReq struct {
	OrgReq
}

type UpdateOrgReq struct {
	OrgReq
}

// 获取组织选项
type GetOrgOptionsReq struct{}

type GetOrgOptionsResp struct {
	List []OrgOption `json:"list"`
}

type OrgOption struct {
	ID       uint64      `json:"id,string"`
	Name     string      `json:"name"`
	Children []OrgOption `json:"children,omitempty"`
}

//endregion

//region 菜单出入参

// 获取菜单列表
type GetMenuListReq struct {
	Title  string `form:"title" binding:"max=100"`
	Types  uint8  `form:"types" binding:"omitempty,oneof=1 2 3"`
	Status uint8  `form:"status" binding:"omitempty,oneof=1 2"`
}

type GetMenuListResp struct {
	List []MenuListItem `json:"list"`
}

type MenuListItem struct {
	ID        uint64         `json:"id,string"`
	PID       uint64         `json:"pid,string"`
	Types     uint8          `json:"types"`
	Title     string         `json:"title"`
	Node      string         `json:"node"`
	Route     *string        `json:"route"`
	Code      *string        `json:"code"`
	Icon      string         `json:"icon"`
	Remark    string         `json:"remark"`
	Sort      uint           `json:"sort"`
	KeepAlive uint8          `json:"keep_alive"`
	Status    uint8          `json:"status"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Children  []MenuListItem `json:"children,omitempty"`
}

// 获取菜单详情
type GetMenuDetailResp struct {
	ID        uint64  `json:"id,string"`
	PID       uint64  `json:"pid,string"`
	Types     uint8   `json:"types"`
	Title     string  `json:"title"`
	Node      string  `json:"node"`
	Route     *string `json:"route"`
	Code      *string `json:"code"`
	Icon      string  `json:"icon"`
	Remark    string  `json:"remark"`
	Sort      uint    `json:"sort"`
	KeepAlive uint8   `json:"keep_alive"`
	Status    uint8   `json:"status"`
}

// 菜单新增和修改共用参数
type MenuReq struct {
	PID       uint64 `json:"pid,string"`
	Types     uint8  `json:"types" binding:"required,oneof=1 2 3"`
	Title     string `json:"title" binding:"required,max=100"`
	Node      string `json:"node" binding:"required,max=255"`
	Route     string `json:"route" binding:"max=255"`
	Code      string `json:"code" binding:"max=100"`
	Icon      string `json:"icon" binding:"max=100"`
	Remark    string `json:"remark" binding:"max=255"`
	Sort      uint   `json:"sort"`
	KeepAlive uint8  `json:"keep_alive" binding:"required,oneof=1 2"`
	Status    uint8  `json:"status" binding:"required,oneof=1 2"`
}

// 新增菜单
type CreateMenuReq struct {
	MenuReq
}

// 修改菜单
type UpdateMenuReq struct {
	MenuReq
}

// 获取当前用户路由
type GetRoutersReq struct{}

type GetRoutersResp struct {
	List []RouterItem `json:"list"`
}

type RouterItem struct {
	ID        uint64       `json:"id,string"`
	Types     uint8        `json:"types"`
	Title     string       `json:"title"`
	Node      string       `json:"node"`
	Route     *string      `json:"route"`
	Icon      string       `json:"icon"`
	KeepAlive uint8        `json:"keep_alive"`
	Children  []RouterItem `json:"children,omitempty"`
}

// 获取角色权限菜单树
type GetPermissionTreeResp struct {
	List        []PermissionTreeItem `json:"list"`
	SelectedIDs []string             `json:"selected_ids"`
}

type PermissionTreeItem struct {
	ID       uint64               `json:"id,string"`
	Name     string               `json:"name"`
	Code     string               `json:"code"`
	Children []PermissionTreeItem `json:"children,omitempty"`
}

//endregion

//region 角色出入参

// 获取角色列表
type GetRoleListReq struct {
	Name   string `form:"name" binding:"max=100"`
	Code   string `form:"code" binding:"max=100"`
	Status uint8  `form:"status" binding:"omitempty,oneof=1 2"`
}

type GetRoleListResp struct {
	List []RoleListItem `json:"list"`
}

type RoleListItem struct {
	ID        uint64         `json:"id,string"`
	ParentID  uint64         `json:"parent_id,string"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Status    uint8          `json:"status"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Children  []RoleListItem `json:"children,omitempty"`
}

// 获取角色详情
type GetRoleDetailResp struct {
	ID       uint64 `json:"id,string"`
	ParentID uint64 `json:"parent_id,string"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Status   uint8  `json:"status"`
}

// 角色新增和修改共用参数
type RoleReq struct {
	ParentID uint64 `json:"parent_id,string"`
	Name     string `json:"name" binding:"required,max=100"`
	Status   uint8  `json:"status" binding:"required,oneof=1 2"`
}

// 新增角色
type CreateRoleReq struct {
	RoleReq
	Code string `json:"code" binding:"required,max=100"`
}

// 修改角色
type UpdateRoleReq struct {
	RoleReq
}

// 获取角色选项
type GetRoleOptionsReq struct{}

type GetRoleOptionsResp []RoleOption

type RoleOption struct {
	Key      uint64       `json:"key,string"`
	Label    string       `json:"label"`
	Code     string       `json:"code"`
	Status   uint8        `json:"status"`
	Children []RoleOption `json:"children,omitempty"`
}

// 修改角色权限
type UpdateRolePermissionReq struct {
	MenuIDs []string `json:"menu_ids" binding:"max=10000,dive,numeric"`
}

//endregion
