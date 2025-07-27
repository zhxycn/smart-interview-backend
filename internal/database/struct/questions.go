package _struct

type Questions struct {
	QID        string `gorm:"column:qid;primary_key"`
	ID         string `gorm:"column:id"`
	Question   string `gorm:"column:question"`
	Answer     string `gorm:"column:answer"`
	Difficulty string `gorm:"column:difficulty"`
	UserAnswer string `gorm:"column:user_answer"`
	Feedback   string `gorm:"column:feedback;type:json"`
}

func (Questions) TableName() string {
	return "questions"
}
