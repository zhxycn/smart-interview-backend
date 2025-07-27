package question

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"strings"
	"time"
)

type Question struct {
	QID        string        `json:"qid"`
	Question   string        `json:"question"`
	Answer     string        `json:"answer"`
	Difficulty string        `json:"difficulty"`
	UserAnswer string        `json:"user_answer,omitempty"`
	Feedback   QFeedbackResp `json:"feedback,omitempty"`
}

type ResponseQuestion struct {
	QID        string `json:"qid"`
	Question   string `json:"question"`
	Difficulty string `json:"difficulty"`
}

type Response struct {
	ID        string             `json:"id"`
	Position  string             `json:"position"`
	Knowledge []string           `json:"knowledge"`
	Count     int                `json:"count"`
	Questions []ResponseQuestion `json:"questions"`
}

type apiResponse struct {
	Data struct {
		Outputs struct {
			Questions []Question `json:"questions"`
		} `json:"outputs"`
	} `json:"data"`
}

func GenerateQuestion(position string, knowledge []string, count int) (string, []Question, error) {
	cfg := config.LoadConfig()
	endpoint := fmt.Sprintf("%s/workflows/run", cfg.DifyEndpoint)

	payload := map[string]interface{}{
		"conversation_id": "",
		"inputs": map[string]interface{}{
			"position":  position,
			"knowledge": strings.Join(knowledge, ", "),
			"count":     count,
		},
		"response_mode":     "blocking",
		"parent_message_id": nil,
		"user":              "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.QuestionApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var respData apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", nil, err
	}

	questions := respData.Data.Outputs.Questions

	for i := range questions {
		qid := uuid.New().String()
		questions[i].QID = qid
	}

	id := uuid.New().String()

	return id, questions, nil
}

func SaveQuestion(uid int64, id, position string, knowledge []string, count int, questions []Question) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection failed")
	}

	timestamp := time.Now()

	knowledgeJSON, err := json.Marshal(knowledge)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO question (id, user, created_at, position, knowledge, count) VALUES (?, ?, ?, ?, ?, ?)",
		id, uid, timestamp, position, string(knowledgeJSON), count)
	if err != nil {
		return err
	}

	for _, question := range questions {
		_, err = db.Exec("INSERT INTO questions (qid, id, question, answer, difficulty) VALUES (?, ?, ?, ?, ?)",
			question.QID, id, question.Question, question.Answer, question.Difficulty)
		if err != nil {
			return err
		}
	}

	return nil
}

func FormatQuestion(id, position string, knowledge []string, q []Question) (Response, error) {
	if len(q) == 0 {
		return Response{}, fmt.Errorf("no questions available")
	}

	resp := Response{
		ID:        id,
		Position:  position,
		Knowledge: knowledge,
		Count:     len(q),
		Questions: make([]ResponseQuestion, 0, len(q)),
	}

	for _, item := range q {
		resp.Questions = append(resp.Questions, ResponseQuestion{
			QID:        item.QID,
			Question:   item.Question,
			Difficulty: item.Difficulty,
		})
	}

	return resp, nil
}
