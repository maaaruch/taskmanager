package main

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"taskmanager/app"
	"taskmanager/model"
	"taskmanager/storage"
	"testing"
)

func TestCommandsPersistAcrossRestart(t *testing.T) {
	s := storage.NewFileStorage(filepath.Join(t.TempDir(), "tasks.json"))
	var out bytes.Buffer
	a, err := app.New(s, &out)
	if err != nil {
		t.Fatal(err)
	}
	commands := []string{
		"add-project Учеба", "add-task 1 Выучить Go",
		"status 2 done", "priority 2 high", "add-tag 2 Важный Тег",
	}
	for _, command := range commands {
		parts := strings.Fields(command)
		if err := dispatch(a, parts[0], parts[1:]); err != nil {
			t.Fatalf("%s: %v", command, err)
		}
	}
	a, err = app.New(s, &out)
	if err != nil {
		t.Fatal(err)
	}
	task := a.Projects[0].Tasks[0]
	if task.Title != "Выучить Go" || task.Status != model.StatusDone || task.Priority != model.PriorityHigh || !task.HasTag("важный тег") {
		t.Fatalf("unexpected task: %+v", task)
	}
	for _, args := range [][]string{{"status", "done"}, {"priority", "high"}, {"tag", "Важный", "Тег"}} {
		out.Reset()
		if err := dispatch(a, "filter", args); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out.String(), "Выучить Go") {
			t.Fatalf("filter %v: %s", args, out.String())
		}
	}
	if err := dispatch(a, "status", []string{"2", "invalid"}); !errors.Is(err, model.ErrInvalidStatus) {
		t.Fatalf("invalid status: %v", err)
	}
	if a.Projects[0].Tasks[0].Status != model.StatusDone {
		t.Fatal("invalid update changed status")
	}
	if err := dispatch(a, "remove-tag", []string{"2", "ВАЖНЫЙ", "ТЕГ"}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(a, "delete-task", []string{"2"}); err != nil {
		t.Fatal(err)
	}
	a, err = app.New(s, &out)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Projects[0].Tasks) != 0 {
		t.Fatal("deletion was not persisted")
	}
	if err := dispatch(a, "add-task", []string{"1", "Новая"}); err != nil {
		t.Fatal(err)
	}
}

type failingStorage struct{ projects []*model.Project }

func (s failingStorage) Load() ([]*model.Project, error) { return s.projects, nil }
func (s failingStorage) Save([]*model.Project) error     { return errors.New("disk unavailable") }

func TestFailedSaveRestoresState(t *testing.T) {
	a, err := app.New(failingStorage{[]*model.Project{{ID: 1, Name: "Test", Tasks: []*model.Task{{ID: 2, Title: "Task", Status: model.StatusTodo, Priority: model.PriorityMid}}}}}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.UpdateTask(2, "status", "done"); err == nil {
		t.Fatal("expected save failure")
	}
	if a.Projects[0].Tasks[0].Status != model.StatusTodo {
		t.Fatal("failed save left a changed task")
	}
	if err := a.DeleteTask(2); err == nil {
		t.Fatal("expected save failure")
	}
	if len(a.Projects[0].Tasks) != 1 {
		t.Fatal("failed save removed task")
	}
	if err := a.AddProject("New"); err == nil {
		t.Fatal("expected save failure")
	}
	if len(a.Projects) != 1 {
		t.Fatal("failed save added project")
	}
}

func TestRejectNullRecords(t *testing.T) {
	for _, projects := range [][]*model.Project{{nil}, {{ID: 1, Name: "Test", Tasks: []*model.Task{nil}}}} {
		if _, err := app.New(failingStorage{projects}, io.Discard); err == nil {
			t.Fatal("expected malformed data error")
		}
	}
}
