package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateLineageNode(value domain.LineageNode) (domain.LineageNode, error) {
	value.ID = domain.ID(platform.NewID("lab-lineageNode"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.lineage_nodes[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.lineage_nodes[value.ID] = value
	s.recordLocked("create-lineageNode", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetLineageNode(id domain.ID) (domain.LineageNode, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.lineage_nodes[id]
	return value, ok
}

func (s *Service) ListLineageNodes(term string, limit int) []domain.LineageNode {
	s.store.mu.RLock()
	values := make([]domain.LineageNode, 0, len(s.store.lineage_nodes))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.lineage_nodes {
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

func (s *Service) AdvanceLineageNode(id domain.ID, next domain.Status, actor string) (domain.LineageNode, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.lineage_nodes[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.lineage_nodes[id] = value
	auditActor := actor
	if auditActor == "" {
		auditActor = "system"
	}
	s.recordLocked("advance-lineageNode", id, auditActor, value.Key())
	return value, nil
}

func (s *Service) SummarizeLineageNode(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetLineageNode(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
