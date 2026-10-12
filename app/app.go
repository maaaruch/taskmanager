package app

import (
	"fmt"
	"io"
	"strings"
	"taskmanager/model"
	"taskmanager/storage"
)

type App struct {
	Storage  storage.Storage
	Projects []*model.Project
	Out      io.Writer
	nextID   int
}

func New(s storage.Storage, out io.Writer) (*App, error) {
	projects, err := s.Load()
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	ids := make(map[int]bool)
	for _, p := range projects {
		if p == nil || p.ID <= 0 || ids[p.ID] || strings.TrimSpace(p.Name) == "" {
			return nil, fmt.Errorf("load: invalid project or duplicate ID")
		}
		ids[p.ID] = true
		for _, t := range p.Tasks {
			if t == nil || t.ID <= 0 || ids[t.ID] || strings.TrimSpace(t.Title) == "" || !t.Status.IsValid() || !t.Priority.IsValid() {
				return nil, fmt.Errorf("load: invalid task or duplicate ID in project %d", p.ID)
			}
			ids[t.ID] = true
		}
	}
	a := &App{
		Storage:  s,
		Projects: projects,
		Out:      out,
	}
	a.nextID = a.computeNextID()
	return a, nil
}

func (a *App) computeNextID() int {
	maxID := 0
	for _, p := range a.Projects {
		if p.ID > maxID {
			maxID = p.ID
		}
		for _, t := range p.Tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
	}
	return maxID + 1
}

func (a *App) NextID() int {
	id := a.nextID
	a.nextID++
	return id
}
