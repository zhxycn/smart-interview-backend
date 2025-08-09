package _struct

import "time"

type Suggestion struct {
	User      int64     `gorm:"column:user;primary_key"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Suggest   string    `gorm:"column:suggestion;type:json"`
}

func (Suggestion) TableName() string {
	return "suggestion"
}
