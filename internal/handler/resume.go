package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/resume"
	"smart-interview/internal/util"
)

func ResumeHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, "failed to parse form")
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to parse form: %v", err))
		return
	}

	file, handler, err := r.FormFile("file") // 文件
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, "failed to get file")
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get file: %v", err))
		return
	}
	defer file.Close()

	middleware.Logger.Log("INFO", fmt.Sprintf("Received file: %v", handler.Filename))

	params := resume.Params{
		TargetPosition: r.FormValue("targetPosition"),
		Experience:     r.FormValue("experience"),
		Industry:       r.FormValue("industry"),
		FocusAreas:     r.FormValue("focusAreas"),
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, "failed to read file")
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to read file: %v", err))
		return
	}

	result, err := resume.Analysis(handler.Filename, buf.Bytes(), params)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, "failed to analyze")
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to analyze: %v", err))
		return
	}

	err = resume.Record(uid, handler.Filename, buf.Bytes(), params.TargetPosition, params.Experience, params.Industry, params.FocusAreas, result)
	if err != nil {
		middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to record resume: %v", err))
	}

	util.WriteResponse(w, http.StatusOK, result)
}
