package resume

import (
	"fmt"
	"smart-interview/internal/database"
	"time"
)

type R struct {
	ID             string    `json:"id"`
	FileName       string    `json:"file_name"`
	CreatedAt      time.Time `json:"created_at"`
	TargetPosition string    `json:"target_position"`
	Experience     string    `json:"experience"`
	Industry       string    `json:"industry"`
	FocusAreas     string    `json:"focus_areas"`
	Score          int       `json:"score"`
}

func ListResumes(uid int64) ([]*R, error) {
	var resumes []*R

	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query(
		`SELECT id, file_name, created_at, target_position, experience, industry, focus_areas, score
				FROM resume WHERE user = ? ORDER BY created_at`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r R
		err := rows.Scan(
			&r.ID, &r.FileName, &r.CreatedAt, &r.TargetPosition, &r.Experience,
			&r.Industry, &r.FocusAreas, &r.Score,
		)
		if err != nil {
			return nil, err
		}
		resumes = append(resumes, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resumes, nil
}
