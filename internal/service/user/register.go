package user

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"smart-interview/internal/database"
	"smart-interview/internal/middleware"
	"time"
)

func generateRandomUID() int64 {
	source := rand.NewSource(time.Now().UnixNano()) // 当前时间作为随机数种子
	random := rand.New(source)
	return random.Int63n(9000000000000000000) + 1000000000000000000
}

func Register(email, password, name string) (int64, error) {
	if email == "" || password == "" || name == "" {
		return 0, errors.New("missing parameters")
	}

	db := database.GetDB()
	if db == nil {
		return 0, errors.New("database connection failed")
	}

	// 检查邮箱是否已存在
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email)
	if err := row.Scan(&count); err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to check email existence: %v", err))
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("email already exists")
	}

	// bcrypt 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to hash password: %v", err))
		return 0, err
	}

	// 生成 UID
	uid := generateRandomUID()
	var uidCount int
	for {
		row := db.QueryRow("SELECT COUNT(*) FROM user WHERE uid = ?", uid) // 检查是否存在
		if err := row.Scan(&uidCount); err != nil {
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to check UID existence: %v", err))
			return 0, err
		}

		if uidCount == 0 {
			break
		}

		uid = generateRandomUID()
	}

	_, err = db.Exec("INSERT INTO user (uid, email, password, name) VALUES (?, ?, ?, ?)", uid, email, string(hashedPassword), name)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to insert new user: %v", err))
		return 0, err
	}

	middleware.Logger.Log("INFO", fmt.Sprintf("New user registered. UID: %d, Email: %s", uid, email))
	return uid, nil // 返回 UID
}
