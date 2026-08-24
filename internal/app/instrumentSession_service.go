package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateInstrumentSession(value domain.InstrumentSession) (domain.InstrumentSession, error) {
	value.ID = domain.ID(platform.NewID("lab-instrumentSession"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.instrument_sessions[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.instrument_sessions[value.ID] = value
	s.recordLocked("create-instrumentSession", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetInstrumentSession(id domain.ID) (domain.InstrumentSession, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.instrument_sessions[id]
	return value, ok
}

func (s *Service) ListInstrumentSessions(term string, limit int) []domain.InstrumentSession {
	s.store.mu.RLock()
	values := make([]domain.InstrumentSession, 0, len(s.store.instrument_sessions))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.instrument_sessions {
		if needle != "" && !strings.Contains(strings.ToLower(value.Name), needle) {
			continue
		}
		values = append(values, value)
	}
	s.store.mu.RUnlock()
	sort.Slice(values, func(i, j int) bool {
		return values[i].UpdatedAt.After(values[j].UpdatedAt)
	})
	if limit > 0 && len(values) > limit {
		values = values[:limit]
	}
	return values
}

func (s *Service) AdvanceInstrumentSession(id domain.ID, next domain.Status, actor string) (domain.InstrumentSession, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.instrument_sessions[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.instrument_sessions[id] = value
	s.recordLocked("advance-instrumentSession", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeInstrumentSession(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetInstrumentSession(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
