package question

import (
	"encoding/json"
	"fmt"
	"smart-interview/internal/database"
)

type QL struct {
	ID        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	Position  string   `json:"position"`
	Knowledge []string `json:"knowledge"`
	Count     int      `json:"count"`
}

func List(uid int64) ([]*QL, error) {
	var ql []*QL

	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query("SELECT id, created_at, position, knowledge, count FROM question WHERE user = ? ORDER BY created_at", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i QL
		var knowledge []byte

		err := rows.Scan(&i.ID, &i.CreatedAt, &i.Position, &knowledge, &i.Count)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(knowledge, &i.Knowledge); err != nil {
			return nil, err
		}

		ql = append(ql, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ql, nil
}
