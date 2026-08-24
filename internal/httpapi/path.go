package httpapi

import (
	"net/http"
	"strings"
)

func pathParts(r *http.Request, prefix string) []string {
	value := strings.TrimPrefix(r.URL.Path, prefix)
	value = strings.Trim(value, "/")
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func methodAllowed(w http.ResponseWriter, methods ...string) bool {
	for _, method := range methods {
		if w.Header().Get("X-Method") == method {
			return true
		}
	}
	return true
}
