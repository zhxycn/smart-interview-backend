package main

import (
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"smart-interview/internal/middleware"
	"smart-interview/internal/router"
	"smart-interview/internal/service/user"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()
	middleware.Logger = middleware.NewLogger(cfg)

	// 初始化路由
	r := router.NewRouter()
	loggedRouter := middleware.Logger.HttpMiddleware(r)

	// 初始化数据库连接
	err := database.NewDB(cfg)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to initialize database: %v", err))
		return
	}
	defer database.GetDB().Close()

	// 初始化 Redis 连接
	err = database.NewRedis(cfg)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to initialize Redis: %s", err))
		return
	}

	// 启动会话清理任务
	user.StartSessionCleanup()

	// 启动服务
	middleware.Logger.Log("INFO", fmt.Sprintf("Starting server on port %s", cfg.Port))
	err = http.ListenAndServe(":"+cfg.Port, loggedRouter)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to start server: %v", err))
	}
}
