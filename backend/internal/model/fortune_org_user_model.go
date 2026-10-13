package model

import "time"

type FortuneOrgUser struct {
	ID           uint64    `gorm:"column:id;primaryKey"                               json:"id,string"`
	UserID       uint64    `gorm:"column:user_id;type:bigint unsigned;not null"       json:"user_id,string"`
	CompanyID    uint64    `gorm:"column:company_id;type:bigint unsigned;not null"    json:"company_id,string"`
	DepartmentID uint64    `gorm:"column:department_id;type:bigint unsigned;not null" json:"department_id,string"`
	RelationType uint8     `gorm:"column:relation_type;type:tinyint unsigned;not null" json:"relation_type"`
	CreatedBy    uint64    `gorm:"column:created_by;type:bigint unsigned;not null"    json:"created_by,string"`
	UpdatedBy    uint64    `gorm:"column:updated_by;type:bigint unsigned;not null"    json:"updated_by,string"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime(3);not null"         json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime(3);not null"         json:"updated_at"`
	DeletedAt    uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null"    json:"deleted_at"`
}

func (FortuneOrgUser) TableName() string {
	return "fortune_org_user"
}
