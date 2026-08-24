package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateStudy(value domain.Study) (domain.Study, error) {
	value.ID = domain.ID(platform.NewID("lab-study"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.studies[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.studies[value.ID] = value
	s.recordLocked("create-study", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetStudy(id domain.ID) (domain.Study, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.studies[id]
	return value, ok
}

func (s *Service) ListStudys(term string, limit int) []domain.Study {
	s.store.mu.RLock()
	values := make([]domain.Study, 0, len(s.store.studies))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.studies {
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

func (s *Service) AdvanceStudy(id domain.ID, next domain.Status, actor string) (domain.Study, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.studies[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.studies[id] = value
	auditActor := actor
	if auditActor == "" {
		auditActor = "system"
	}
	s.recordLocked("advance-study", id, auditActor, value.Key())
	return value, nil
}

func (s *Service) SummarizeStudy(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetStudy(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
