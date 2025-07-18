package interview

import (
	"encoding/json"
	"fmt"
	"smart-interview/internal/database"
	"time"
)

type Data struct {
	ID           string           `json:"id"`
	User         int              `json:"user"`
	Position     string           `json:"position"`
	Level        string           `json:"level"`
	CreatedAt    time.Time        `json:"created_at"`
	StartTime    time.Time        `json:"start_time"`
	EndTime      time.Time        `json:"end_time"`
	Conversation json.RawMessage  `json:"conversation"`
	Facial       json.RawMessage  `json:"facial"`
	Score        int              `json:"score"`
	Communicate  int              `json:"communicate"`
	Specialized  int              `json:"specialized"`
	Expression   int              `json:"expression"`
	Strain       int              `json:"strain"`
	Appearance   int              `json:"appearance"`
	Appraise     *json.RawMessage `json:"appraise"`
}

func Result(id string, uid int64) (Data, error) {
	var interview Data

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
