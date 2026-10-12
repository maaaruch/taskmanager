package model

import (
	"testing"
)

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"Валидный статус ToDo", StatusTodo, true},
		{"Валидный статус InProgress", StatusInProgress, true},
		{"Валидный статус Done", StatusDone, true},
		{"Невалидный статус aboba", Status("aboba"), false},
		{"Пустая строка", Status(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("Status.IsValid() для %s = %v; ожидалось %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestPriority_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		want     bool
	}{
		{"Валидный приоритет Low", PriorityLow, true},
		{"Валидный приоритет Mid", PriorityMid, true},
		{"Валидный приоритет High", PriorityHigh, true},
		{"Невалидный приоритет aboba", Priority("aboba"), false},
		{"Пустая строка", Priority(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.priority.IsValid()
			if got != tt.want {
				t.Errorf("Priority.IsValid() для %s = %v; ожидалось %v", tt.priority, got, tt.want)
			}
		})
	}
}
