package interview

import (
	"context"
	"encoding/json"
	"fmt"
	"smart-interview/internal/database"
	"time"
)

func Record(id, text, role string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("interview ID cannot be empty")
	}

	rdb := database.GetRedis()
	if rdb == nil {
		return false, fmt.Errorf("redis connection failed")
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

func GetRecord(id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("interview ID cannot be empty")
	}

	rdb := database.GetRedis()
	if rdb == nil {
		return "", fmt.Errorf("redis connection failed")
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("interview:%s:conversation", id)

	data, err := rdb.LRange(ctx, redisKey, 0, -1).Result()
	if err != nil {
		return "", err
	}

	var records []map[string]string
	for _, item := range data {
		var record map[string]string
		if err = json.Unmarshal([]byte(item), &record); err != nil {
			return "", err
		}
		records = append(records, record)
	}

	jsonRecords, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	return string(jsonRecords), nil
}
