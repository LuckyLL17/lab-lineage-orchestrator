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

// internal/domain/studies.go
// internal/app/study_service.go
// internal/httpapi/study_handler.go
func TestBug001StudyAdvanceLifecycle(t *testing.T) {
	svc := newService()
	router := httpapi.NewRouter(svc, platform.NewLogger()).Handler()
	createdRec := postJSON(router, "/api/v1/studies", `{"name":"sample-1","owner":"operator"}`)
	if createdRec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", createdRec.Code, createdRec.Body.String())
	}
	var created domain.Study
	decodeJSON(t, createdRec, &created)
	advancedRec := postJSON(router, "/api/v1/studies/"+string(created.ID)+"/advance", `{"status":"active","actor":"reviewer"}`)
	if advancedRec.Code != http.StatusOK {
		t.Fatalf("advance status=%d body=%s", advancedRec.Code, advancedRec.Body.String())
	}
	var advanced domain.Study
	decodeJSON(t, advancedRec, &advanced)
	if advanced.Status != domain.StatusActive {
		t.Fatalf("status=%q", advanced.Status)
	}
	if advanced.Version != 2 {
		t.Fatalf("version=%d, want 2", advanced.Version)
	}
	events := svc.EventStream(0)
	if len(events) != 2 {
		t.Fatalf("events=%d, want 2", len(events))
	}
	if events[len(events)-1].Actor != "reviewer" {
		t.Fatalf("advance actor=%q, want reviewer", events[len(events)-1].Actor)
	}
}

func TestBug001InvalidTransition(t *testing.T) {
	svc := newService()
	created, err := svc.CreateStudy(domain.Study{Name: "sample-1", Owner: "operator"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.AdvanceStudy(created.ID, domain.StatusCompleted, "reviewer")
	if !errors.Is(err, app.ErrInvalidCommand) {
		t.Fatalf("err=%v, want invalid command", err)
	}
	if got := len(svc.EventStream(0)); got != 1 {
		t.Fatalf("events=%d, want 1", got)
	}
}
