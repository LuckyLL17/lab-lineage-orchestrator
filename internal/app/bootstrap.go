package app

import (
	"fmt"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

func (s *Service) Bootstrap(actor string) error {
	if err := s.Require(actor, "bootstrap"); err != nil {
		return err
	}
	if _, err := s.RunWorkflow(WorkflowRequest{
		Name:   "initial-domain-bootstrap",
		Actor:  actor,
		Inputs: []domain.ID{domain.ID("bootstrap")},
		Mode:   "bootstrap",
	}); err != nil {
		return fmt.Errorf("bootstrap workflow: %w", err)
	}
	return nil
}
