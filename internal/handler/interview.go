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
	Id   string `json:"id"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
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

	uid, ok := util.GetUserID(r)
	if !ok {
		util.WriteResponse(w, http.StatusUnauthorized, nil)
		middleware.Logger.Log("ERROR", "Unauthorized access")
		return
	}

	if req.Msg == "start" {
		ok, err := interview.StartInterview(uid, req.Id)
		if !ok {
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to start interview: %v", err))
			return
		}
		middleware.Logger.Log("INFO", fmt.Sprintf("Interview started. User %d, ID %s", uid, req.Id))
	}

	if req.Msg == "stop" {
		ok, err := interview.StopInterview(uid, req.Id)
		if !ok {
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to stop interview: %v", err))
			return
		}
		util.WriteResponse(w, http.StatusOK, nil)
		middleware.Logger.Log("INFO", fmt.Sprintf("Interview stopped. User %d, ID %s", uid, req.Id))
		return
	}

	response, err := interview.Interview(cfg.TencentAppId, cfg.TencentSecretId, cfg.TencentSecretKey, req.Data, req.Msg, req.Id)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("%v", err))
		return
	}

	// 记录用户内容
	ok, err = interview.Record(req.Id, response.ASR, "user")
	if !ok {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to record user data: %v", err))
		return
	}

	// 记录AI内容
	ok, err = interview.Record(req.Id, response.Text, "assistant")
	if !ok {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to record assistant data: %v", err))
		return
	}

	middleware.Logger.Log("DEBUG", fmt.Sprintf("[Interview] ASR: %s, Text: %s", response.ASR, response.Text))

	util.WriteResponse(w, http.StatusOK, Response{
		ASR:   response.ASR,
		Text:  response.Text,
		Audio: response.Audio,
	})
}
