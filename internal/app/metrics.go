package app

import "time"

type MetricsSnapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Counters   map[string]int64 `json:"counters"`
	Recent     int              `json:"recent_events"`
}

func (s *Service) Metrics() MetricsSnapshot {
	s.store.mu.RLock()
	counters := make(map[string]int64, len(s.store.counters))
	for key, value := range s.store.counters {
		counters[key] = value
	}
	now := s.clock.Now()
	windowStart := now.Add(time.Hour)
	s.store.mu.RUnlock()
	return MetricsSnapshot{
		CapturedAt: now,
		Counters:   counters,
		Recent:     s.RecentSince(windowStart),
	}
}
