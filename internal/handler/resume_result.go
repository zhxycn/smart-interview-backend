package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/resume"
	"smart-interview/internal/util"
)

func ResumeResultHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	id := r.URL.Query().Get("id")
	if id == "" {
		util.WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	result, err := resume.GetResume(uid, id)
	if err != nil {
		if err.Error() == "not found" {
			util.WriteResponse(w, http.StatusNotFound, nil)
			return
		}
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get interview result: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
