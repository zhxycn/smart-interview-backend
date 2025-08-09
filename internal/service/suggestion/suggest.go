package suggestion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/database"
	"time"
)

func RunWorkflow(recent RecentFeedback) (string, error) {
	cfg := config.LoadConfig()

	url := fmt.Sprintf("%s/workflows/run", cfg.DifyEndpoint)

	input, _ := json.Marshal(recent)

	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"input": input,
		},
		"response_mode": "blocking",
		"user":          "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.SuggestionApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData struct {
		Data struct {
			Outputs struct {
				StructuredOutput json.RawMessage `json:"structured_output"`
			} `json:"outputs"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	output := string(respData.Data.Outputs.StructuredOutput)
	return output, nil
}

func Record(uid int64, data string) error {
	db := database.GetDB()

	timestamp := time.Now()
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM suggestion WHERE user = ?)", uid).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err := db.Exec("UPDATE suggestion SET created_at = ?, suggestion = ? WHERE user = ?", timestamp, data, uid)
		if err != nil {
			return err
		}
	} else {
		_, err := db.Exec("INSERT INTO suggestion (user, created_at, suggestion) VALUES (?, ?, ?)", uid, timestamp, data)
		if err != nil {
			return err
		}
	}

	return nil
}
