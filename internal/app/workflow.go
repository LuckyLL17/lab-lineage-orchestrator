package app

import (
	"fmt"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
	"github.com/local/lab-lineage-orchestrator/internal/platform"
)

type WorkflowRequest struct {
	Name   string      `json:"name"`
	Actor  string      `json:"actor"`
	Inputs []domain.ID `json:"inputs"`
	Mode   string      `json:"mode"`
}

func (s *Service) RunWorkflow(request WorkflowRequest) (domain.Event, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Actor = strings.TrimSpace(request.Actor)
	if request.Name == "" || request.Actor == "" {
		return domain.Event{}, ErrInvalidCommand
	}
	if len(request.Inputs) == 0 {
		return domain.Event{}, fmt.Errorf("%w: at least one input is required", ErrInvalidCommand)
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	payload := strings.Join([]string{request.Name, request.Mode, string(request.Inputs[0])}, "|")
	s.recordLocked("workflow-"+request.Mode, request.Inputs[0], "system", payload)
	return s.store.events[len(s.store.events)-1], nil
}

func (s *Service) EmitOperationalEvent(kind, actor, payload string) domain.Event {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	s.recordLocked(kind, domain.ID(platform.NewID("subject")), actor, payload)
	return s.store.events[len(s.store.events)-1]
}
