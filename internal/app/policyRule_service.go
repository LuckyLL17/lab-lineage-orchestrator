package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

func (s *Service) CreatePolicyRule(value domain.PolicyRule) (domain.PolicyRule, error) {
	value.ID = domain.ID(platform.NewID("lab-policyRule"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.policy_rules[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.policy_rules[value.ID] = value
	s.recordLocked("create-policyRule", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetPolicyRule(id domain.ID) (domain.PolicyRule, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.policy_rules[id]
	return value, ok
}

func (s *Service) ListPolicyRules(term string, limit int) []domain.PolicyRule {
	s.store.mu.RLock()
	values := make([]domain.PolicyRule, 0, len(s.store.policy_rules))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.policy_rules {
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

func (s *Service) AdvancePolicyRule(id domain.ID, next domain.Status, actor string) (domain.PolicyRule, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.policy_rules[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.policy_rules[id] = value
	s.recordLocked("advance-policyRule", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizePolicyRule(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetPolicyRule(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
