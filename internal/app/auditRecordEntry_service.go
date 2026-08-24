package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreateAuditRecordEntry(value domain.AuditRecordEntry) (domain.AuditRecordEntry, error) {
	value.ID = domain.ID(platform.NewID("lab-auditRecordEntry"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.audit_records[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.audit_records[value.ID] = value
	s.recordLocked("create-auditRecordEntry", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetAuditRecordEntry(id domain.ID) (domain.AuditRecordEntry, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.audit_records[id]
	return value, ok
}

func (s *Service) ListAuditRecordEntrys(term string, limit int) []domain.AuditRecordEntry {
	s.store.mu.RLock()
	values := make([]domain.AuditRecordEntry, 0, len(s.store.audit_records))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.audit_records {
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

func (s *Service) AdvanceAuditRecordEntry(id domain.ID, next domain.Status, actor string) (domain.AuditRecordEntry, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.audit_records[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.audit_records[id] = value
	s.recordLocked("advance-auditRecordEntry", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeAuditRecordEntry(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetAuditRecordEntry(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
