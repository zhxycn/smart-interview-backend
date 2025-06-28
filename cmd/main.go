package main

import (
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"smart-interview/internal/middleware"
	"smart-interview/internal/router"
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

	// 启动服务
	middleware.Logger.Log("INFO", fmt.Sprintf("Starting server on port %s", cfg.Port))
	err = http.ListenAndServe(":"+cfg.Port, loggedRouter)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to start server: %v", err))
	}
}
