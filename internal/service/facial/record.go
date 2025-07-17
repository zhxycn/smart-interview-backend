package facial

import (
	"context"
	"encoding/json"
	"fmt"
	"smart-interview/internal/database"
	"time"
)

func Record(id string, data interface{}) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("interview ID cannot be empty")
	}

	rdb := database.GetRedis()
	if rdb == nil {
		return false, fmt.Errorf("redis connection failed")
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("interview:%s:facial", id)

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
