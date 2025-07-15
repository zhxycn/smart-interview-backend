package database

type User struct {
	UID      int64  `gorm:"primaryKey;column:uid"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Name     string `gorm:"column:name"`
	Job      string `gorm:"column:job;default:user"`
}

func (User) TableName() string {
	return "user"
}
