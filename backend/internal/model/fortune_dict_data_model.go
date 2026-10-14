package model

import "time"

type FortuneDictData struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"                  json:"id,string"`
	DictCode  string    `gorm:"column:dict_code;type:varchar(100);not null"         json:"dict_code"`
	Label     string    `gorm:"column:label;type:varchar(100);not null"             json:"label"`
	Value     string    `gorm:"column:value;type:varchar(255);not null"             json:"value"`
	Color     string    `gorm:"column:color;type:varchar(32);not null"              json:"color"`
	IsDefault uint8     `gorm:"column:is_default;type:tinyint unsigned;not null"    json:"is_default"`
	Sort      uint      `gorm:"column:sort;type:int unsigned;not null"              json:"sort"`
	Status    uint8     `gorm:"column:status;type:tinyint unsigned;not null"        json:"status"`
	Remark    string    `gorm:"column:remark;type:varchar(255);not null"            json:"remark"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null"     json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"     json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"         json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"         json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"     json:"deleted_at"`
}

func (FortuneDictData) TableName() string {
	return "fortune_dict_data"
}
