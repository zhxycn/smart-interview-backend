package interview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"smart-interview/internal/middleware"
)

type StructuredOutput struct {
	Score       int `json:"score"`
	Communicate int `json:"communicate"`
	Specialized int `json:"specialized"`
	Expression  int `json:"expression"`
	Strain      int `json:"strain"`
	Appearance  int `json:"appearance"`
	Appraise    struct {
		Specialize struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"specialize"`
		Skill struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"skill"`
		Express struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"express"`
		Logic struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"logic"`
		Innovation struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"innovation"`
		Stress struct {
			Score   string `json:"score"`
			Comment string `json:"comment"`
		} `json:"stress"`
		Emotion struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"emotion"`
		Learning struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
		} `json:"learning"`
	} `json:"appraise"`
	Process []struct {
		Question string `json:"question"`
		Score    int    `json:"score"`
		Time     int    `json:"time"`
		Progress int    `json:"progress"`
		Ability  string `json:"ability"`
	} `json:"process"`
	Problem    []string `json:"problem"`
	Suggestion []string `json:"suggestion"`
}

func GenerateInterviewResult(uid int64, id string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection failed")
	}

	data, err := Result(id, uid)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get interview data: %v", err))
		return err
	}

	conf := config.LoadConfig()

	url := fmt.Sprintf("%s/workflows/run", conf.DifyEndpoint)

	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"conversation": func() string {
				if data.Conversation != nil {
					return string(*data.Conversation)
				}
				return ""
			}(),
			"facial": func() string {
				if data.Facial != nil {
					return string(*data.Facial)
				}
				return ""
			}(),
			"position": data.Position,
			"level":    data.Level,
		},
		"response_mode": "blocking",
		"user":          "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to create request: %v", err))
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.ResultApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to send request: %v", err))
		return err
	}
	defer resp.Body.Close()

	var respData struct {
		Data struct {
			Outputs struct {
				StructuredOutput StructuredOutput `json:"structured_output"`
			} `json:"outputs"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to decode response: %v", err))
		return err
	}

	appraiseJSON, err := json.Marshal(respData.Data.Outputs.StructuredOutput.Appraise)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to marshal: %v", err))
		return err
	}

	processJSON, err := json.Marshal(respData.Data.Outputs.StructuredOutput.Process)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to marshal: %v", err))
	}

	problemJSON, err := json.Marshal(respData.Data.Outputs.StructuredOutput.Problem)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to marshal: %v", err))
		return err
	}

	suggestionJSON, err := json.Marshal(respData.Data.Outputs.StructuredOutput.Suggestion)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to marshal: %v", err))
		return err
	}

	_, err = db.Exec(
		"UPDATE interview SET score = ?, communicate = ?, specialized = ?, expression = ?, strain = ?, appearance = ?, appraise = ?, process = ?, problem = ?, suggestion = ? WHERE id = ? AND user = ?",
		respData.Data.Outputs.StructuredOutput.Score,
		respData.Data.Outputs.StructuredOutput.Communicate,
		respData.Data.Outputs.StructuredOutput.Specialized,
		respData.Data.Outputs.StructuredOutput.Expression,
		respData.Data.Outputs.StructuredOutput.Strain,
		respData.Data.Outputs.StructuredOutput.Appearance,
		appraiseJSON,
		processJSON,
		problemJSON,
		suggestionJSON,
		id,
		uid,
	)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to update interview result: %v", err))
		return err
	}

	middleware.Logger.Log("INFO", fmt.Sprintf("Interview result generated. ID %s, User %d", id, uid))

	return nil
}
