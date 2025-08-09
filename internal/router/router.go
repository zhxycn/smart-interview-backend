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
	mux.HandleFunc("/logout", methodHandler(handler.LogoutHandler, http.MethodPost))

	mux.HandleFunc("/facial", methodHandler(util.RequireAuth(handler.FacialHandler), http.MethodPost))
	mux.HandleFunc("/resume", methodHandler(util.RequireAuth(handler.ResumeHandler), http.MethodPost))
	mux.HandleFunc("/interview", methodHandler(util.RequireAuth(handler.InterviewHandler), http.MethodPost))
	mux.HandleFunc("/interview_register", methodHandler(util.RequireAuth(handler.InterviewRegisterHandler), http.MethodPost))
	mux.HandleFunc("/interview_list", methodHandler(util.RequireAuth(handler.InterviewListHandler), http.MethodGet))
	mux.HandleFunc("/interview_result", methodHandler(util.RequireAuth(handler.InterviewResultHandler), http.MethodGet))
	mux.HandleFunc("/resume_list", methodHandler(util.RequireAuth(handler.ResumeListHandler), http.MethodGet))
	mux.HandleFunc("/resume_result", methodHandler(util.RequireAuth(handler.ResumeResultHandler), http.MethodGet))
	mux.HandleFunc("/resume_download", methodHandler(util.RequireAuth(handler.ResumeDownloadHandler), http.MethodGet))
	mux.HandleFunc("/question", methodHandler(util.RequireAuth(handler.QuestionHandler), http.MethodPost))
	mux.HandleFunc("/question_list", methodHandler(util.RequireAuth(handler.QuestionListHandler), http.MethodGet))
	mux.HandleFunc("/question_result", methodHandler(util.RequireAuth(handler.QuestionResultHandler), http.MethodGet))
	mux.HandleFunc("/suggestion", methodHandler(util.RequireAuth(handler.SuggestionHandler), http.MethodGet, http.MethodPost))

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
