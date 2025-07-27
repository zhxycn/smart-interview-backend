package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/question"
	"smart-interview/internal/util"
)

type QuestionRequest struct {
	IsQuestion bool                     `json:"is_question"`
	ID         string                   `json:"id"`
	Position   string                   `json:"position"`
	Knowledge  []string                 `json:"knowledge"`
	Count      int                      `json:"count"`
	Data       []question.ResultRequest `json:"data"`
}

func QuestionHandler(w http.ResponseWriter, r *http.Request) {
	var req QuestionRequest
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
		middleware.Logger.Log("ERROR", "Unauthorized access")
		return
	}

	if req.IsQuestion {
		id, result, err := question.GenerateQuestion(req.Position, req.Knowledge, req.Count)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to generate questions: %v", err))
			return
		}

		go func() {
			if err := question.SaveQuestion(uid, id, req.Position, req.Knowledge, req.Count, result); err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to save questions: %v", err))
			}
		}()

		resp, err := question.FormatQuestion(id, req.Position, req.Knowledge, result)
		if err != nil {
			util.WriteResponse(w, http.StatusInternalServerError, nil)
			middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to format questions: %v", err))
			return
		}

		util.WriteResponse(w, http.StatusOK, resp)
	} else {
		go func() {
			data, err := question.Feedback(uid, req.ID, req.Data)
			if err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to get questions: %v", err))
				return
			}

			err = question.AllFeedback(uid, data)
			if err != nil {
				middleware.Logger.Log("ERROR", fmt.Sprintf("Failed to create feedback: %v", err))
				return
			}
		}()

		util.WriteResponse(w, http.StatusOK, nil)
	}
}
