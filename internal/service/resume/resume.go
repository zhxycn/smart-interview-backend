package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"smart-interview/internal/config"
)

type Params struct {
	TargetPosition string
	Experience     string
	Industry       string
	FocusAreas     string
}

type StructuredOutput struct {
	OverallScore   int `json:"overallScore"`
	DetailedScores []struct {
		Category string `json:"category"`
		Score    int    `json:"score"`
		Comment  string `json:"comment"`
	} `json:"detailedScores"`
	Suggestions []struct {
		Priority int    `json:"priority"`
		Title    string `json:"title"`
		Content  string `json:"content"`
	} `json:"suggestions"`
	KeywordAnalysis struct {
		Matched   []string `json:"matched"`
		Suggested []string `json:"suggested"`
	} `json:"keywordAnalysis"`
	Position string `json:"position"`
}

func uploadFile(fileName string, fileData []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	if _, err = fw.Write(fileData); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/files/upload", config.LoadConfig().DifyEndpoint), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LoadConfig().ResumeApiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		FileID string `json:"id"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.FileID, nil
}

func RunWorkflow(fileID string, params Params) (StructuredOutput, error) {
	conf := config.LoadConfig()

	url := fmt.Sprintf("%s/workflows/run", conf.DifyEndpoint)

	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"transfer_method": "local_file",
					"upload_file_id":  fileID,
					"url":             "",
					"type":            "document",
				},
			},
			"targetPosition": params.TargetPosition,
			"experience":     params.Experience,
			"industry":       params.Industry,
			"focusAreas":     params.FocusAreas,
		},
		"response_mode": "blocking",
		"user":          "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return StructuredOutput{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.ResumeApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return StructuredOutput{}, err
	}
	defer resp.Body.Close()

	var respData struct {
		Data struct {
			Outputs struct {
				StructuredOutput StructuredOutput `json:"structured_output"`
			} `json:"outputs"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return StructuredOutput{}, err
	}

	return respData.Data.Outputs.StructuredOutput, nil
}

func Analysis(fileName string, fileData []byte, params Params) (StructuredOutput, error) {
	fileID, err := uploadFile(fileName, fileData)
	if err != nil {
		return StructuredOutput{}, err
	}

	result, err := RunWorkflow(fileID, params)
	if err != nil {
		return StructuredOutput{}, err
	}

	return result, nil
}
