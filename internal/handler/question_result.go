package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/question"
	"smart-interview/internal/util"
)

func QuestionResultHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		util.WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	uid, ok := util.GetUserID(r)
	if !ok {
		util.WriteResponse(w, http.StatusUnauthorized, nil)
		return
	}

	result, err := question.Result(id, uid)
	if err != nil {
		if err.Error() == "not found" {
			util.WriteResponse(w, http.StatusNotFound, nil)
			return
		}
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get question result: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
