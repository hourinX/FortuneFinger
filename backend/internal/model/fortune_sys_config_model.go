package model

import "time"

type FortuneSysConfig struct {
	ID        uint64    `gorm:"column:id;primaryKey"                            json:"id,string"`
	Code      string    `gorm:"column:code;type:varchar(100);not null"          json:"code"`
	Config    string    `gorm:"column:config;type:text;not null"                json:"config"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null" json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null" json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"      json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"      json:"updated_at"`
}

func (FortuneSysConfig) TableName() string {
	return "fortune_sys_config"
}
