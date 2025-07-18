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
	Score        int       `gorm:"column:score;default:-1"`
	Communicate  int       `gorm:"column:communicate;default:-1"`
	Specialized  int       `gorm:"column:specialized;default:-1"`
	Expression   int       `gorm:"column:expression;default:-1"`
	Strain       int       `gorm:"column:strain;default:-1"`
	Appearance   int       `gorm:"column:appearance;default:-1"`
	Appraise     string    `gorm:"column:appraise;type:json"`
}

func (Interview) TableName() string {
	return "interview"
}
