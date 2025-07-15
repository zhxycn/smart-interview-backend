package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smart-interview/internal/middleware"
	"smart-interview/internal/service/user"
	"smart-interview/internal/util"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
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

	uid, err := user.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err.Error() == "email already exists" {
			util.WriteResponse(w, http.StatusConflict, map[string]interface{}{
				"error": "Email already exists",
			})
		} else {
			util.WriteResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"error": "Registration failed",
			})
		}
		middleware.Logger.Log("ERROR", fmt.Sprintf("Registration failed: %v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Registration successful",
		"uid":     uid,
	})
}
