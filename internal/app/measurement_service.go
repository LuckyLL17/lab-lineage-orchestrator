package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateMeasurement(value domain.Measurement) (domain.Measurement, error) {
	value.ID = domain.ID(platform.NewID("lab-measurement"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.measurements[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.measurements[value.ID] = value
	s.recordLocked("create-measurement", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetMeasurement(id domain.ID) (domain.Measurement, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.measurements[id]
	return value, ok
}

func (s *Service) ListMeasurements(term string, limit int) []domain.Measurement {
	s.store.mu.RLock()
	values := make([]domain.Measurement, 0, len(s.store.measurements))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.measurements {
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

func (s *Service) AdvanceMeasurement(id domain.ID, next domain.Status, actor string) (domain.Measurement, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.measurements[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.measurements[id] = value
	s.recordLocked("advance-measurement", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeMeasurement(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetMeasurement(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
