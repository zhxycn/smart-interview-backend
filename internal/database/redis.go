package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"smart-interview/internal/config"
	"smart-interview/internal/middleware"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func NewRedis(cfg *config.Config) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("[REDIS] Failed to connect: %s", err))
		return err
	}
	middleware.Logger.Log("INFO", "[REDIS] Connected successfully")
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}
