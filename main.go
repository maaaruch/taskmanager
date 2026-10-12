package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"taskmanager/app"
	"taskmanager/model"
	"taskmanager/storage"
)

func main() {
	fileStorage := storage.NewFileStorage("tasks.json")
	a, err := app.New(fileStorage, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка инициализации:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "Task Manager. Введите 'help' для списка команд.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		if err := dispatch(a, cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, "ошибка:", err)
		}

		if cmd == "exit" || cmd == "quit" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка чтения:", err)
	}
}

func dispatch(a *app.App, cmd string, args []string) error {
	switch cmd {
	case "help":
		printHelp(a.Out)
	case "list":
		return a.ListProjects()
	case "stats":
		return a.Stats()
	case "add-project":
		if len(args) < 1 {
			return fmt.Errorf("usage: add-project <name>")
		}
		return a.AddProject(strings.Join(args, " "))
	case "add-task":
		if len(args) < 2 {
			return fmt.Errorf("usage: add-task <projectID> <title>")
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid project ID: %q", args[0])
		}
		return a.AddTask(pid, strings.Join(args[1:], " "))
	case "status", "priority", "add-tag", "remove-tag":
		if len(args) < 2 {
			return fmt.Errorf("usage: %s <taskID> <value>", cmd)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %q", args[0])
		}
		return a.UpdateTask(id, cmd, strings.Join(args[1:], " "))
	case "delete-task":
		if len(args) != 1 {
			return fmt.Errorf("usage: delete-task <taskID>")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %q", args[0])
		}
		return a.DeleteTask(id)
	case "filter":
		if len(args) < 2 {
			return fmt.Errorf("usage: filter <field> <value>")
		}
		field, value := args[0], strings.Join(args[1:], " ")
		switch field {
		case "status":
			return a.FilterByStatus(model.Status(value))
		case "priority":
			return a.FilterByPriority(model.Priority(value))
		case "tag":
			return a.FilterByTag(value)
		default:
			return fmt.Errorf("unknown filter field: %q", field)
		}
	case "exit", "quit":
		fmt.Fprintln(a.Out, "Пока!")
		return nil
	default:
		return fmt.Errorf("unknown command: %q", cmd)
	}
	return nil
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, `Доступные команды:
  list                              — список всех проектов
  add-project <name>                — создать проект
  add-task <projectID> <title>      — добавить задачу
  filter status <value>             — фильтр по статусу (todo|in_progress|done)
  filter priority <value>           — фильтр по приоритету (low|medium|high)
  filter tag <tag>                  — фильтр по тегу
  status <taskID> <value>            — статус (todo|in_progress|done)
  priority <taskID> <value>          — приоритет (low|medium|high)
  add-tag <taskID> <tag>             — добавить тег
  remove-tag <taskID> <tag>          — удалить тег
  delete-task <taskID>               — удалить задачу
  stats                             — статистика
  exit                              — выход`)
}
