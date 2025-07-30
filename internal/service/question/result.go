package question

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"smart-interview/internal/database"
)

type Data struct {
	ID        string       `json:"id"`
	User      int64        `json:"user"`
	CreatedAt string       `json:"created_at"`
	Position  string       `json:"position"`
	Knowledge []string     `json:"knowledge"`
	Count     int          `json:"count"`
	Feedback  FeedbackResp `json:"feedback"`
	Questions []Question   `json:"questions"`
}

func Result(id string, uid int64) (Data, error) {
	var d Data

	db := database.GetDB()
	if db == nil {
		return d, fmt.Errorf("database connection failed")
	}

	rows, err := db.Query(`SELECT * FROM question WHERE id = ? AND user = ?`, id, uid)
	if err != nil {
		return d, err
	}
	defer rows.Close()

	if !rows.Next() {
		return d, fmt.Errorf("not found")
	}

	var knowledge []byte
	var feedback []byte

	err = rows.Scan(&d.ID, &d.User, &d.CreatedAt, &d.Position, &knowledge, &d.Count, &feedback)
	if err != nil {
		return d, err
	}

	if err := json.Unmarshal(knowledge, &d.Knowledge); err != nil {
		return d, err
	}

	if err := json.Unmarshal(feedback, &d.Feedback); err != nil {
		return d, err
	}

	rows, err = db.Query("SELECT qid, question, answer, difficulty, user_answer, feedback FROM questions WHERE id = ?", id)
	if err != nil {
		return d, err
	}
	defer rows.Close()

	for rows.Next() {
		var question Question
		var qfeedback []byte
		var userAnswer sql.NullString

		err = rows.Scan(&question.QID, &question.Question, &question.Answer, &question.Difficulty, &userAnswer, &qfeedback)
		if err != nil {
			return d, err
		}

		if userAnswer.Valid {
			question.UserAnswer = userAnswer.String
		} else {
			question.UserAnswer = ""
		}

		if len(qfeedback) > 0 {
			if err = json.Unmarshal(qfeedback, &question.Feedback); err != nil {
				return d, err
			}
		}

		d.Questions = append(d.Questions, question)
	}

	return d, nil
}
