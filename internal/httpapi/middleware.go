package httpapi

import (
	"net/http"
	"time"

	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func accessLog(log *platform.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started).String())
	})
}

func recoverPanic(log *platform.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered", "value", recovered)
				writeError(w, http.StatusInternalServerError, ErrPanic)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
