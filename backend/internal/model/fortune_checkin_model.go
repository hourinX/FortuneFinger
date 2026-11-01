package model

import "time"

type FortuneCheckin struct {
	ID           uint64    `gorm:"column:id;primaryKey"                            json:"id,string"`
	ActivityID   uint64    `gorm:"column:activity_id;type:bigint unsigned;not null" json:"activity_id,string"`
	OrgID        uint64    `gorm:"column:org_id;type:bigint unsigned;not null"      json:"org_id,string"`
	RosterID     uint64    `gorm:"column:roster_id;type:bigint unsigned;not null"   json:"roster_id,string"`
	Name         string    `gorm:"column:name;type:varchar(100);not null"           json:"name"`
	Mobile       string    `gorm:"column:mobile;type:varchar(32);not null"          json:"mobile"`
	Email        string    `gorm:"column:email;type:varchar(254);not null"          json:"email"`
	IdentityHash string    `gorm:"column:identity_hash;type:char(64);not null"      json:"identity_hash"`
	Avatar       uint64    `gorm:"column:avatar;type:bigint unsigned;not null"      json:"avatar,string"`
	CheckedInAt  time.Time `gorm:"column:checked_in_at;type:datetime(3);not null"   json:"checked_in_at"`
	CreatedBy    uint64    `gorm:"column:created_by;type:bigint unsigned;not null"  json:"created_by,string"`
	UpdatedBy    uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"  json:"updated_by,string"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime(3);not null"      json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime(3);not null"      json:"updated_at"`
	DeletedAt    uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"  json:"deleted_at"`
}

func (FortuneCheckin) TableName() string {
	return "fortune_checkin"
}
