package handler

import (
	"fmt"
	"net/http"
	"smart-interview/internal/service/resume"
	"smart-interview/internal/util"
)

func ResumeDownloadHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	id := r.URL.Query().Get("id")
	if id == "" {
		util.WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	fileName, fileData, err := resume.GetResumeFile(uid, id)
	if err != nil {
		if err.Error() == "not found" {
			util.WriteResponse(w, http.StatusNotFound, nil)
			return
		}
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fileData)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	util.WriteResponse(w, http.StatusOK, nil)
}
