package app

import (
	"fmt"
	"strings"
	"taskmanager/cli"
	"taskmanager/filter"
	"taskmanager/model"
	"taskmanager/stats"
	"time"
)

func (a *App) AddProject(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("имя проекта не может быть пустым")
	}
	now := time.Now()
	p := &model.Project{
		Timestamp: model.Timestamp{CreatedAt: now, UpdatedAt: now},
		ID:        a.NextID(),
		Name:      name,
		Tasks:     []*model.Task{},
	}
	return a.change(func() error {
		a.Projects = append(a.Projects, p)
		return nil
	})
}

func (a *App) AddTask(projectID int, title string) error {
	p := a.findProject(projectID)
	if p == nil {
		return fmt.Errorf("%w: id=%d", model.ErrProjectNotFound, projectID)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return model.ErrEmptyTitle
	}
	t := &model.Task{
		ID:       a.NextID(),
		Title:    title,
		Status:   model.StatusTodo,
		Priority: model.PriorityMid,
		Timestamp: model.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return a.change(func() error { return p.AddTask(t) })
}

func (a *App) ListProjects() error {
	return cli.PrintProjects(a.Out, a.Projects)
}

func (a *App) Stats() error {
	allTasks := a.allTasks()
	r := stats.Build(allTasks, 3)
	return cli.PrintReport(a.Out, r)
}

func (a *App) FilterByStatus(s model.Status) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: %q", model.ErrInvalidStatus, s)
	}
	tasks := filter.Apply(a.allTasks(), filter.ByStatus(s))
	return cli.PrintTasks(a.Out, tasks)
}

func (a *App) save() error {
	if err := a.Storage.Save(a.Projects); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (a *App) findProject(id int) *model.Project {
	for _, p := range a.Projects {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (a *App) allTasks() []*model.Task {
	var all []*model.Task
	for _, p := range a.Projects {
		all = append(all, p.Tasks...)
	}
	return all
}

func (a *App) FilterByPriority(p model.Priority) error {
	if !p.IsValid() {
		return fmt.Errorf("%w: %q", model.ErrInvalidPriority, p)
	}
	return cli.PrintTasks(a.Out, filter.Apply(a.allTasks(), filter.ByPriority(p)))
}

func (a *App) FilterByTag(tag string) error {
	if strings.TrimSpace(tag) == "" {
		return model.ErrEmptyTag
	}
	return cli.PrintTasks(a.Out, filter.Apply(a.allTasks(), filter.ByTag(tag)))
}

func (a *App) UpdateTask(id int, field, value string) error {
	for _, p := range a.Projects {
		for _, t := range p.Tasks {
			if t.ID != id {
				continue
			}
			return a.change(func() error {
				var err error
				switch field {
				case "status":
					err = t.SetStatus(model.Status(value))
				case "priority":
					err = t.SetPriority(model.Priority(value))
				case "add-tag":
					err = t.AddTag(value)
				case "remove-tag":
					err = t.RemoveTag(value)
				default:
					return fmt.Errorf("unknown task field: %q", field)
				}
				if err == nil {
					p.UpdatedAt = time.Now()
				}
				return err
			})
		}
	}
	return fmt.Errorf("%w: id=%d", model.ErrTaskNotFound, id)
}

func (a *App) DeleteTask(id int) error {
	for _, p := range a.Projects {
		for i, t := range p.Tasks {
			if t.ID == id {
				return a.change(func() error {
					p.Tasks = append(p.Tasks[:i:i], p.Tasks[i+1:]...)
					p.UpdatedAt = time.Now()
					return nil
				})
			}
		}
	}
	return fmt.Errorf("%w: id=%d", model.ErrTaskNotFound, id)
}

func (a *App) change(update func() error) error {
	previous := make([]*model.Project, len(a.Projects))
	for i, p := range a.Projects {
		cp := *p
		cp.Tasks = make([]*model.Task, len(p.Tasks))
		for j, t := range p.Tasks {
			ct := *t
			ct.Tags = append([]string(nil), t.Tags...)
			cp.Tasks[j] = &ct
		}
		previous[i] = &cp
	}
	err := update()
	if err == nil {
		err = a.save()
	}
	if err != nil {
		a.Projects = previous
	}
	return err
}
