package interview

import (
	"fmt"
	"smart-interview/internal/database"
	"time"
)

type I struct {
	ID          string     `json:"id"`
	Position    string     `json:"position"`
	Level       string     `json:"level"`
	CreatedAt   time.Time  `json:"created_at"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Score       int        `json:"score"`
	Communicate int        `json:"communicate"`
	Specialized int        `json:"specialized"`
	Expression  int        `json:"expression"`
	Strain      int        `json:"strain"`
	Appearance  int        `json:"appearance"`
}

func List(uid int64) ([]*I, error) {
	var interviews []*I

	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query(
		`SELECT id, position, level, created_at, start_time, end_time, score, communicate, specialized, expression, strain, appearance
				FROM interview WHERE user = ? ORDER BY created_at`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i I
		err := rows.Scan(
			&i.ID, &i.Position, &i.Level, &i.CreatedAt, &i.StartTime, &i.EndTime, &i.Score,
			&i.Communicate, &i.Specialized, &i.Expression, &i.Strain, &i.Appearance,
		)
		if err != nil {
			return nil, err
		}
		if i.EndTime == nil {
			continue // 处理EndTime为空，即未结束的面试
		}
		interviews = append(interviews, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return interviews, nil
}
