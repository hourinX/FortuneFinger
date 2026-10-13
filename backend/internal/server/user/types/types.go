package types

//region 公共类型

type IDPath struct {
	ID uint64 `uri:"id" binding:"required"`
}

//endregion

//region 认证出入参

// 用户登录
type LoginReq struct {
	Username          string `json:"username" binding:"required,max=64"`
	Password          string `json:"password" binding:"required,max=255"`
	CaptchaID         string `json:"captcha_id" binding:"required"`
	CaptchaCode       string `json:"captcha" binding:"required"`
	VerificationToken string `json:"verification_token"`
	Remember          bool   `json:"remember"`
}

type LoginResp struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"-"`
	RefreshExpiresAt int64  `json:"-"`
	Remember         bool   `json:"-"`
}

// 刷新登录凭证
type RefreshTokenReq struct {
}

type RefreshTokenResp struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"-"`
	RefreshExpiresAt int64  `json:"-"`
	Remember         bool   `json:"-"`
}

// 退出登录
type LogoutReq struct{}

// 获取验证码
type GetCaptchaReq struct{}

type GetCaptchaResp struct {
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

// 创建增强验证
type CreateEnhancedVerificationReq struct {
	Scene   string `json:"scene" binding:"required,oneof=login"`
	Subject string `json:"subject" binding:"required,max=64"`
}

type CreateEnhancedVerificationResp struct {
	VerificationID string `json:"verification_id"`
	Mode           string `json:"mode"`
	Image          string `json:"image"`
	Prompt         string `json:"prompt"`
	Thumb          string `json:"thumb"`
	ThumbX         int    `json:"thumb_x"`
	ThumbY         int    `json:"thumb_y"`
	ThumbWidth     int    `json:"thumb_width"`
	ThumbHeight    int    `json:"thumb_height"`
	RequiredPoints int    `json:"required_points"`
	ExpiresIn      int64  `json:"expires_in"`
}

// 校验增强验证
type VerifyEnhancedVerificationReq struct {
	VerificationID string                      `json:"verification_id" binding:"required"`
	Scene          string                      `json:"scene" binding:"required,oneof=login"`
	Subject        string                      `json:"subject" binding:"required,max=64"`
	Points         []EnhancedVerificationPoint `json:"points" binding:"required,min=1,max=128,dive"`
}

type EnhancedVerificationPoint struct {
	X int   `json:"x" binding:"min=0"`
	Y int   `json:"y" binding:"min=0"`
	T int64 `json:"t" binding:"min=0"`
}

type VerifyEnhancedVerificationResp struct {
	VerificationToken string `json:"verification_token"`
	ExpiresIn         int64  `json:"expires_in"`
}

//endregion

//region 当前用户信息出入参

// 获取当前用户信息
type GetUserInfoReq struct{}

type GetUserInfoResp struct {
	ID       uint64 `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleCode string `json:"role_code"`
}

//endregion

//region 用户管理出入参

// 获取用户列表
type GetUserListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Username string `form:"username" binding:"max=64"`
	Nickname string `form:"nickname" binding:"max=64"`
	Status   uint8  `form:"status" binding:"omitempty,oneof=1 2"`
	OrgID    uint64 `form:"org_id"`
}

type GetUserListResp struct {
	List  []UserListItem `json:"list"`
	Total int64          `json:"total"`
}

type UserListItem struct {
	ID             uint64 `json:"id,string"`
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	Status         uint8  `json:"status"`
	RoleID         uint64 `json:"role_id,string"`
	RoleName       string `json:"role_name"`
	RoleCode       string `json:"role_code"`
	DepartmentName string `json:"department_name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// 获取用户详情
type GetUserDetailResp struct {
	ID           uint64 `json:"id,string"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	Status       uint8  `json:"status"`
	RoleID       uint64 `json:"role_id,string"`
	RoleName     string `json:"role_name"`
	RoleCode     string `json:"role_code"`
	DepartmentID uint64 `json:"department_id,string"`
}

// 用户新增和修改共用参数
type UserReq struct {
	Username string `json:"username" binding:"required,max=64"`
	Nickname string `json:"nickname" binding:"required,max=64"`
	Avatar   string `json:"avatar" binding:"max=512"`
	Email    string `json:"email" binding:"omitempty,email,max=254"`
	Mobile   string `json:"mobile" binding:"max=32"`
	Status   uint8  `json:"status" binding:"required,oneof=1 2"`
}

// 新增用户
type CreateUserReq struct {
	UserReq
	Password     string `json:"password" binding:"required,min=6,max=64"`
	DepartmentID uint64 `json:"department_id,string" binding:"required"`
}

// 修改用户
type UpdateUserReq struct {
	UserReq
	Password     string `json:"password" binding:"omitempty,min=6,max=64"`
	DepartmentID uint64 `json:"department_id,string"`
}

// 分配用户角色
type AssignUserRoleReq struct {
	RoleID uint64 `json:"role_id,string" binding:"required"`
}

//endregion
