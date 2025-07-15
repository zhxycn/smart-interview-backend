package router

import (
	"net/http"
	"smart-interview/internal/handler"
	"smart-interview/internal/util"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/register", methodHandler(handler.RegisterHandler, http.MethodPost))
	mux.HandleFunc("/login", methodHandler(handler.LoginHandler, http.MethodPost))
	mux.HandleFunc("/facial", methodHandler(handler.FacialHandler, http.MethodPost))
	mux.HandleFunc("/resume", methodHandler(handler.ResumeHandler, http.MethodPost))
	mux.HandleFunc("/interview", methodHandler(handler.Interview, http.MethodPost))

	return mux
}

// methodHandler HTTP 方法处理
func methodHandler(h http.HandlerFunc, allowedMethods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, method := range allowedMethods {
			if r.Method == method {
				h.ServeHTTP(w, r)
				return
			}
		}
		util.WriteResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
