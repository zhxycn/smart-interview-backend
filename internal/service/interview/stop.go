package interview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart-interview/internal/database"
	"time"
)

func StopInterview(uid int64, id string) (bool, error) {
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database connection failed")
	}

	rdb := database.GetRedis()
	if rdb == nil {
		return false, errors.New("redis connection failed")
	}

	timestamp := time.Now()

	ctx := context.Background()
	conversationRedisKey := fmt.Sprintf("interview:%s:conversation", id)
	facialRedisKey := fmt.Sprintf("interview:%s:facial", id)

	conversationData, err := rdb.LRange(ctx, conversationRedisKey, 0, -1).Result()
	if err != nil {
		return false, err
	}

	facialData, err := rdb.LRange(ctx, facialRedisKey, 0, -1).Result()
	if err != nil {
		return false, err
	}

	conversationJSON, err := json.Marshal(conversationData)
	if err != nil {
		return false, err
	}

	facialJSON, err := json.Marshal(facialData)
	if err != nil {
		return false, err
	}

	_, err = db.Exec(
		"UPDATE interview SET end_time = ?, conversation = ?, facial = ? WHERE id = ? AND user = ?",
		timestamp, conversationJSON, facialJSON, id, uid,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
