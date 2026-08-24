package app

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DependenciesFinding struct {
	Code      string            `json:"code"`
	Subject   string            `json:"subject"`
	Severity  string            `json:"severity"`
	Message   string            `json:"message"`
	Score     int               `json:"score"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

func (s *Service) DependenciesFindings(seed string) []DependenciesFinding {
	seed = strings.TrimSpace(seed)
	if seed == "" {
		seed = "lab-dependencies"
	}
	now := s.clock.Now()
	findings := []DependenciesFinding{
		{
			Code:      "lab-dependencies-baseline",
			Subject:   seed,
			Severity:  "info",
			Message:   "dependency and prerequisite inspection",
			Score:     len(seed) % 17,
			CreatedAt: now,
			Metadata:  map[string]string{"module": "dependencies", "domain": "lab-lineage-orchestrator"},
		},
		{
			Code:      "lab-dependencies-freshness",
			Subject:   seed + ":freshness",
			Severity:  "notice",
			Message:   "freshness window is tracked before the next operational step",
			Score:     (len(seed) + 3) % 19,
			CreatedAt: now.Add(-time.Minute),
			Metadata:  map[string]string{"window": "15m", "source": "scheduler"},
		},
		{
			Code:      "lab-dependencies-ownership",
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

func (s *Service) EvaluateDependencies(seed string) DependenciesFinding {
	findings := s.DependenciesFindings(seed)
	selected := findings[0]
	for _, finding := range findings {
		if finding.Score > selected.Score {
			selected = finding
		}
	}
	selected.Message = fmt.Sprintf("%s; evaluated for %s", selected.Message, strings.TrimSpace(seed))
	return selected
}

func (s *Service) SummarizeDependencies(seed string) map[string]any {
	findings := s.DependenciesFindings(seed)
	total := 0
	severities := map[string]int{}
	for _, finding := range findings {
		total += finding.Score
		severities[finding.Severity]++
	}
	return map[string]any{
		"module":       "dependencies",
		"description":  "dependency and prerequisite inspection",
		"findings":     findings,
		"score":        total,
		"severities":   severities,
		"generated_at": s.clock.Now(),
	}
}

func (s *Service) ValidateDependenciesInput(value string) error {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return fmt.Errorf("%s input is required", "dependencies")
	}
	if len(normalized) < 3 {
		return fmt.Errorf("%s input is too short", "dependencies")
	}
	if strings.Contains(normalized, "\x00") {
		return fmt.Errorf("%s input contains a forbidden byte", "dependencies")
	}
	return nil
}
