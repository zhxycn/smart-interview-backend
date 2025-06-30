package database

type User struct {
	UID int `gorm:"primaryKey;column:uid"`
}

func (User) TableName() string {
	return "user"
}
