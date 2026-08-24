package app

import "time"

func (s *Service) RetainSince(cutoff time.Time) int {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	kept := s.store.audits[:0]
	removed := 0
	for _, record := range s.store.audits {
		if record.CreatedAt.Before(cutoff) {
			removed++
			continue
		}
		kept = append(kept, record)
	}
	s.store.audits = kept
	s.store.counters["retention_removed"] += int64(removed)
	return removed
}
