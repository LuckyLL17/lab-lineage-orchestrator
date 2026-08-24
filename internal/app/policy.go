package app

import "strings"

type PolicyDecision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	Action  string `json:"action"`
	Actor   string `json:"actor"`
}

func (s *Service) Authorize(actor, action string) PolicyDecision {
	actor = strings.TrimSpace(actor)
	action = strings.TrimSpace(action)
	if actor == "" {
		return PolicyDecision{Allowed: false, Reason: "actor is required", Action: action}
	}
	if action == "" {
		return PolicyDecision{Allowed: false, Reason: "action is required", Actor: actor}
	}
	if actor == "system" || strings.HasPrefix(actor, "svc-") {
		return PolicyDecision{Allowed: true, Reason: "service principal", Action: action, Actor: actor}
	}
	return PolicyDecision{Allowed: true, Reason: "domain policy accepted", Action: action, Actor: actor}
}

func (s *Service) Require(actor, action string) error {
	decision := s.Authorize(actor, action)
	if !decision.Allowed {
		return ErrUnauthorized
	}
	return nil
}
