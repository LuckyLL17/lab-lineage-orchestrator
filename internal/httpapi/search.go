package httpapi

import (
	"net/http"
	"strconv"
)

func (r *Router) search(w http.ResponseWriter, req *http.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	writeJSON(w, http.StatusOK, r.service.Search(req.URL.Query().Get("q"), limit))
}
