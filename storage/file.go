package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"taskmanager/model"
)

type FileStorage struct {
	Path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

func (f *FileStorage) Save(projects []*model.Project) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal projects: %w", err)
	}

	tmpPath := f.Path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, f.Path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func (f *FileStorage) Load() ([]*model.Project, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*model.Project{}, nil
		}
		return nil, fmt.Errorf("read file %s: %w", f.Path, err)
	}

	var projects []*model.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("unmarshal projects: %w", err)
	}

	if projects == nil {
		projects = []*model.Project{}
	}
	return projects, nil
}

var _ Storage = (*FileStorage)(nil)
