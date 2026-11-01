package types

type GetCheckinTemplateReq struct {
	Code string `form:"code" binding:"required,len=8,alphanum"`
}

type GetCheckinTemplateResp struct {
	Template string `json:"template"`
}

type GetCheckinStateReq struct {
	Code string `form:"code" binding:"required,len=8,alphanum"`
}

type GetCheckinStateResp struct {
	Status            string  `json:"status"`
	ServerTime        string  `json:"server_time"`
	ActivityStartTime string  `json:"activity_start_time"`
	ActivityEndTime   string  `json:"activity_end_time"`
	CheckinStartTime  string  `json:"checkin_start_time"`
	CheckinEndTime    string  `json:"checkin_end_time"`
	CheckedIn         uint64  `json:"checked_in"`
	Total             *uint64 `json:"total"`
}

type GetCheckinQRCodeReq struct {
	Code string `form:"code" binding:"required,len=8,alphanum"`
}

type GetCheckinQRCodeResp struct {
	Token        string `json:"token"`
	QRCode       string `json:"qrcode"`
	ExpiresAt    string `json:"expires_at"`
	RefreshAfter uint8  `json:"refresh_after"`
}

type CreateCheckinSessionReq struct {
	Token string `json:"token" binding:"required"`
}

type CreateCheckinSessionResp struct {
	SessionID string                   `json:"session_id"`
	ExpiresAt string                   `json:"expires_at"`
	Activity  CheckinActivityResp      `json:"activity"`
	Config    CheckinSessionConfigResp `json:"config"`
}

type CheckinActivityResp struct {
	Name        string `json:"name"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

type CheckinSessionConfigResp struct {
	Mode           uint8 `json:"mode"`
	RequireName    bool  `json:"require_name"`
	RequireCompany bool  `json:"require_company"`
	AllowAvatar    bool  `json:"allow_avatar"`
	AllowMessage   bool  `json:"allow_message"`
}

type CheckinConfig struct {
	Mode           uint8               `json:"mode"`
	RequireName    bool                `json:"requireName"`
	RequireCompany bool                `json:"requireCompany"`
	AllowAvatar    bool                `json:"allowAvatar"`
	Screen         CheckinScreenConfig `json:"screen"`
}

type CheckinScreenConfig struct {
	Template string `json:"template"`
}
