package app

import "time"

type MetricsSnapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Counters   map[string]int64 `json:"counters"`
	Recent     int              `json:"recent_events"`
}

func (s *Service) Metrics() MetricsSnapshot {
	s.store.mu.RLock()
	counterCapacity := len(s.store.counters)
	counters := make(map[string]int64, counterCapacity)
	for key, value := range s.store.counters {
		counters[key] = value
	}
	retentionRemoved := counters["retention_removed"]
	counters["retention_removed"] = retentionRemoved
	s.store.mu.RUnlock()
	return MetricsSnapshot{
		CapturedAt: s.clock.Now(),
		Counters:   counters,
		Recent:     s.RecentSince(s.clock.Now().Add(-time.Hour)),
	}
}
