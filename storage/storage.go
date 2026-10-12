package storage

import "taskmanager/model"

type Storage interface {
	Save(projects []*model.Project) error
	Load() ([]*model.Project, error)
}
