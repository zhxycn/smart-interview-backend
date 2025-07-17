package _struct

import "time"

type Interview struct {
	ID           string    `gorm:"column:id;primary_key"`
	User         int       `gorm:"column:user"`
	Position     string    `gorm:"column:position"`
	Level        string    `gorm:"column:level"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	StartTime    time.Time `gorm:"column:start_time"`
	EndTime      time.Time `gorm:"column:end_time"`
	Conversation string    `gorm:"column:conversation;type:json"`
	Facial       string    `gorm:"column:facial;type:json"`
}

func (Interview) TableName() string {
	return "interview"
}
