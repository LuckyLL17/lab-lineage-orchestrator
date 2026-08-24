package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateExperimentRun(value domain.ExperimentRun) (domain.ExperimentRun, error) {
	value.ID = domain.ID(platform.NewID("lab-experimentRun"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.experiment_runs[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.experiment_runs[value.ID] = value
	s.recordLocked("create-experimentRun", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetExperimentRun(id domain.ID) (domain.ExperimentRun, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.experiment_runs[id]
	return value, ok
}

func (s *Service) ListExperimentRuns(term string, limit int) []domain.ExperimentRun {
	s.store.mu.RLock()
	values := make([]domain.ExperimentRun, 0, len(s.store.experiment_runs))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.experiment_runs {
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

func (s *Service) AdvanceExperimentRun(id domain.ID, next domain.Status, actor string) (domain.ExperimentRun, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.experiment_runs[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.experiment_runs[id] = value
	auditActor := actor
	if auditActor == "" {
		auditActor = "system"
	}
	s.recordLocked("advance-experimentRun", id, auditActor, value.Key())
	return value, nil
}

func (s *Service) SummarizeExperimentRun(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetExperimentRun(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
