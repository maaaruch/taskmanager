package stats

import (
	"sort"
	"taskmanager/model"
)

type Report struct {
	ByStatus   map[string]int `json:"by_status"`
	ByPriority map[string]int `json:"by_priority"`
	TopTags    []string       `json:"top_tags"`
	Total      int            `json:"total"`
}

func Build(tasks []*model.Task, topN int) Report {
	return Report{
		ByStatus:   CountByStatus(tasks),
		ByPriority: CountByPriority(tasks),
		TopTags:    TopTags(tasks, topN),
		Total:      len(tasks),
	}
}

func CountBy(tasks []*model.Task, keyFn func(*model.Task) string) map[string]int {
	counts := make(map[string]int)
	for _, t := range tasks {
		key := keyFn(t)
		counts[key]++
	}
	return counts
}

func CountByStatus(tasks []*model.Task) map[string]int {
	return CountBy(tasks, func(t *model.Task) string {
		return string(t.Status)
	})
}

func CountByPriority(tasks []*model.Task) map[string]int {
	return CountBy(tasks, func(t *model.Task) string {
		return string(t.Priority)
	})
}

func TopTags(tasks []*model.Task, n int) []string {
	if n <= 0 {
		return []string{}
	}
	counts := make(map[string]int)
	for _, t := range tasks {
		for _, tag := range t.Tags {
			counts[tag]++
		}
	}

	type tagCount struct {
		tag   string
		count int
	}
	pairs := make([]tagCount, 0, len(counts))
	for tag, count := range counts {
		pairs = append(pairs, tagCount{tag, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].tag < pairs[j].tag
	})

	if n > len(pairs) {
		n = len(pairs)
	}
	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, pairs[i].tag)
	}
	return result
}
