package model

import "time"

type FortuneStorage struct {
	ID        uint64    `gorm:"column:id;primaryKey"                             json:"id,string"`
	Provider  string    `gorm:"column:provider;type:varchar(20);not null"        json:"provider"`
	Bucket    string    `gorm:"column:bucket;type:varchar(100);not null"         json:"bucket"`
	Path      string    `gorm:"column:path;type:varchar(512);not null"           json:"path"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"           json:"name"`
	Ext       string    `gorm:"column:ext;type:varchar(20);not null"             json:"ext"`
	MIMEType  string    `gorm:"column:mime_type;type:varchar(100);not null"      json:"mime_type"`
	Size      uint64    `gorm:"column:size;type:bigint unsigned;not null"        json:"size"`
	Hash      string    `gorm:"column:hash;type:char(64);not null"               json:"hash"`
	Scene     string    `gorm:"column:scene;type:varchar(32);not null"           json:"scene"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null"  json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"  json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"       json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"       json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"  json:"deleted_at"`
}

func (FortuneStorage) TableName() string {
	return "fortune_storage"
}
