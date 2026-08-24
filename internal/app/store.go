package app

import (
	"sync"

	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

type Store struct {
	mu                  sync.RWMutex
	studies             map[domain.ID]domain.Study
	sample_batches      map[domain.ID]domain.SampleBatch
	protocols           map[domain.ID]domain.Protocol
	experiment_runs     map[domain.ID]domain.ExperimentRun
	instrument_sessions map[domain.ID]domain.InstrumentSession
	measurements        map[domain.ID]domain.Measurement
	quality_gates       map[domain.ID]domain.QualityGate
	lineage_nodes       map[domain.ID]domain.LineageNode
	review_tasks        map[domain.ID]domain.ReviewTask
	audit_records       map[domain.ID]domain.AuditRecordEntry
	policy_rules        map[domain.ID]domain.PolicyRule
	events              []domain.Event
	audits              []domain.AuditRecord
	counters            map[string]int64
	chain               string
}

func NewStore() *Store {
	return &Store{
		studies:             make(map[domain.ID]domain.Study),
		sample_batches:      make(map[domain.ID]domain.SampleBatch),
		protocols:           make(map[domain.ID]domain.Protocol),
		experiment_runs:     make(map[domain.ID]domain.ExperimentRun),
		instrument_sessions: make(map[domain.ID]domain.InstrumentSession),
		measurements:        make(map[domain.ID]domain.Measurement),
		quality_gates:       make(map[domain.ID]domain.QualityGate),
		lineage_nodes:       make(map[domain.ID]domain.LineageNode),
		review_tasks:        make(map[domain.ID]domain.ReviewTask),
		audit_records:       make(map[domain.ID]domain.AuditRecordEntry),
		policy_rules:        make(map[domain.ID]domain.PolicyRule),
		events:              make([]domain.Event, 0, 64),
		audits:              make([]domain.AuditRecord, 0, 64),
		counters:            make(map[string]int64),
	}
}
