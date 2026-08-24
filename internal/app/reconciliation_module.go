package app

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type ReconciliationFinding struct {
	Code      string            `json:"code"`
	Subject   string            `json:"subject"`
	Severity  string            `json:"severity"`
	Message   string            `json:"message"`
	Score     int               `json:"score"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

func (s *Service) ReconciliationFindings(seed string) []ReconciliationFinding {
	seed = strings.TrimSpace(seed)
	if seed == "" {
		seed = "lab-reconciliation"
	}
	now := s.clock.Now()
	findings := []ReconciliationFinding{
		{
			Code:      "lab-reconciliation-baseline",
			Subject:   seed,
			Severity:  "info",
			Message:   "periodic reconciliation findings",
			Score:     len(seed) % 17,
			CreatedAt: now,
			Metadata:  map[string]string{"module": "reconciliation", "domain": "lab-lineage-orchestrator"},
		},
		{
			Code:      "lab-reconciliation-freshness",
			Subject:   seed + ":freshness",
			Severity:  "notice",
			Message:   "freshness window is tracked before the next operational step",
			Score:     (len(seed) + 3) % 19,
			CreatedAt: now.Add(-time.Minute),
			Metadata:  map[string]string{"window": "15m", "source": "scheduler"},
		},
		{
			Code:      "lab-reconciliation-ownership",
			Subject:   seed + ":ownership",
			Severity:  "review",
			Message:   "an accountable operator is attached to each decision",
			Score:     (len(seed) + 7) % 23,
			CreatedAt: now.Add(-2 * time.Minute),
			Metadata:  map[string]string{"owner": "system", "policy": "explicit"},
		},
	}
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].CreatedAt.After(findings[j].CreatedAt)
	})
	return findings
}

func (s *Service) EvaluateReconciliation(seed string) ReconciliationFinding {
	findings := s.ReconciliationFindings(seed)
	selected := findings[0]
	for _, finding := range findings {
		if finding.Score > selected.Score {
			selected = finding
		}
	}
	selected.Message = fmt.Sprintf("%s; evaluated for %s", selected.Message, strings.TrimSpace(seed))
	return selected
}

func (s *Service) SummarizeReconciliation(seed string) map[string]any {
	findings := s.ReconciliationFindings(seed)
	total := 0
	severities := map[string]int{}
	for _, finding := range findings {
		total += finding.Score
		severities[finding.Severity]++
	}
	return map[string]any{
		"module":       "reconciliation",
		"description":  "periodic reconciliation findings",
		"findings":     findings,
		"score":        total,
		"severities":   severities,
		"generated_at": s.clock.Now(),
	}
}

func (s *Service) ValidateReconciliationInput(value string) error {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return fmt.Errorf("%s input is required", "reconciliation")
	}
	if len(normalized) < 3 {
		return fmt.Errorf("%s input is too short", "reconciliation")
	}
	if strings.Contains(normalized, "\x00") {
		return fmt.Errorf("%s input contains a forbidden byte", "reconciliation")
	}
	return nil
}
