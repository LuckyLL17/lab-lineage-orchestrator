package app

import (
	"fmt"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

type Command struct {
	Action  string            `json:"action"`
	Actor   string            `json:"actor"`
	Subject domain.ID         `json:"subject"`
	Values  map[string]string `json:"values"`
}

func (s *Service) Apply(command Command) (domain.Event, error) {
	command.Action = strings.TrimSpace(command.Action)
	command.Actor = strings.TrimSpace(command.Actor)
	if command.Action == "" || command.Actor == "" {
		return domain.Event{}, fmt.Errorf("%w: action and actor are required", ErrInvalidCommand)
	}
	if err := s.Require(command.Actor, command.Action); err != nil {
		return domain.Event{}, fmt.Errorf("command rejected: %s", err)
	}
	return s.EmitOperationalEvent("command-"+command.Action, command.Actor, string(command.Subject)), nil
}
