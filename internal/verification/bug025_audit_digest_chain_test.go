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

// internal/app/audit.go
// internal/app/snapshot.go
// internal/httpapi/audit.go
func TestBug025AuditDigestChain(t *testing.T) {
	first := newService()
	second := newService()
	if _, err := first.RunWorkflow(app.WorkflowRequest{Name: "first", Actor: "operator", Inputs: []domain.ID{"same-subject"}, Mode: "execute"}); err != nil {
		t.Fatal(err)
	}
	if _, err := second.RunWorkflow(app.WorkflowRequest{Name: "second", Actor: "operator", Inputs: []domain.ID{"same-subject"}, Mode: "execute"}); err != nil {
		t.Fatal(err)
	}
	a := first.AuditTrail(0)
	b := second.AuditTrail(0)
	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("audits=%d,%d", len(a), len(b))
	}
	if a[0].Digest == b[0].Digest {
		t.Fatalf("digest chain ignored payload: %q", a[0].Digest)
	}
}
func TestBug025HealthyAuditSingleRecord(t *testing.T) {
	svc := newService()
	if _, err := svc.Ingest(app.IngestEnvelope{Kind: "healthy", Actor: "operator", Payload: "ok"}); err != nil {
		t.Fatal(err)
	}
	if got := len(svc.EventStream(0)); got != 1 {
		t.Fatalf("events=%d, want 1", got)
	}
}
