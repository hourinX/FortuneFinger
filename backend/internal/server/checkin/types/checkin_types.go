package types

type GetCheckinTemplateReq struct{}

type GetCheckinTemplateResp struct {
	Template string `json:"template"`
}

type GetCheckinStateReq struct{}

type GetCheckinStateResp struct {
	Status           string  `json:"status"`
	ServerTime       string  `json:"server_time"`
	CheckinStartTime string  `json:"checkin_start_time"`
	CheckinEndTime   string  `json:"checkin_end_time"`
	CheckedIn        uint64  `json:"checked_in"`
	Total            *uint64 `json:"total"`
}

type GetCheckinQRCodeReq struct{}

type GetCheckinQRCodeResp struct {
	Token        string `json:"token"`
	QRCode       string `json:"qrcode"`
	ExpiresAt    string `json:"expires_at"`
	RefreshAfter uint8  `json:"refresh_after"`
}

type CheckinConfig struct {
	Mode             uint8               `json:"mode"`
	CheckinStartTime string              `json:"checkin_start_time"`
	CheckinEndTime   string              `json:"checkin_end_time"`
	Screen           CheckinScreenConfig `json:"screen"`
}

type CheckinScreenConfig struct {
	Template string `json:"template"`
}
