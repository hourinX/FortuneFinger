package types

type GetCheckinTemplateReq struct{}

type GetCheckinTemplateResp struct {
	Template string `json:"template"`
}

type CheckinConfig struct {
	Screen CheckinScreenConfig `json:"screen"`
}

type CheckinScreenConfig struct {
	Template string `json:"template"`
}
