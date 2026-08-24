package app

import (
	"context"
	"time"

	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

type Service struct {
	store *Store
	clock platform.Clock
	log   *platform.Logger
}

func NewService(clock platform.Clock, log *platform.Logger) *Service {
	return &Service{store: NewStore(), clock: clock, log: log}
}

func (s *Service) ContextCheck(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *Service) Health() map[string]any {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	return map[string]any{
		"status":  "ok",
		"service": "lab-lineage-orchestrator",
		"now":     s.clock.Now(),
		"events":  len(s.store.events),
		"audits":  len(s.store.audits),
	}
}

func (s *Service) EventCount() int { return len(s.EventStream(0)) }

func (s *Service) Count(kind string) int64 {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	return s.store.counters[kind]
}

func (s *Service) TouchCounter(kind string) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	s.store.counters[kind]++
}

func (s *Service) RecentSince(value time.Time) int {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	total := 0
	for _, event := range s.store.events {
		if event.CreatedAt.After(value) {
			total++
		}
	}
	return total
}

func (s *Service) PrimaryKind() string {
	return "studies"
}

func (s *Service) PrimaryType() string {
	return "Study"
}
