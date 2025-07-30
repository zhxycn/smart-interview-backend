package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/question"
	"smart-interview/internal/util"
)

func QuestionListHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	result, err := question.List(uid)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get question list: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
