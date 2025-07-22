package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/resume"
	"smart-interview/internal/util"
)

func ResumeListHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	result, err := resume.ListResumes(uid)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get resume list: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
