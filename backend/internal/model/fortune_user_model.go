package model

import "time"

type FortuneUser struct {
	ID        uint64    `gorm:"column:id;primaryKey"                            json:"id,string"`
	Username  string    `gorm:"column:username;type:varchar(64);not null"       json:"username"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64);not null"       json:"nickname"`
	Password  string    `gorm:"column:password;type:varchar(255);not null"      json:"-"`
	Avatar    string    `gorm:"column:avatar;type:varchar(512);not null"        json:"avatar"`
	Email     string    `gorm:"column:email;type:varchar(254);not null"         json:"email"`
	Mobile    string    `gorm:"column:mobile;type:varchar(32);not null"         json:"mobile"`
	Status    uint8     `gorm:"column:status;type:tinyint unsigned;not null"    json:"status"`
	CreatedBy uint64    `gorm:"column:created_by;type:bigint unsigned;not null" json:"created_by,string"`
	UpdatedBy uint64    `gorm:"column:updated_by;type:bigint unsigned;not null" json:"updated_by,string"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null"     json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null"     json:"updated_at"`
	DeletedAt uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null" json:"deleted_at"`
}

func (FortuneUser) TableName() string {
	return "fortune_user"
}
