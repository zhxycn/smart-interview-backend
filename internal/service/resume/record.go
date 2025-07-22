package resume

import (
	"fmt"
	"github.com/google/uuid"
	"smart-interview/internal/database"
	"time"
)

type Resume struct {
	ID             string           `json:"id"`
	User           int64            `json:"user"`
	FileName       string           `json:"file_name"`
	CreatedAt      time.Time        `json:"created_at"`
	TargetPosition string           `json:"target_position"`
	Experience     string           `json:"experience"`
	Industry       string           `json:"industry"`
	FocusAreas     string           `json:"focus_areas"`
	Score          int              `json:"score"`
	Feedback       StructuredOutput `json:"feedback"`
}

func Record(user int64, fileName string, fileData []byte, targetPosition, experience, industry, focusAreas string, feedback StructuredOutput) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection failed")
	}

	resumeId := uuid.New().String()

	timestamp := time.Now()

	_, err := db.Exec(
		`INSERT INTO resume (id, user, file_name, file_data, created_at, target_position, experience, industry, focus_areas, score, feedback) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		resumeId, user, fileName, fileData, timestamp, targetPosition, experience, industry, focusAreas, feedback.OverallScore, feedback,
	)
	if err != nil {
		return fmt.Errorf("failed to record resume: %w", err)
	}

	return nil
}

func GetResume(uid int64, id string) (Resume, error) {
	var resume Resume

	db := database.GetDB()
	if db == nil {
		return resume, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query(
		`SELECT id, user, file_name, created_at, target_position, experience, industry, focus_areas, score, feedback FROM resume WHERE id = ? AND user = ?`,
		id, uid,
	)
	if err != nil {
		return resume, err
	}
	defer rows.Close()

	if !rows.Next() {
		return resume, fmt.Errorf("not found")
	}
	err = rows.Scan(
		&resume.ID,
		&resume.User,
		&resume.FileName,
		&resume.CreatedAt,
		&resume.TargetPosition,
		&resume.Experience,
		&resume.Industry,
		&resume.FocusAreas,
		&resume.Score,
		&resume.Feedback,
	)
	if err != nil {
		return resume, err
	}

	return resume, nil
}

func GetResumeFile(uid int64, id string) (string, []byte, error) {
	db := database.GetDB()
	if db == nil {
		return "", nil, fmt.Errorf("database connection failed")
	}

	var fileName string
	var fileData []byte
	err := db.QueryRow(
		`SELECT file_name, file_data FROM resume WHERE id = ? AND user = ?`,
		id, uid,
	).Scan(&fileName, &fileData)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get resume file: %v", err)
	}

	return fileName, fileData, nil
}
