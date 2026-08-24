package httpapi

import "net/http"

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, r.service.Health())
}

func (r *Router) describe(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, r.service.Describe())
}
