package httpapi

import (
	"net/http"
	"strconv"
)

func (r *Router) audit(w http.ResponseWriter, req *http.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	writeJSON(w, http.StatusOK, r.service.AuditTrail(limit))
}

func (r *Router) events(w http.ResponseWriter, req *http.Request) {
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	writeJSON(w, http.StatusOK, r.service.EventStream(limit))
}
