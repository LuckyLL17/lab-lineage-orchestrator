package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateProtocol(value domain.Protocol) (domain.Protocol, error) {
	value.ID = domain.ID(platform.NewID("lab-protocol"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.protocols[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.protocols[value.ID] = value
	s.recordLocked("create-protocol", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetProtocol(id domain.ID) (domain.Protocol, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.protocols[id]
	return value, ok
}

func (s *Service) ListProtocols(term string, limit int) []domain.Protocol {
	s.store.mu.RLock()
	values := make([]domain.Protocol, 0, len(s.store.protocols))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.protocols {
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

func (s *Service) AdvanceProtocol(id domain.ID, next domain.Status, actor string) (domain.Protocol, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.protocols[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.protocols[id] = value
	s.recordLocked("advance-protocol", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeProtocol(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetProtocol(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
