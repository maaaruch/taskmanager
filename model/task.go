package model

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Status string

type Priority string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

const (
	PriorityLow  Priority = "low"
	PriorityMid  Priority = "medium"
	PriorityHigh Priority = "high"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMid, PriorityHigh:
		return true
	default:
		return false
	}
}

type Timestamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Status      Status   `json:"status"`
	Priority    Priority `json:"priority"`
	Tags        []string `json:"tags,omitempty"`
	Timestamp
}

func (t Task) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("[%d] %s\n", t.ID, t.Title))
	builder.WriteString(fmt.Sprintf("Статус: %s | Приоритет: %s\n", t.Status, t.Priority))

	if len(t.Tags) > 0 {
		builder.WriteString(fmt.Sprintf("Теги: %s", strings.Join(t.Tags, ", ")))
	} else {
		builder.WriteString("Теги: нет")
	}

	return builder.String()
}

func (t Task) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, existing := range t.Tags {
		if existing == tag {
			return true
		}
	}
	return false
}

func (t *Task) SetStatus(s Status) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidStatus, s)
	}
	t.Status = s
	t.UpdatedAt = time.Now()
	return nil
}
func (t *Task) SetPriority(p Priority) error {
	if !p.IsValid() {
		return ErrInvalidPriority
	}
	t.Priority = p
	t.UpdatedAt = time.Now()
	return nil
}
func (t *Task) AddTag(tag string) error {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return ErrEmptyTag
	}
	if t.HasTag(tag) {
		return nil
	}
	t.Tags = append(t.Tags, tag)
	t.UpdatedAt = time.Now()
	return nil
}
func (t *Task) RemoveTag(tag string) error {
	tag = strings.ToLower(strings.TrimSpace(tag))
	removeIdx := -1
	for i, t := range t.Tags {
		if t == tag {
			removeIdx = i
			break
		}
	}
	if removeIdx == -1 {
		return fmt.Errorf("tag %q not found", tag)
	}
	t.Tags = slices.Delete(t.Tags, removeIdx, removeIdx+1)
	t.UpdatedAt = time.Now()
	return nil
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("status: %w", err)
	}
	candidate := Status(raw)
	if !candidate.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidStatus, raw)
	}
	*s = candidate
	return nil
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("priority: %w", err)
	}
	candidate := Priority(raw)
	if !candidate.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidPriority, raw)
	}
	*p = candidate
	return nil
}
