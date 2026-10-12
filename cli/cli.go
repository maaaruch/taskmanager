package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"taskmanager/model"
	"taskmanager/stats"
)

func PrintProjects(w io.Writer, projects []*model.Project) error {
	if len(projects) == 0 {
		_, err := fmt.Fprintln(w, "Проектов нет.")
		return err
	}
	for _, p := range projects {
		if err := PrintProject(w, p); err != nil {
			return err
		}
	}
	return nil
}

func PrintProject(w io.Writer, p *model.Project) error {
	if _, err := fmt.Fprintf(w, "=== [%d] %s (%d задач) ===\n", p.ID, p.Name, p.TaskCount()); err != nil {
		return err
	}
	for _, t := range p.Tasks {
		if _, err := fmt.Fprintln(w, t); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w)
	return err
}

func PrintTasks(w io.Writer, tasks []*model.Task) error {
	if len(tasks) == 0 {
		_, err := fmt.Fprintln(w, "Задач нет.")
		return err
	}
	for _, t := range tasks {
		if _, err := fmt.Fprintln(w, t); err != nil {
			return err
		}
	}
	return nil
}

func PrintReport(w io.Writer, r stats.Report) error {
	var b strings.Builder

	fmt.Fprintf(&b, "Всего задач: %d\n\n", r.Total)

	b.WriteString("По статусам:\n")
	writeCounts(&b, r.ByStatus)

	b.WriteString("\nПо приоритетам:\n")
	writeCounts(&b, r.ByPriority)

	b.WriteString("\nТоп тегов:\n")
	if len(r.TopTags) == 0 {
		b.WriteString("  (нет)\n")
	} else {
		for i, tag := range r.TopTags {
			fmt.Fprintf(&b, "  %d. %s\n", i+1, tag)
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func writeCounts(b *strings.Builder, counts map[string]int) {
	if len(counts) == 0 {
		b.WriteString("  (нет)\n")
		return
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(b, "  %s: %d\n", k, counts[k])
	}
}
