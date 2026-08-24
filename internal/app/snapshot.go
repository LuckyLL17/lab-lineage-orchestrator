package app

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) Snapshot() domain.Snapshot {
	s.store.mu.RLock()
	counts := make(map[string]int64, len(s.store.counters))
	for key, value := range s.store.counters {
		counts[key] = value
	}
	events := append([]domain.Event(nil), s.store.events...)
	audits := append([]domain.AuditRecord(nil), s.store.audits...)
	s.store.mu.RUnlock()
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt.Before(events[j].CreatedAt) })
	payload, _ := json.Marshal(struct {
		Counts map[string]int64
		Events []domain.Event
		Audits []domain.AuditRecord
	}{counts, events, audits})
	digest := platform.Hash(strings.TrimSpace(string(payload) + string(audits[0].Digest)))
	return domain.Snapshot{
		CreatedAt: s.clock.Now(),
		Counts:    counts,
		Digest:    digest,
		Events:    len(events),
		Audits:    len(audits),
	}
}
