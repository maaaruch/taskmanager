# Task Manager

Консольный менеджер проектов и задач на Go, без сторонних зависимостей.

## Запуск

```sh
cd ~/Desktop/taskmanager
go run .
```

Или запустите собранную программу: `./taskmanager`.
Данные автоматически сохраняются в `tasks.json` в текущей рабочей папке.
Запускайте программу из одной и той же папки, чтобы открывать прежние задачи.

## Пример

```text
add-project Учеба
add-task 1 Разобраться с Go
list
status 2 in_progress
priority 2 high
add-tag 2 учеба
filter tag учеба
status 2 done
stats
exit
```

ID проектов и задач выдаются автоматически. Используйте ID из команды `list`.
Названия и теги могут содержать пробелы, кавычки не нужны.
Все команды: `help`. Удалить задачу: `delete-task <taskID>`.
Удалить тег: `remove-tag <taskID> <tag>`.

## Проверка и сборка

```sh
go test ./...
go vet ./...
go build -o taskmanager .
```
