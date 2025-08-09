package suggestion

import (
	"fmt"
	"smart-interview/internal/database"
)

type RecentFeedback struct {
	Interview []InterviewRF `json:"interview"`
	Question  []QuestionRF  `json:"question"`
	Resume    []ResumeRF    `json:"resume"`
}

type InterviewRF struct {
	Problem    []string `json:"problem"`
	Suggestion []string `json:"suggestion"`
}

type QuestionRF struct {
	Feedback string `json:"feedback"`
}

type ResumeRF struct {
	Feedback string `json:"feedback"`
}

func GetRecentFeedback(uid int64) (RecentFeedback, error) {
	db := database.GetDB()
	if db == nil {
		return RecentFeedback{}, fmt.Errorf("database connection failed")
	}

	var rf RecentFeedback

	// 模拟面试
	rows, err := db.Query(
		"SELECT problem, suggestion FROM interview WHERE problem AND suggestion IS NOT NULL AND user=? ORDER BY end_time DESC LIMIT 5",
		uid,
	)
	if err != nil {
		return rf, err
	}

	var interview []InterviewRF

	for rows.Next() {
		var problem, suggestion string
		if err := rows.Scan(&problem, &suggestion); err != nil {
			return rf, err
		}
		interview = append(interview, InterviewRF{
			Problem:    []string{problem},
			Suggestion: []string{suggestion},
		})
	}
	if err := rows.Close(); err != nil {
		return rf, err
	}

	rf.Interview = interview

	// 面试题
	rows, err = db.Query(
		"SELECT feedback FROM question WHERE feedback IS NOT NULL AND user=? ORDER BY created_at DESC LIMIT 5",
		uid,
	)
	if err != nil {
		return rf, err
	}

	var question []QuestionRF

	for rows.Next() {
		var feedback string
		if err := rows.Scan(&feedback); err != nil {
			return rf, err
		}
		question = append(question, QuestionRF{
			Feedback: feedback,
		})
	}
	if err := rows.Close(); err != nil {
		return rf, err
	}

	rf.Question = question

	// 简历评价
	rows, err = db.Query(
		"SELECT feedback FROM resume WHERE feedback IS NOT NULL AND user=? ORDER BY created_at DESC LIMIT 5",
		uid,
	)
	if err != nil {
		return rf, err
	}

	var resume []ResumeRF
	for rows.Next() {
		var feedback string
		if err := rows.Scan(&feedback); err != nil {
			return rf, err
		}
		resume = append(resume, ResumeRF{
			Feedback: feedback,
		})
	}
	if err := rows.Close(); err != nil {
		return rf, err
	}

	rf.Resume = resume

	return rf, nil
}
