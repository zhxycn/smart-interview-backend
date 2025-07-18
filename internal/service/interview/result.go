package interview

import (
	"fmt"
	"smart-interview/internal/database"
	"smart-interview/internal/database/struct"
)

func Result(id string, uid int64) (_struct.Interview, error) {
	var interview _struct.Interview

	db := database.GetDB()
	if db == nil {
		return interview, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query(`SELECT * FROM interview WHERE id = ? AND user = ?`, id, uid)
	if err != nil {
		return interview, err
	}
	defer rows.Close()

	if !rows.Next() {
		return interview, fmt.Errorf("not found")
	}
	err = rows.Scan(
		&interview.ID,
		&interview.User,
		&interview.Position,
		&interview.Level,
		&interview.CreatedAt,
		&interview.StartTime,
		&interview.EndTime,
		&interview.Conversation,
		&interview.Facial,
		&interview.Score,
		&interview.Communicate,
		&interview.Specialized,
		&interview.Expression,
		&interview.Strain,
		&interview.Appearance,
		&interview.Appraise,
	)
	if err != nil {
		return interview, err
	}

	return interview, nil
}
