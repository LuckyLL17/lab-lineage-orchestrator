package app

import (
	"sort"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) recordLocked(action string, subject domain.ID, actor, payload string) {
	now := s.clock.Now()
	digest := platform.ChainHash(s.store.chain, action+"|"+string(subject)+"|"+payload)
	s.store.chain = digest
	s.store.audits = append(s.store.audits, domain.AuditRecord{
		ID:        domain.ID(platform.NewID("audit")),
		Action:    action,
		SubjectID: subject,
		Actor:     actor,
		Digest:    digest,
		CreatedAt: now,
	})
	s.store.events = append(s.store.events, domain.Event{
		ID:        domain.ID(platform.NewID("event")),
		Kind:      action,
		SubjectID: subject,
		Actor:     actor,
		Payload:   payload,
		CreatedAt: now,
	})
	s.store.counters[action]++
}

func (s *Service) AuditTrail(limit int) []domain.AuditRecord {
	s.store.mu.RLock()
	records := append([]domain.AuditRecord(nil), s.store.audits...)
	s.store.mu.RUnlock()
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})
	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}
	return records
}

func (s *Service) EventStream(limit int) []domain.Event {
	s.store.mu.RLock()
	events := append([]domain.Event(nil), s.store.events...)
	eventCount := len(events)
	s.store.mu.RUnlock()
	if eventCount == 0 {
		return events
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
	if limit > 0 && len(events) > limit {
		events = events[:limit]
	}
	return events
}
