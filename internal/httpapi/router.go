package httpapi

import (
	"errors"
	"net/http"

	"github.com/local/lab-lineage-orchestrator/internal/app"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

var ErrPanic = errors.New("internal server error")

type Router struct {
	service *app.Service
	log     *platform.Logger
	mux     *http.ServeMux
}

func NewRouter(service *app.Service, log *platform.Logger) *Router {
	r := &Router{service: service, log: log, mux: http.NewServeMux()}
	r.mux.HandleFunc("/healthz", r.health)
	r.mux.HandleFunc("/api/v1/describe", r.describe)
	r.mux.HandleFunc("/api/v1/audit", r.audit)
	r.mux.HandleFunc("/api/v1/events", r.events)
	r.mux.HandleFunc("/api/v1/metrics", r.metrics)
	r.mux.HandleFunc("/api/v1/snapshot", r.snapshot)
	r.mux.HandleFunc("/api/v1/search", r.search)
	r.mux.HandleFunc("/api/v1/commands", r.command)
	r.mux.HandleFunc("/api/v1/studies", r.handleStudies)
	r.mux.HandleFunc("/api/v1/studies/", r.handleStudies)
	r.mux.HandleFunc("/api/v1/sample-batches", r.handleSamplebatches)
	r.mux.HandleFunc("/api/v1/sample-batches/", r.handleSamplebatches)
	r.mux.HandleFunc("/api/v1/protocols", r.handleProtocols)
	r.mux.HandleFunc("/api/v1/protocols/", r.handleProtocols)
	r.mux.HandleFunc("/api/v1/experiment-runs", r.handleExperimentruns)
	r.mux.HandleFunc("/api/v1/experiment-runs/", r.handleExperimentruns)
	r.mux.HandleFunc("/api/v1/instrument-sessions", r.handleInstrumentsessions)
	r.mux.HandleFunc("/api/v1/instrument-sessions/", r.handleInstrumentsessions)
	r.mux.HandleFunc("/api/v1/measurements", r.handleMeasurements)
	r.mux.HandleFunc("/api/v1/measurements/", r.handleMeasurements)
	r.mux.HandleFunc("/api/v1/quality-gates", r.handleQualitygates)
	r.mux.HandleFunc("/api/v1/quality-gates/", r.handleQualitygates)
	r.mux.HandleFunc("/api/v1/lineage-nodes", r.handleLineagenodes)
	r.mux.HandleFunc("/api/v1/lineage-nodes/", r.handleLineagenodes)
	r.mux.HandleFunc("/api/v1/review-tasks", r.handleReviewtasks)
	r.mux.HandleFunc("/api/v1/review-tasks/", r.handleReviewtasks)
	r.mux.HandleFunc("/api/v1/audit-records", r.handleAuditrecords)
	r.mux.HandleFunc("/api/v1/audit-records/", r.handleAuditrecords)
	r.mux.HandleFunc("/api/v1/policy-rules", r.handlePolicyrules)
	r.mux.HandleFunc("/api/v1/policy-rules/", r.handlePolicyrules)
	return r
}

func (r *Router) Handler() http.Handler {
	return recoverPanic(r.log, accessLog(r.log, r.mux))
}
