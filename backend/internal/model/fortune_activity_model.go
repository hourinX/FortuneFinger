package model

import "time"

type FortuneActivity struct {
	ID                uint64    `gorm:"column:id;primaryKey"                                 json:"id,string"`
	Code              string    `gorm:"column:code;type:char(8);not null;uniqueIndex"        json:"code"`
	Name              string    `gorm:"column:name;type:varchar(100);not null"               json:"name"`
	OrgID             uint64    `gorm:"column:org_id;type:bigint unsigned;not null"           json:"org_id,string"`
	Venue             string    `gorm:"column:venue;type:varchar(255);not null"               json:"venue"`
	ActivityStartTime time.Time `gorm:"column:activity_start_time;type:datetime(3);not null"  json:"activity_start_time"`
	ActivityEndTime   time.Time `gorm:"column:activity_end_time;type:datetime(3);not null"    json:"activity_end_time"`
	CheckinStartTime  time.Time `gorm:"column:checkin_start_time;type:datetime(3);not null"   json:"checkin_start_time"`
	CheckinEndTime    time.Time `gorm:"column:checkin_end_time;type:datetime(3);not null"     json:"checkin_end_time"`
	CheckinTemplate   string    `gorm:"column:checkin_template;type:varchar(50);not null"     json:"checkin_template"`
	Status            uint8     `gorm:"column:status;type:tinyint unsigned;not null"          json:"status"`
	Remark            string    `gorm:"column:remark;type:varchar(255);not null"              json:"remark"`
	CreatedBy         uint64    `gorm:"column:created_by;type:bigint unsigned;not null"       json:"created_by,string"`
	UpdatedBy         uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"       json:"updated_by,string"`
	CreatedAt         time.Time `gorm:"column:created_at;type:datetime(3);not null"           json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:datetime(3);not null"           json:"updated_at"`
	DeletedAt         uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"       json:"deleted_at"`
}

func (FortuneActivity) TableName() string {
	return "fortune_activity"
}
