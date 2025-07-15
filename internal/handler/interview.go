package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/interview"
	"smart-interview/internal/util"
)

type AudioRequest struct {
	Data  string `json:"data"`
	IsEnd bool   `json:"is_end"`
}

type Response struct {
	ASR   string `json:"asr"`
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
	Error string `json:"error,omitempty"`
}

func Interview(w http.ResponseWriter, r *http.Request) {
	var req AudioRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("%v", err))
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("%v", err))
		return
	}

	cfg := config.LoadConfig()
	response, err := interview.Interview(cfg.TencentAppId, cfg.TencentSecretId, cfg.TencentSecretKey, req.Data)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("%v", err))
		return
	}

	middleware.Logger.Log("DEBUG", fmt.Sprintf("[Interview] ASR: %s, Text: %s", response.ASR, response.Text))

	util.WriteResponse(w, http.StatusOK, Response{
		ASR:   response.ASR,
		Text:  response.Text,
		Audio: response.Audio,
	})
}
