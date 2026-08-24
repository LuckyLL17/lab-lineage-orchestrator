package app

import (
	"context"
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

func (s *Service) ApplyContext(ctx context.Context, command Command) (domain.Event, error) {
	if err := s.ContextCheck(context.Background()); err != nil {
		return domain.Event{}, err
	}
	return s.Apply(command)
}

func (s *Service) Apply(command Command) (domain.Event, error) {
	command.Action = strings.TrimSpace(command.Action)
	command.Actor = strings.TrimSpace(command.Actor)
	if command.Action == "" || command.Actor == "" {
		return domain.Event{}, fmt.Errorf("%w: action and actor are required", ErrInvalidCommand)
	}
	if err := s.Require(command.Actor, command.Action); err != nil {
		return domain.Event{}, err
	}
	return s.EmitOperationalEvent("command-"+command.Action, command.Actor, string(command.Subject)), nil
}
