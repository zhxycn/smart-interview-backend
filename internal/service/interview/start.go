package interview

import (
	"errors"
	"smart-interview/internal/database"
	"time"
)

func StartInterview(uid int64, id string) (bool, error) {
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database connection failed")
	}

	timestamp := time.Now()

	_, err := db.Exec("UPDATE interview SET start_time = ? WHERE id = ? AND user = ?", timestamp, id, uid)
	if err != nil {
		return false, err
	}

	return true, nil
}
