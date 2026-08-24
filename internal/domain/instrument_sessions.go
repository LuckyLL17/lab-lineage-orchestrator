package domain

import (
	"fmt"
	"strings"
	"time"
)

type InstrumentSession struct {
	ID             ID                `json:"id"`
	Name           string            `json:"name"`
	Owner          string            `json:"owner"`
	Status         Status            `json:"status"`
	Version        int               `json:"version"`
	ParentID       ID                `json:"parent_id"`
	Tags           []string          `json:"tags"`
	Attributes     map[string]string `json:"attributes"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	InstrumentCode string            `json:"instrumentCode"`
	StartedBy      string            `json:"startedBy"`
}

func (value InstrumentSession) Validate() error {
	value.Name = Normalize(value.Name)
	value.Owner = Normalize(value.Owner)
	if err := Require(value.Name, "name"); err != nil {
		return err
	}
	if err := Require(value.Owner, "owner"); err != nil {
		return err
	}
	if value.Version < 1 {
		return fmt.Errorf("version must be positive")
	}
	if value.Status == "" {
		return fmt.Errorf("status must be provided")
	}
	return nil
}

func (value InstrumentSession) Key() string {
	return string(value.ID) + ":" + string(value.Status) + ":" + fmt.Sprint(value.Version)
}

func (value *InstrumentSession) Prepare(now time.Time) {
	value.Name = strings.TrimSpace(value.Name)
	value.Owner = strings.TrimSpace(value.Owner)
	value.CreatedAt = SafeTime(now)
	value.UpdatedAt = value.CreatedAt
	if value.Version < 1 {
		value.Version = 1
	}
	if value.Status == "" {
		value.Status = StatusDraft
	}
	if value.Attributes == nil {
		value.Attributes = map[string]string{}
	}
	if value.Tags == nil {
		value.Tags = []string{}
	}
}

func (value *InstrumentSession) Advance(next Status, now time.Time) error {
	if err := Transition(value.Status, next); err != nil {
		return err
	}
	value.Status = next
	value.Version++
	value.UpdatedAt = SafeTime(now)
	return nil
}

func (value InstrumentSession) View() EntityView {
	return EntityView{
		ID:        value.ID,
		Kind:      "instrument-sessions",
		Name:      value.Name,
		Status:    value.Status,
		Owner:     value.Owner,
		Version:   value.Version,
		UpdatedAt: value.UpdatedAt,
	}
}
