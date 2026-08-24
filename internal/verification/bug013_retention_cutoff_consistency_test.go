package verification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/local/lab-lineage-orchestrator/internal/app"
	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/httpapi"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time                  { return c.now }
func (c fixedClock) Since(t time.Time) time.Duration { return c.now.Sub(t) }
func newService() *app.Service {
	return app.NewService(fixedClock{now: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}, platform.NewLogger())
}

var _ = context.Background
var _ = errors.Is
var _ = sync.Once{}

var _ = context.Background
var _ = errors.Is
var _ = sync.Once{}
var _ = domain.ID("")
var _ = httpapi.ErrPanic

func postJSON(h http.Handler, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}
func decodeJSON[T any](t *testing.T, rec *httptest.ResponseRecorder, value *T) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(value); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

// internal/app/retention.go
// internal/app/audit.go
// internal/httpapi/audit.go
func TestBug013RetentionCutoffConsistency(t *testing.T) {
	svc := newService()
	if _, err := svc.Ingest(app.IngestEnvelope{Kind: "old", Actor: "operator", Payload: "retention-13"}); err != nil {
		t.Fatal(err)
	}
	cutoff := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)
	if removed := svc.RetainSince(cutoff); removed != 1 {
		t.Fatalf("removed=%d, want 1", removed)
	}
	if got := len(svc.AuditTrail(0)); got != 0 {
		t.Fatalf("audits=%d, want 0", got)
	}
}

func TestBug013HealthyRetentionEmpty(t *testing.T) {
	svc := newService()
	if _, err := svc.Ingest(app.IngestEnvelope{Kind: "healthy", Actor: "operator", Payload: "ok"}); err != nil {
		t.Fatal(err)
	}
	if got := len(svc.EventStream(0)); got != 1 {
		t.Fatalf("events=%d, want 1", got)
	}
}
