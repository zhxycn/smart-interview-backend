package _struct

import "time"

type Session struct {
	SessionID    string    `gorm:"primaryKey;column:session_id"`
	UID          int64     `gorm:"column:uid"`
	IPAddress    string    `gorm:"column:ip_address"`
	UserAgent    string    `gorm:"column:user_agent"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	LastActivity time.Time `gorm:"column:last_activity"`
}

func (Session) TableName() string {
	return "session"
}
