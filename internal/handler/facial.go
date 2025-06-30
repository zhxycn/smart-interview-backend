package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smart-interview/internal/config"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/facial"
	"smart-interview/internal/util"
	"strings"
)

type FacialRequest struct {
	ImageData   string `json:"image_data"`
	ImageFormat string `json:"image_format"`
}

func stripBase64Header(data string) string {
	if idx := strings.Index(data, "base64,"); idx != -1 {
		return data[idx+7:]
	}
	return data
}

func Facial(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()
	serverID := "s67c9c78c"

	var req FacialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, "invalid request body")
		middleware.Logger.Log("ERROR", fmt.Sprintf("invalid request body: %v", err))
		return
	}

	imageData := stripBase64Header(req.ImageData)

	result, err := facial.Detect(cfg.XunfeiAppId, cfg.XunfeiApiKey, cfg.XunfeiApiSecret, serverID, imageData, req.ImageFormat)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, err.Error())
		middleware.Logger.Log("ERROR", fmt.Sprintf("facial detection error: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, result)
}
