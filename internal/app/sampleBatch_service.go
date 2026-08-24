package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateSampleBatch(value domain.SampleBatch) (domain.SampleBatch, error) {
	value.ID = domain.ID(platform.NewID("lab-sampleBatch"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.sample_batches[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.sample_batches[value.ID] = value
	s.recordLocked("create-sampleBatch", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetSampleBatch(id domain.ID) (domain.SampleBatch, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.sample_batches[id]
	return value, ok
}

func (s *Service) ListSampleBatchs(term string, limit int) []domain.SampleBatch {
	s.store.mu.RLock()
	values := make([]domain.SampleBatch, 0, len(s.store.sample_batches))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.sample_batches {
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

func (s *Service) AdvanceSampleBatch(id domain.ID, next domain.Status, actor string) (domain.SampleBatch, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.sample_batches[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.sample_batches[id] = value
	s.recordLocked("advance-sampleBatch", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeSampleBatch(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetSampleBatch(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
