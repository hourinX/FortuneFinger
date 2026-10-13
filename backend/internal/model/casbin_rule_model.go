package model

type CasbinRule struct {
	ID    uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id,string"`
	PType string `gorm:"column:ptype;type:varchar(100);not null" json:"ptype"`
	V0    string `gorm:"column:v0;type:varchar(100);not null" json:"v0"`
	V1    string `gorm:"column:v1;type:varchar(100);not null" json:"v1"`
	V2    string `gorm:"column:v2;type:varchar(100);not null" json:"v2"`
	V3    string `gorm:"column:v3;type:varchar(100);not null" json:"v3"`
	V4    string `gorm:"column:v4;type:varchar(100);not null" json:"v4"`
	V5    string `gorm:"column:v5;type:varchar(100);not null" json:"v5"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
