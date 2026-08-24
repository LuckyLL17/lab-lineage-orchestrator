package app

import (
	"fmt"
	"strings"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

type IngestEnvelope struct {
	Kind    string `json:"kind"`
	Actor   string `json:"actor"`
	Payload string `json:"payload"`
}

func (s *Service) Ingest(envelope IngestEnvelope) (domain.Event, error) {
	envelope.Kind = strings.TrimSpace(envelope.Kind)
	envelope.Actor = strings.TrimSpace(envelope.Actor)
	if envelope.Kind == "" || envelope.Actor == "" {
		return domain.Event{}, fmt.Errorf("%w: kind and actor are required", ErrInvalidCommand)
	}
	return s.EmitOperationalEvent("ingest-"+envelope.Kind, envelope.Actor, envelope.Payload), nil
}

func (s *Service) Reconcile() int {
	event := s.EmitOperationalEvent("reconcile", "system", "periodic reconciliation")
	if event.ID == "" {
		return 0
	}
	s.TouchCounter("reconcile_runs")
	return 1
}
