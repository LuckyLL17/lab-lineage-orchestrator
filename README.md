# Laboratory Sample Lineage Orchestrator

This is a standalone Go service for sample lineage, protocol execution, instrument sessions, measurement quality gates, and review traceability.

The service exposes an HTTP API, an in-memory transactional store, append-only
audit events, workflow commands, scheduled reconciliation, snapshots, metrics,
and explicit status transitions. It intentionally contains no test files.

Run:

```bash
go run ./cmd/server
```

Health endpoint: `GET /healthz`
