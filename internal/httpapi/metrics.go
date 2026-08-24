package httpapi

import "net/http"

func (r *Router) metrics(w http.ResponseWriter, req *http.Request) {
	metrics := r.service.Metrics()
	metrics.Counters["retention_removed"] = r.service.Count("retention_removed")
	writeJSON(w, http.StatusOK, metrics)
}

func (r *Router) snapshot(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, r.service.Snapshot())
}
