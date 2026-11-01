package model

import "time"

type FortuneCheckinFingerprint struct {
	ID                  uint64    `gorm:"column:id;primaryKey"                                      json:"id,string"`
	ActivityID          uint64    `gorm:"column:activity_id;type:bigint unsigned;not null"          json:"activity_id,string"`
	CheckinID           uint64    `gorm:"column:checkin_id;type:bigint unsigned;not null"           json:"checkin_id,string"`
	IdentityHash        string    `gorm:"column:identity_hash;type:char(64);not null"               json:"identity_hash"`
	SessionHash         string    `gorm:"column:session_hash;type:char(64);not null"                json:"session_hash"`
	BrowserInstanceHash string    `gorm:"column:browser_instance_hash;type:char(64);not null"       json:"browser_instance_hash"`
	FingerprintHash     string    `gorm:"column:fingerprint_hash;type:char(64);not null"            json:"fingerprint_hash"`
	FingerprintVersion  uint16    `gorm:"column:fingerprint_version;type:smallint unsigned;not null" json:"fingerprint_version"`
	FingerprintJSHash   string    `gorm:"column:fingerprintjs_hash;type:char(64);not null"          json:"fingerprintjs_hash"`
	CanvasHash          string    `gorm:"column:canvas_hash;type:char(64);not null"                 json:"canvas_hash"`
	WebGLHash           string    `gorm:"column:webgl_hash;type:char(64);not null"                  json:"webgl_hash"`
	AudioHash           string    `gorm:"column:audio_hash;type:char(64);not null"                  json:"audio_hash"`
	ScreenHash          string    `gorm:"column:screen_hash;type:char(64);not null"                 json:"screen_hash"`
	UserAgentHash       string    `gorm:"column:user_agent_hash;type:char(64);not null"             json:"user_agent_hash"`
	ClientHintsHash     string    `gorm:"column:client_hints_hash;type:char(64);not null"           json:"client_hints_hash"`
	LanguagesHash       string    `gorm:"column:languages_hash;type:char(64);not null"              json:"languages_hash"`
	ScreenWidth         uint      `gorm:"column:screen_width;type:int unsigned;not null"            json:"screen_width"`
	ScreenHeight        uint      `gorm:"column:screen_height;type:int unsigned;not null"           json:"screen_height"`
	AvailWidth          uint      `gorm:"column:avail_width;type:int unsigned;not null"             json:"avail_width"`
	AvailHeight         uint      `gorm:"column:avail_height;type:int unsigned;not null"            json:"avail_height"`
	DevicePixelRatio    float64   `gorm:"column:device_pixel_ratio;type:decimal(6,3) unsigned;not null" json:"device_pixel_ratio"`
	ColorDepth          uint8     `gorm:"column:color_depth;type:tinyint unsigned;not null"          json:"color_depth"`
	Timezone            string    `gorm:"column:timezone;type:varchar(64);not null"                 json:"timezone"`
	Platform            string    `gorm:"column:platform;type:varchar(64);not null"                 json:"platform"`
	Browser             string    `gorm:"column:browser;type:varchar(64);not null"                  json:"browser"`
	HardwareConcurrency uint16    `gorm:"column:hardware_concurrency;type:smallint unsigned;not null" json:"hardware_concurrency"`
	DeviceMemory        float64   `gorm:"column:device_memory;type:decimal(6,2) unsigned;not null"   json:"device_memory"`
	MaxTouchPoints      uint16    `gorm:"column:max_touch_points;type:smallint unsigned;not null"   json:"max_touch_points"`
	TouchSupported      uint8     `gorm:"column:touch_supported;type:tinyint unsigned;not null"     json:"touch_supported"`
	IPHash              string    `gorm:"column:ip_hash;type:char(64);not null"                     json:"ip_hash"`
	Components          string    `gorm:"column:components;type:json;not null"                      json:"components"`
	Confidence          float64   `gorm:"column:confidence;type:decimal(5,4) unsigned;not null"     json:"confidence"`
	SimilarityScore     uint8     `gorm:"column:similarity_score;type:tinyint unsigned;not null"    json:"similarity_score"`
	RiskReason          string    `gorm:"column:risk_reason;type:varchar(500);not null"             json:"risk_reason"`
	Decision            uint8     `gorm:"column:decision;type:tinyint unsigned;not null"            json:"decision"`
	CreatedBy           uint64    `gorm:"column:created_by;type:bigint unsigned;not null"           json:"created_by,string"`
	UpdatedBy           uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"           json:"updated_by,string"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime(3);not null"               json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(3);not null"               json:"updated_at"`
	DeletedAt           uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"           json:"deleted_at"`
}

func (FortuneCheckinFingerprint) TableName() string {
	return "fortune_checkin_fingerprint"
}
