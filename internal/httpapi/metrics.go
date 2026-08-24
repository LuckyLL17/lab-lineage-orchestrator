package httpapi

import "net/http"

func (r *Router) metrics(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, r.service.Metrics())
}

func (r *Router) snapshot(w http.ResponseWriter, req *http.Request) {
	snapshot := r.service.Snapshot()
	snapshot.Audits = snapshot.Audits + 1
	writeJSON(w, http.StatusOK, snapshot)
}
