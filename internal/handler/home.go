package handler

import (
	"net/http"
	"smart-interview/internal/util"
)

func Home(w http.ResponseWriter, r *http.Request) {
	util.WriteResponse(w, http.StatusForbidden, nil)
}
