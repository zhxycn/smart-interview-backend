package _struct

import "time"

type Resume struct {
	ID             string    `gorm:"column:id;primary_key"`
	User           int       `gorm:"column:user"`
	FileName       string    `gorm:"column:file_name"`
	FileData       []byte    `gorm:"column:file_data"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	TargetPosition string    `gorm:"column:target_position"`
	Experience     string    `gorm:"column:experience"`
	Industry       string    `gorm:"column:industry"`
	FocusAreas     string    `gorm:"column:focus_areas"`
	Score          int       `gorm:"column:score"`
	Feedback       string    `gorm:"column:feedback;type:json"`
}

func (Resume) TableName() string {
	return "resume"
}
