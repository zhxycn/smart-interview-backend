package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/interview"
	"smart-interview/internal/util"
)

type InterviewRequest struct {
	Position string `json:"position"` // 岗位
	Level    string `json:"level"`    // 难度
}

type InterviewResponse struct {
	ID string `json:"id"`
}

func InterviewRegister(w http.ResponseWriter, r *http.Request) {
	var req InterviewRequest
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

	uid, ok := util.GetUserID(r)
	if !ok {
		util.WriteResponse(w, http.StatusUnauthorized, nil)
		return
	}

	interviewId, err := interview.Register(uid, req.Position, req.Level)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, InterviewResponse{ID: interviewId})
}
