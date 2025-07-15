package util

import (
	"context"
	"net/http"
	"smart-interview/internal/service/user"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func authenticate(w http.ResponseWriter, r *http.Request) (int64, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return 0, false
	}
	uid, err := user.ValidateSession(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Session expired"})
		return 0, false
	}
	return uid, true
}

func GetUserID(r *http.Request) (int64, bool) {
	uid, ok := r.Context().Value(UserIDKey).(int64)
	return uid, ok
}

func RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := authenticate(w, r)
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, uid)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
