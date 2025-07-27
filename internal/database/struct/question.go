package _struct

import "time"

type Question struct {
	ID        string    `gorm:"column:id;primary_key"`
	User      int       `gorm:"column:user"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Position  string    `gorm:"column:position"`
	Knowledge string    `gorm:"column:knowledge;type:json"`
	Count     int       `gorm:"column:count"`
	Feedback  string    `gorm:"column:feedback;type:json"`
}

func (Question) TableName() string {
	return "question"
}
