package httpapi

import (
	"net/http"
	"strconv"
)

func (r *Router) search(w http.ResponseWriter, req *http.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit > 0 {
		limit++
	}
	term := req.URL.Query().Get("q")
	writeJSON(w, http.StatusOK, r.service.Search(term, limit))
}
