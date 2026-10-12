package filter

import "taskmanager/model"

type Filter func(t *model.Task) bool

func ByStatus(s model.Status) Filter {
	return func(t *model.Task) bool {
		return t.Status == s
	}
}

func ByPriority(p model.Priority) Filter {
	return func(t *model.Task) bool {
		return t.Priority == p
	}
}

func ByTag(tag string) Filter {
	return func(t *model.Task) bool {
		return t.HasTag(tag)
	}
}

func Apply(tasks []*model.Task, f Filter) []*model.Task {
	result := make([]*model.Task, 0, len(tasks))
	for _, t := range tasks {
		if f(t) {
			result = append(result, t)
		}
	}
	return result
}
