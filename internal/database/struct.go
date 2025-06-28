package database

type DB struct {
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
}

func (DB) TableName() string {
	return "db"
}
