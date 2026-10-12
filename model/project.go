package model

import (
	"fmt"
	"strings"
	"time"
)

type Project struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Tasks []*Task `json:"tasks"`
	Timestamp
}

func (p *Project) TaskCount() int { return len(p.Tasks) }

func (p *Project) AddTask(t *Task) error {
	if t == nil || strings.TrimSpace(t.Title) == "" {
		return ErrEmptyTitle
	}
	for _, existing := range p.Tasks {
		if existing.ID == t.ID {
			return fmt.Errorf("task ID %d already exists", t.ID)
		}
	}
	p.Tasks = append(p.Tasks, t)
	p.UpdatedAt = time.Now()
	return nil
}
