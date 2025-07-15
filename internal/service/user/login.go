package user

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"smart-interview/internal/database"
	"smart-interview/internal/database/struct"
	"smart-interview/internal/middleware"
)

func Login(email, password string) (*_struct.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("missing parameters")
	}

	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database connection failed")
	}

	var user _struct.User
	var hashedPassword string
	row := db.QueryRow("SELECT uid, email, password, name, job FROM user WHERE email = ?", email)
	err := row.Scan(&user.UID, &user.Email, &hashedPassword, &user.Name, &user.Job)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to query user: %v", err))
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		middleware.Logger.Log("WARN", fmt.Sprintf("Invalid login attempt for email: %s", email))
		return nil, errors.New("invalid credentials")
	}

	middleware.Logger.Log("INFO", fmt.Sprintf("User logged in. UID: %d, Email: %s", user.UID, user.Email))
	return &user, nil
}
