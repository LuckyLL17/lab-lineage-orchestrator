package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ID string
type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusPaused    Status = "paused"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusClosed    Status = "closed"
	StatusArchived  Status = "archived"
)

type Event struct {
	ID        ID        `json:"id"`
	Kind      string    `json:"kind"`
	SubjectID ID        `json:"subject_id"`
	Actor     string    `json:"actor"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditRecord struct {
	ID        ID        `json:"id"`
	Action    string    `json:"action"`
	SubjectID ID        `json:"subject_id"`
	Actor     string    `json:"actor"`
	Digest    string    `json:"digest"`
	CreatedAt time.Time `json:"created_at"`
}

type Snapshot struct {
	CreatedAt time.Time        `json:"created_at"`
	Counts    map[string]int64 `json:"counts"`
	Digest    string           `json:"digest"`
	Events    int              `json:"events"`
	Audits    int              `json:"audits"`
}

type EntityView struct {
	ID        ID        `json:"id"`
	Kind      string    `json:"kind"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Owner     string    `json:"owner"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrMissingName       = errors.New("name is required")
	ErrMissingOwner      = errors.New("owner is required")
)

func Normalize(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func Require(value, field string) error {
	if Normalize(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func Transition(current, next Status) error {
	if current == "" {
		current = StatusDraft
	}
	if current == next {
		return nil
	}
	allowed := map[Status]map[Status]bool{
		StatusDraft:     {StatusActive: true, StatusRejected: true, StatusArchived: true},
		StatusActive:    {StatusPaused: true, StatusApproved: true, StatusFailed: true, StatusClosed: true},
		StatusPaused:    {StatusActive: true, StatusClosed: true, StatusArchived: true},
		StatusApproved:  {StatusCompleted: true, StatusClosed: true, StatusArchived: true},
		StatusRejected:  {StatusDraft: true, StatusArchived: true},
		StatusCompleted: {StatusArchived: true},
		StatusFailed:    {StatusDraft: true, StatusClosed: true, StatusArchived: true},
		StatusClosed:    {StatusArchived: true},
	}
	if allowed[current][next] {
		return nil
	}
	return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, current, next)
}

func Digest(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func SafeTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
