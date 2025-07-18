package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/interview"
	"smart-interview/internal/util"
)

func InterviewList(w http.ResponseWriter, r *http.Request) {
	uid, ok := util.GetUserID(r)
	if !ok {
		util.WriteResponse(w, http.StatusUnauthorized, nil)
		return
	}

	result, err := interview.List(uid)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get interview list: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
