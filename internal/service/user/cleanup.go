package user

import (
	"fmt"
	"smart-interview/internal/middleware"
	"time"
)

func StartSessionCleanup() {
	go func() {
		// 每小时清理一次过期会话
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			count := DeleteExpiredSessions()
			if count > 0 {
				middleware.Logger.Log("INFO", fmt.Sprintf("Session cleanup: removed %d expired sessions", count))
			}
		}
	}()

	count := DeleteExpiredSessions() // 立即执行
	if count > 0 {
		middleware.Logger.Log("INFO", fmt.Sprintf("Initial session cleanup: removed %d expired sessions", count))
	}

	middleware.Logger.Log("INFO", "Session cleanup task started")
}
