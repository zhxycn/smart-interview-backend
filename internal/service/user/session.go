package user

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"smart-interview/internal/database"
	"smart-interview/internal/database/struct"
	"smart-interview/internal/middleware"
	"time"
)

func GenerateSessionID(uid int64) string {
	timestamp := time.Now().Unix()
	nBig, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to generate random number: %v", err))
		nBig = big.NewInt(int64(time.Now().Nanosecond() % 9000))
	}
	randomNumber := nBig.Int64() + 1000

	sessionRaw := fmt.Sprintf("%d%d%d", uid, timestamp, randomNumber)

	hash := sha256.Sum256([]byte(sessionRaw))

	return hex.EncodeToString(hash[:])
}

func CreateSession(uid int64, r *http.Request) (string, error) {
	db := database.GetDB()
	if db == nil {
		return "", errors.New("database connection failed")
	}

	sessionID := GenerateSessionID(uid)             // 会话ID
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 有效期

	// 获取客户端 IP 和 UA
	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()

	timestamp := time.Now()

	_, err := db.Exec(
		"INSERT INTO session (session_id, uid, ip_address, user_agent,created_at, expires_at, last_activity) VALUES (?,?, ?, ?, ?, ?, ?)",
		sessionID, uid, ipAddress, userAgent, timestamp, expiresAt, timestamp,
	)

	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to create session: %v", err))
		return "", err
	}

	middleware.Logger.Log("INFO", fmt.Sprintf("Session created for user %d", uid))
	return sessionID, nil
}

func ValidateSession(sessionID string) (int64, error) {
	db := database.GetDB()
	if db == nil {
		return 0, errors.New("database connection failed")
	}

	// 查询会话
	var uid int64
	var expiresAt time.Time
	err := db.QueryRow(
		"SELECT uid, expires_at FROM session WHERE session_id = ?",
		sessionID,
	).Scan(&uid, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("session not found")
		}
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to validate session: %v", err))
		return 0, err
	}

	// 删除过期会话
	if time.Now().After(expiresAt) {
		_, _ = db.Exec("DELETE FROM session WHERE session_id = ?", sessionID)
		return 0, errors.New("session expired")
	}

	// 更新最后活动时间
	_, _ = db.Exec("UPDATE session SET last_activity = ? WHERE session_id = ?", time.Now(), sessionID)

	return uid, nil
}

func DeleteSession(sessionID string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection failed")
	}

	_, err := db.Exec("DELETE FROM session WHERE session_id = ?", sessionID)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to delete session: %v", err))
		return err
	}

	return nil
}

func DeleteExpiredSessions() int64 {
	db := database.GetDB()
	if db == nil {
		return 0
	}

	result, err := db.Exec("DELETE FROM session WHERE expires_at < ?", time.Now())
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to delete expired sessions: %v", err))
		return 0
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		middleware.Logger.Log("INFO", fmt.Sprintf("Deleted %d expired sessions", rowsAffected))
	}

	return rowsAffected
}

func GetUserSessions(uid int64) ([]_struct.Session, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database connection failed")
	}

	rows, err := db.Query(
		"SELECT session_id, uid, ip_address, user_agent, created_at, expires_at, last_activity FROM session WHERE uid = ? AND expires_at > ?",
		uid, time.Now(),
	)

	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get user sessions: %v", err))
		return nil, err
	}
	defer rows.Close()

	var sessions []_struct.Session
	for rows.Next() {
		var session _struct.Session
		err := rows.Scan(
			&session.SessionID,
			&session.UID,
			&session.IPAddress,
			&session.UserAgent,
			&session.CreatedAt,
			&session.ExpiresAt,
			&session.LastActivity,
		)

		if err != nil {
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to scan session: %v", err))
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func getIPAddress(r *http.Request) string {
	// 从X-Forwarded-For头获取
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// 从X-Real-IP头获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 从RemoteAddr获取
	return r.RemoteAddr
}
