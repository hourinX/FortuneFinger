package model

import "time"

type FortuneOrg struct {
	ID        uint64    `gorm:"column:id;primaryKey"                            json:"id,string"`
	ParentID  uint64    `gorm:"column:parent_id;type:bigint unsigned;not null"  json:"parent_id,string"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"          json:"name"`
	Code      string    `gorm:"column:code;type:varchar(100);not null"          json:"code"`
	Types     uint8     `gorm:"column:types;type:tinyint unsigned;not null"     json:"types"`
	Sort      uint      `gorm:"column:sort;type:int unsigned;not null"          json:"sort"`
	Status    uint8     `gorm:"column:status;type:tinyint unsigned;not null"    json:"status"`
	Remark    string    `gorm:"column:remark;type:varchar(255);not null"        json:"remark"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null" json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null" json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"     json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"     json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null" json:"deleted_at"`
}

func (FortuneOrg) TableName() string {
	return "fortune_org"
}
