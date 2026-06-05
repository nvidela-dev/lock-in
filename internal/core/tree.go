package core

import (
	"strconv"
	"strings"
	"time"
)

func setStatusRecursive(task *Task, status Status, now time.Time) {
	task.Status = status
	task.UpdatedAt = now
	for i := range task.Subtasks {
		setStatusRecursive(&task.Subtasks[i], status, now)
	}
}

func descendantCount(task *Task) int {
	count := len(task.Subtasks)
	for i := range task.Subtasks {
		count += descendantCount(&task.Subtasks[i])
	}
	return count
}

func walkVisible(items *[]TaskItem, tasks []Task, path []int, prefix string, depth int) {
	for i := range tasks {
		number := taskNumber(prefix, i)
		nextPath := appendPath(path, i)
		*items = append(*items, TaskItem{
			Number: number,
			Depth:  depth,
			Task:   &tasks[i],
			Path:   nextPath,
		})
		if !tasks[i].Collapsed {
			walkVisible(items, tasks[i].Subtasks, nextPath, number, depth+1)
		}
	}
}

func walkAll(items *[]TaskItem, tasks []Task, path []int, prefix string, depth int) {
	for i := range tasks {
		number := taskNumber(prefix, i)
		nextPath := appendPath(path, i)
		*items = append(*items, TaskItem{
			Number: number,
			Depth:  depth,
			Task:   &tasks[i],
			Path:   nextPath,
		})
		walkAll(items, tasks[i].Subtasks, nextPath, number, depth+1)
	}
}

func taskNumber(prefix string, index int) string {
	part := strconv.Itoa(index + 1)
	if prefix == "" {
		return part
	}
	return prefix + "." + part
}

func parseNumber(number string) ([]int, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return nil, ErrInvalidNumber
	}
	raw := strings.Split(number, ".")
	parts := make([]int, 0, len(raw))
	for _, part := range raw {
		if part == "" {
			return nil, ErrInvalidNumber
		}
		value, err := strconv.Atoi(part)
		if err != nil || value < 1 {
			return nil, ErrInvalidNumber
		}
		parts = append(parts, value-1)
	}
	return parts, nil
}

func taskAtPath(tasks []Task, path []int) *Task {
	if len(path) == 0 {
		return nil
	}
	idx := path[0]
	if idx < 0 || idx >= len(tasks) {
		return nil
	}
	if len(path) == 1 {
		return &tasks[idx]
	}
	return taskAtPath(tasks[idx].Subtasks, path[1:])
}

func appendPath(path []int, index int) []int {
	next := make([]int, 0, len(path)+1)
	next = append(next, path...)
	next = append(next, index)
	return next
}
