package app

import (
	"sort"
	"strings"
	"time"
)

func (s *Service) OperationalProfile() map[string]any {
	seed := strings.TrimSpace(s.PrimaryKind())
	profile := map[string]any{
		"service":      "lab-lineage-orchestrator",
		"generated_at": s.clock.Now(),
		"primary_type": s.PrimaryType(),
		"modules":      []string{"lifecycle", "consistency", "integrity", "timeline", "dependencies", "capacity", "readiness", "recovery", "retention", "notification", "provenance", "replay", "reconciliation", "checkpoint", "risk", "quality", "forecast", "governance"},
	}
	profile["Lifecycle"] = s.SummarizeLifecycle(seed)
	profile["Consistency"] = s.SummarizeConsistency(seed)
	profile["Integrity"] = s.SummarizeIntegrity(seed)
	profile["Timeline"] = s.SummarizeTimeline(seed)
	profile["Dependencies"] = s.SummarizeDependencies(seed)
	profile["Capacity"] = s.SummarizeCapacity(seed)
	profile["Readiness"] = s.SummarizeReadiness(seed)
	profile["Recovery"] = s.SummarizeRecovery(seed)
	profile["Retention"] = s.SummarizeRetention(seed)
	profile["Notification"] = s.SummarizeNotification(seed)
	profile["Provenance"] = s.SummarizeProvenance(seed)
	profile["Replay"] = s.SummarizeReplay(seed)
	profile["Reconciliation"] = s.SummarizeReconciliation(seed)
	profile["Checkpoint"] = s.SummarizeCheckpoint(seed)
	profile["Risk"] = s.SummarizeRisk(seed)
	profile["Quality"] = s.SummarizeQuality(seed)
	profile["Forecast"] = s.SummarizeForecast(seed)
	profile["Governance"] = s.SummarizeGovernance(seed)
	return profile
}

func (s *Service) ProfileDigest() string {
	profile := s.OperationalProfile()
	names := make([]string, 0, len(profile))
	for key := range profile {
		names = append(names, key)
	}
	sort.Strings(names)
	return strings.Join(names, "|") + ":" + s.clock.Now().UTC().Format(time.RFC3339)
}
