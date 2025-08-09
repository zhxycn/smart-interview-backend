package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/suggestion"
	"smart-interview/internal/util"
)

func SuggestionHandler(w http.ResponseWriter, r *http.Request) {
	uid, _ := util.GetUserID(r)

	if r.Method == http.MethodGet {
		result, err := suggestion.GetSuggestion(uid)
		if err != nil {
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get suggestion: %v", err))
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			return
		}

		if result == nil {
			util.WriteResponse(w, http.StatusNotFound, nil)
			return
		}

		var data interface{}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to marshal result: %v", err))
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			return
		}
		err = json.Unmarshal(jsonBytes, &data)

		util.WriteResponse(w, http.StatusOK, result)
		return
	}

	if r.Method == http.MethodPost {
		go func() {
			recent, err := suggestion.GetRecentFeedback(uid)
			if err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get recent feedback: %v", err))
				util.WriteResponse(w, http.StatusInternalServerError, nil)
				return
			}

			result, err := suggestion.RunWorkflow(recent)
			if err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to run workflow: %v", err))
				util.WriteResponse(w, http.StatusInternalServerError, nil)
				return
			}

			err = suggestion.Record(uid, result)
			if err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to record suggestion: %v", err))
			}
		}()

		util.WriteResponse(w, http.StatusOK, nil)
		return
	}

	util.WriteResponse(w, http.StatusMethodNotAllowed, nil)
}
