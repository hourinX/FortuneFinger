package model

import "time"

type FortuneMenu struct {
	ID        uint64    `gorm:"column:id;primaryKey"                                 json:"id,string"`
	PID       uint64    `gorm:"column:pid;type:bigint unsigned;not null"             json:"pid,string"`
	Types     uint8     `gorm:"column:types;type:tinyint unsigned;not null"          json:"types"`
	Title     string    `gorm:"column:title;type:varchar(100);not null"              json:"title"`
	Node      string    `gorm:"column:node;type:varchar(255);not null"               json:"node"`
	Route     *string   `gorm:"column:route;type:varchar(255)"                       json:"route"`
	Code      *string   `gorm:"column:code;type:varchar(100)"                        json:"code"`
	Icon      string    `gorm:"column:icon;type:varchar(100);not null"               json:"icon"`
	Remark    string    `gorm:"column:remark;type:varchar(255);not null"             json:"remark"`
	Sort      uint      `gorm:"column:sort;type:int unsigned;not null"               json:"sort"`
	KeepAlive uint8     `gorm:"column:keep_alive;type:tinyint unsigned;not null"     json:"keep_alive"`
	Status    uint8     `gorm:"column:status;type:tinyint unsigned;not null"         json:"status"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null"      json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"      json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"          json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"          json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"      json:"deleted_at"`
}

func (FortuneMenu) TableName() string {
	return "fortune_menu"
}
