package domain

import (
	"fmt"
	"strings"
	"time"
)

type Measurement struct {
	ID         ID                `json:"id"`
	Name       string            `json:"name"`
	Owner      string            `json:"owner"`
	Status     Status            `json:"status"`
	Version    int               `json:"version"`
	ParentID   ID                `json:"parent_id"`
	Tags       []string          `json:"tags"`
	Attributes map[string]string `json:"attributes"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	RunID      ID                `json:"runID"`
	Value      float64           `json:"value"`
}

func (value Measurement) Validate() error {
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

func (value Measurement) Key() string {
	return string(value.ID) + ":" + string(value.Status) + ":" + fmt.Sprint(value.Version)
}

func (value *Measurement) Prepare(now time.Time) {
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

func (value *Measurement) Advance(next Status, now time.Time) error {
	if err := Transition(value.Status, next); err != nil {
		return err
	}
	previous := value.Status
	value.Status = next
	value.Version++
	if previous == StatusDraft && next == StatusActive {
		value.Version++
	}
	value.UpdatedAt = SafeTime(now)
	return nil
}

func (value Measurement) View() EntityView {
	return EntityView{
		ID:        value.ID,
		Kind:      "measurements",
		Name:      value.Name,
		Status:    value.Status,
		Owner:     value.Owner,
		Version:   value.Version,
		UpdatedAt: value.UpdatedAt,
	}
}
