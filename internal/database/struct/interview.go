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
	Score        int       `gorm:"column:score;default:0"`
	Communicate  int       `gorm:"column:communicate;default:0"`
	Specialized  int       `gorm:"column:specialized;default:0"`
	Expression   int       `gorm:"column:expression;default:0"`
	Strain       int       `gorm:"column:strain;default:0"`
	Appearance   int       `gorm:"column:appearance;default:0"`
	Appraise     string    `gorm:"column:appraise;type:json"`
	Process      string    `gorm:"column:process;type:json"`
	Problem      string    `gorm:"column:problem;type:json"`
	Suggestion   string    `gorm:"column:suggestion;type:json"`
}

func (Interview) TableName() string {
	return "interview"
}
