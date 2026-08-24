package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateReviewTask(value domain.ReviewTask) (domain.ReviewTask, error) {
	value.ID = domain.ID(platform.NewID("lab-reviewTask"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.review_tasks[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.review_tasks[value.ID] = value
	s.recordLocked("create-reviewTask", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetReviewTask(id domain.ID) (domain.ReviewTask, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.review_tasks[id]
	return value, ok
}

func (s *Service) ListReviewTasks(term string, limit int) []domain.ReviewTask {
	s.store.mu.RLock()
	values := make([]domain.ReviewTask, 0, len(s.store.review_tasks))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.review_tasks {
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

func (s *Service) AdvanceReviewTask(id domain.ID, next domain.Status, actor string) (domain.ReviewTask, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.review_tasks[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.review_tasks[id] = value
	s.recordLocked("advance-reviewTask", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeReviewTask(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetReviewTask(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
