package interview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart-interview/internal/database"
	"time"
)

func Record(id, text, role string) (bool, error) {
	rdb := database.GetRedis()
	if rdb == nil {
		return false, errors.New("redis connection failed")
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("interview:%s:conversation", id)

	data := map[string]string{
		"role":    role,
		"content": text,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	if err = rdb.RPush(ctx, redisKey, jsonData).Err(); err != nil {
		return false, err
	}
	if err = rdb.Expire(ctx, redisKey, time.Hour).Err(); err != nil {
		return false, err
	}

	return true, nil
}
