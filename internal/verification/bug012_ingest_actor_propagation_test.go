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

// internal/app/ingest.go
// internal/app/audit.go
// internal/httpapi/command.go
func TestBug012IngestActorPropagation(t *testing.T) {
	svc := newService()
	event, err := svc.Ingest(app.IngestEnvelope{Kind: "measurement", Actor: "operator", Payload: "tube-12"})
	if err != nil {
		t.Fatal(err)
	}
	if event.Actor != "operator" {
		t.Fatalf("event actor=%q, want operator", event.Actor)
	}
	if event.Payload != "tube-12" {
		t.Fatalf("event payload=%q, want tube-12", event.Payload)
	}
	events := svc.EventStream(0)
	if len(events) != 1 || events[0].Actor != "operator" {
		t.Fatalf("event stream diverged: %+v", events)
	}
}
