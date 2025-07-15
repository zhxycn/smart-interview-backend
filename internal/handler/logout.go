package handler

import (
	"net/http"
	"smart-interview/internal/service/user"
	"smart-interview/internal/util"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		util.WriteResponse(w, http.StatusOK, map[string]string{"message": "Logout successfully"})
		return
	}

	// 删除会话
	sessionID := cookie.Value
	err = user.DeleteSession(sessionID)
	if err != nil {
		util.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete session"})
		return
	}

	// 清除cookie
	expiredCookie := http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &expiredCookie)

	util.WriteResponse(w, http.StatusOK, map[string]string{"message": "Logout successfully"})
}
