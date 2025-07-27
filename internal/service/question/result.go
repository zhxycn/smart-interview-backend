package question

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"strings"
)

type ResultRequest struct {
	QID        string `json:"qid"`
	UserAnswer string `json:"user_answer"`
}

type QData struct {
	ID        string     `json:"id"`
	Position  string     `json:"position"`
	Knowledge []string   `json:"knowledge"`
	Count     int        `json:"count"`
	Questions []Question `json:"questions"`
}

type QFeedbackResp struct {
	Appraise    string `json:"appraise"`
	Correctness string `json:"correctness"`
}

type FeedbackResp struct {
	Weak       []string `json:"weak"`
	Suggestion string   `json:"suggestion"`
}

type qApiResponse struct {
	Data struct {
		Outputs struct {
			Appraises QFeedbackResp `json:"appraises"`
		} `json:"outputs"`
	} `json:"data"`
}

type ApiResponse struct {
	Data struct {
		Outputs struct {
			Appraise FeedbackResp `json:"appraise"`
		} `json:"outputs"`
	} `json:"data"`
}

func Feedback(uid int64, id string, data []ResultRequest) (QData, error) {
	db := database.GetDB()
	if db == nil {
		return QData{}, fmt.Errorf("database connection failed")
	}

	var q QData
	q.ID = id

	row := db.QueryRow("SELECT position, knowledge, count FROM question WHERE id = ? AND user = ?", id, uid)

	var knowledgeStr string

	if err := row.Scan(&q.Position, &knowledgeStr, &q.Count); err != nil {
		return QData{}, err
	}

	if err := json.Unmarshal([]byte(knowledgeStr), &q.Knowledge); err != nil {
		return QData{}, err
	}

	for _, item := range data {
		var question Question

		row := db.QueryRow("SELECT qid, question, answer, difficulty FROM questions WHERE id = ? AND qid = ?", id, item.QID)

		if err := row.Scan(&question.QID, &question.Question, &question.Answer, &question.Difficulty); err != nil {
			return QData{}, err
		}

		question.UserAnswer = item.UserAnswer

		qf, err := QFeedback(question.Question, q.Position, question.Difficulty, question.Answer, item.UserAnswer, q.Knowledge)

		question.Feedback = qf

		qfJSON, err := json.Marshal(qf)
		if err != nil {
			return QData{}, err
		}

		_, err = db.Exec("UPDATE questions SET user_answer = ?, feedback = ? WHERE id = ? AND qid = ?", item.UserAnswer, qfJSON, id, item.QID)
		if err != nil {
			return QData{}, err
		}
	}

	rows, err := db.Query("SELECT qid, question, answer, difficulty, user_answer, feedback FROM questions WHERE id = ?", id)
	if err != nil {
		return QData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var question Question
		var feedbackJSON sql.NullString
		var userAnswer sql.NullString

		if err := rows.Scan(&question.QID, &question.Question, &question.Answer, &question.Difficulty, &userAnswer, &feedbackJSON); err != nil {
			return QData{}, err
		}

		if userAnswer.Valid {
			question.UserAnswer = userAnswer.String
		} else {
			question.UserAnswer = ""
		}

		if feedbackJSON.Valid && feedbackJSON.String != "" {
			var feedback QFeedbackResp
			if err := json.Unmarshal([]byte(feedbackJSON.String), &feedback); err != nil {
				return QData{}, err
			}
			question.Feedback = feedback
		}

		q.Questions = append(q.Questions, question)
	}

	return q, nil
}

func QFeedback(question, position, difficulty, answer, userAnswer string, knowledge []string) (QFeedbackResp, error) {
	cfg := config.LoadConfig()
	endpoint := fmt.Sprintf("%s/workflows/run", cfg.DifyEndpoint)

	payload := map[string]interface{}{
		"conversation_id": "",
		"inputs": map[string]interface{}{
			"position":    position,
			"knowledge":   strings.Join(knowledge, ", "),
			"user_answer": userAnswer,
			"question":    question,
			"answer":      answer,
			"difficulty":  difficulty,
		},
		"response_mode":     "blocking",
		"parent_message_id": nil,
		"user":              "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return QFeedbackResp{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.QuestionJudgmentApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return QFeedbackResp{}, err
	}
	defer resp.Body.Close()

	var respData qApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return QFeedbackResp{}, err
	}

	return respData.Data.Outputs.Appraises, nil
}

func AllFeedback(uid int64, data QData) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection failed")
	}

	cfg := config.LoadConfig()
	endpoint := fmt.Sprintf("%s/workflows/run", cfg.DifyEndpoint)

	questionsJSON, err := json.Marshal(data.Questions)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"conversation_id": "",
		"inputs": map[string]interface{}{
			"position":  data.Position,
			"knowledge": strings.Join(data.Knowledge, ", "),
			"data":      string(questionsJSON),
		},
		"response_mode":     "blocking",
		"parent_message_id": nil,
		"user":              "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.QuestionResultApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respData ApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return err
	}

	respJSON, err := json.Marshal(respData.Data.Outputs.Appraise)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE question SET feedback = ? WHERE id = ? AND user = ?", respJSON, data.ID, uid)
	if err != nil {
		return err
	}

	return nil
}
