package model

import "time"

type FortuneRound struct {
	ID          uint64    `gorm:"column:id;primaryKey"                          json:"id,string"`
	ActivityID  uint64    `gorm:"column:activity_id;type:bigint unsigned;not null" json:"activity_id,string"`
	PrizeID     uint64    `gorm:"column:prize_id;type:bigint unsigned;not null"  json:"prize_id,string"`
	Name        string    `gorm:"column:name;type:varchar(100);not null"         json:"name"`
	DrawDate    time.Time `gorm:"column:draw_date;type:date;not null"            json:"draw_date"`
	DrawTime    string    `gorm:"column:draw_time;type:time;not null"            json:"draw_time"`
	WinnerCount uint      `gorm:"column:winner_count;type:int unsigned;not null" json:"winner_count"`
	CreatedBy   uint64    `gorm:"column:created_by;type:bigint unsigned;not null" json:"created_by,string"`
	UpdatedBy   uint64    `gorm:"column:updated_by;type:bigint unsigned;not null" json:"updated_by,string"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3);not null"     json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(3);not null"     json:"updated_at"`
	DeletedAt   uint64    `gorm:"column:deleted_at;type:bigint unsigned;not null" json:"deleted_at"`
}

func (FortuneRound) TableName() string {
	return "fortune_round"
}
