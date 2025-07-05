package interview

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"smart-interview/internal/config"
	"smart-interview/internal/middleware"
	"sort"
	"strconv"
	"time"
)

type FlashResult struct {
	Text string `json:"text"`
}

type FlashResponse struct {
	Code        int           `json:"code"`
	Message     string        `json:"message"`
	FlashResult []FlashResult `json:"flash_result"`
}

type WorkflowResponse struct {
	ASR   string `json:"asr"`
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
}

type DifyWorkflowResponse struct {
	Answer         string `json:"answer"`
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
}

func Interview(appId, secretId, secretKey, base64Audio string) (WorkflowResponse, error) {
	asrResult, err := RecognizeAudio(appId, secretId, secretKey, base64Audio)
	if err != nil {
		return WorkflowResponse{}, err
	}

	difyResponse, err := CallInterviewWorkflow(asrResult)
	if err != nil {
		return WorkflowResponse{}, fmt.Errorf("failed to call Dify workflow: %v", err)
	}

	audio, err := CallTTSWorkflow(difyResponse.Answer)
	if err != nil {
		return WorkflowResponse{}, fmt.Errorf("failed to call TTS workflow: %v", err)
	}

	return WorkflowResponse{
		ASR:   asrResult,
		Text:  difyResponse.Answer,
		Audio: audio,
	}, nil
}

func RecognizeAudio(appId, secretId, secretKey, base64Audio string) (string, error) {
	// 解码base64音频
	audioData, err := base64.StdEncoding.DecodeString(base64Audio)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 audio: %v", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	engineType := "16k_zh"
	voiceFormat := "wav"

	paramsMap := map[string]string{
		"engine_type":         engineType,
		"extra_punc":          "0",
		"filter_punc":         "0",
		"first_channel_only":  "0",
		"reinforce_hotword":   "0",
		"secretid":            secretId,
		"speaker_diarization": "0",
		"timestamp":           timestamp,
		"voice_format":        voiceFormat,
		"word_info":           "0",
	}

	var keys []string
	for k := range paramsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryString string
	for i, k := range keys {
		if i > 0 {
			queryString += "&"
		}
		queryString += k + "=" + url.QueryEscape(paramsMap[k])
	}

	urlPath := fmt.Sprintf("/asr/flash/v1/%s", appId)
	origin := "POST" + "asr.cloud.tencent.com" + urlPath + "?" + queryString
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(origin))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	requestUrl := "https://asr.cloud.tencent.com" + urlPath + "?" + queryString
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Host", "asr.cloud.tencent.com")
	req.Header.Set("Authorization", signature)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(audioData)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var flashResp FlashResponse
	if err := json.Unmarshal(respBody, &flashResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(respBody))
	}
	if flashResp.Code != 0 {
		return "", fmt.Errorf("%s", flashResp.Message)
	}

	var recognizedText string
	if len(flashResp.FlashResult) > 0 {
		recognizedText = flashResp.FlashResult[0].Text
	}

	middleware.Logger.Log("DEBUG", fmt.Sprintf("[ASR] Recognized text: %s", recognizedText))

	return recognizedText, nil
}

func CallInterviewWorkflow(text string) (DifyWorkflowResponse, error) {
	endpoint := fmt.Sprintf("%s/chat-messages", config.LoadConfig().DifyEndpoint)

	payload := map[string]interface{}{
		"query":           text,
		"conversation_id": "",
		"inputs": map[string]interface{}{
			"data": nil,
		},
		"response_mode":     "blocking",
		"parent_message_id": nil,
		"user":              "admin",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return DifyWorkflowResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LoadConfig().InterviewApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DifyWorkflowResponse{}, err
	}
	defer resp.Body.Close()

	var respData DifyWorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return DifyWorkflowResponse{}, err
	}

	return respData, nil
}

func CallTTSWorkflow(text string) (string, error) {
	endpoint := "https://api.siliconflow.cn/v1/audio/speech"

	payload := map[string]interface{}{
		"model":           config.LoadConfig().SiliconflowModel,
		"input":           text,
		"voice":           config.LoadConfig().SiliconflowVoice,
		"response_format": "wav",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+config.LoadConfig().SiliconflowToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取音频二进制数据
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 转为base64字符串
	base64Audio := base64.StdEncoding.EncodeToString(audioData)
	return base64Audio, nil
}
