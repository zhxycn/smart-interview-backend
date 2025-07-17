package interview

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"smart-interview/internal/database"
	"smart-interview/internal/middleware"
	"time"
)

func Register(uid int64, position, level string) (string, error) {
	db := database.GetDB()
	if db == nil {
		return "", errors.New("database connection failed")
	}

	interviewId := uuid.New().String()

	timestamp := time.Now()

	_, err := db.Exec("INSERT INTO interview (id, user, position, level, created_at) VALUES (?, ?, ?, ?, ?)", interviewId, uid, position, level, timestamp)
	if err != nil {
		return "", err
	}

	middleware.Logger.Log("INFO", fmt.Sprintf("Interview registered successfully, ID: %s", interviewId))

	return interviewId, nil
}
