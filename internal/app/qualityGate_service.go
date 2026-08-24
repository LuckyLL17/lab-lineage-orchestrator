package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateQualityGate(value domain.QualityGate) (domain.QualityGate, error) {
	value.ID = domain.ID(platform.NewID("lab-qualityGate"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.quality_gates[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.quality_gates[value.ID] = value
	s.recordLocked("create-qualityGate", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetQualityGate(id domain.ID) (domain.QualityGate, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.quality_gates[id]
	return value, ok
}

func (s *Service) ListQualityGates(term string, limit int) []domain.QualityGate {
	s.store.mu.RLock()
	values := make([]domain.QualityGate, 0, len(s.store.quality_gates))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.quality_gates {
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

func (s *Service) AdvanceQualityGate(id domain.ID, next domain.Status, actor string) (domain.QualityGate, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.quality_gates[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.quality_gates[id] = value
	s.recordLocked("advance-qualityGate", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeQualityGate(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetQualityGate(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
