package model

import "errors"

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidStatus   = errors.New("invalid status")
	ErrInvalidPriority = errors.New("invalid priority")
	ErrEmptyTitle      = errors.New("title cannot be empty")
	ErrEmptyTag        = errors.New("tag cannot be empty")
)
