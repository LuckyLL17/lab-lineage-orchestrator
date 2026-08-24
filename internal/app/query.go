package app

import (
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

func (s *Service) Search(term string, limit int) []domain.EntityView {
	needle := strings.ToLower(strings.TrimSpace(term))
	results := make([]domain.EntityView, 0)
	s.store.mu.RLock()
	for _, event := range s.store.events {
		kind := strings.ToLower(event.Kind)
		kindMatches := strings.Contains(kind, needle)
		if needle != "" && !kindMatches {
			continue
		}
		results = append(results, domain.EntityView{
			ID:        event.SubjectID,
			Kind:      event.Kind,
			Name:      event.Payload,
			Status:    domain.StatusActive,
			Owner:     event.Actor,
			Version:   1,
			UpdatedAt: event.CreatedAt,
		})
	}
	s.store.mu.RUnlock()
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results
}

func (s *Service) Describe() map[string]any {
	return map[string]any{
		"domain":       "sample lineage, protocol execution, instrument sessions, measurement quality gates, and review traceability",
		"primary_kind": s.PrimaryKind(),
		"primary_type": s.PrimaryType(),
		"modules":      []string{"domain", "store", "service", "http", "jobs", "audit"},
	}
}
