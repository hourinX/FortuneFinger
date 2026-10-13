package model

import "time"

type FortuneRole struct {
	ID        uint64    `gorm:"column:id;primaryKey"                                 json:"id,string"`
	ParentID  uint64    `gorm:"column:parent_id;type:bigint unsigned;not null"       json:"parent_id,string"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"               json:"name"`
	Code      string    `gorm:"column:code;type:varchar(100);not null"               json:"code"`
	Status    uint8     `gorm:"column:status;type:tinyint unsigned;not null"         json:"status"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null"      json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"      json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"          json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"          json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"      json:"deleted_at"`
}

func (FortuneRole) TableName() string {
	return "fortune_role"
}
