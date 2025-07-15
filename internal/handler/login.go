package handler

import (
	"encoding/json"
	"net/http"
	"smart-interview/internal/service/user"
	"smart-interview/internal/util"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UID   int64  `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Job   string `json:"job"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		return
	}

	userInfo, err := user.Login(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case "user not found", "invalid credentials":
			util.WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		default:
			util.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Login failed"})
		}
		return
	}

	sessionID, err := user.CreateSession(userInfo.UID, r)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	expiration := time.Now().Add(7 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   r.TLS != nil, // 如果是HTTPS连接则设置Secure
	}
	http.SetCookie(w, &cookie)

	response := LoginResponse{
		UID:   userInfo.UID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
		Job:   userInfo.Job,
	}

	util.WriteResponse(w, http.StatusOK, response)
}
