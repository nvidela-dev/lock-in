package core

import (
	"strconv"
	"strings"
	"time"
)

func (p *Project) AddTask(title, id string, now time.Time) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", ErrEmptyName
	}
	task := Task{
		ID:        id,
		Title:     title,
		Status:    StatusReady,
		CreatedAt: now,
		UpdatedAt: now,
	}
	p.Tasks = append(p.Tasks, task)
	p.UpdatedAt = now
	return strconv.Itoa(len(p.Tasks)), nil
}

func (p *Project) AddSubtask(parentNumber, title, id string, now time.Time) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", ErrEmptyName
	}
	parent, path, err := p.FindTask(parentNumber)
	if err != nil {
		return "", err
	}
	task := Task{
		ID:        id,
		Title:     title,
		Status:    StatusReady,
		CreatedAt: now,
		UpdatedAt: now,
	}
	parent.Subtasks = append(parent.Subtasks, task)
	parent.Collapsed = false
	parent.UpdatedAt = now
	p.touchPath(path[:len(path)-1], now)
	p.UpdatedAt = now
	return parentNumber + "." + strconv.Itoa(len(parent.Subtasks)), nil
}

func (p *Project) SetStatus(number string, status Status, now time.Time) error {
	task, path, err := p.FindTask(number)
	if err != nil {
		return err
	}
	task.Status = NormalizeStatus(status)
	task.UpdatedAt = now
	p.touchPath(path[:len(path)-1], now)
	p.UpdatedAt = now
	return nil
}

func (p *Project) SetStatusCascade(number string, status Status, now time.Time) error {
	task, path, err := p.FindTask(number)
	if err != nil {
		return err
	}
	setStatusRecursive(task, NormalizeStatus(status), now)
	p.touchPath(path[:len(path)-1], now)
	p.UpdatedAt = now
	return nil
}

func (p *Project) DescendantCount(number string) (int, error) {
	task, _, err := p.FindTask(number)
	if err != nil {
		return 0, err
	}
	return descendantCount(task), nil
}

func (p *Project) DeleteTask(number string, now time.Time) error {
	path, err := parseNumber(number)
	if err != nil {
		return err
	}
	if len(path) == 1 {
		idx := path[0]
		if idx < 0 || idx >= len(p.Tasks) {
			return ErrInvalidNumber
		}
		p.Tasks = append(p.Tasks[:idx], p.Tasks[idx+1:]...)
		p.UpdatedAt = now
		return nil
	}

	parentPath := path[:len(path)-1]
	parent := p.taskAtPath(parentPath)
	if parent == nil {
		return ErrInvalidNumber
	}
	idx := path[len(path)-1]
	if idx < 0 || idx >= len(parent.Subtasks) {
		return ErrInvalidNumber
	}
	parent.Subtasks = append(parent.Subtasks[:idx], parent.Subtasks[idx+1:]...)
	parent.UpdatedAt = now
	p.touchPath(parentPath[:len(parentPath)-1], now)
	p.UpdatedAt = now
	return nil
}

func (p *Project) Collapse(number string, now time.Time) error {
	task, path, err := p.FindTask(number)
	if err != nil {
		return err
	}
	task.Collapsed = true
	task.UpdatedAt = now
	p.touchPath(path[:len(path)-1], now)
	p.UpdatedAt = now
	return nil
}

func (p *Project) Expand(number string, now time.Time) error {
	task, path, err := p.FindTask(number)
	if err != nil {
		return err
	}
	task.Collapsed = false
	task.UpdatedAt = now
	p.touchPath(path[:len(path)-1], now)
	p.UpdatedAt = now
	return nil
}

func (p *Project) ExpandAncestors(number string, now time.Time) error {
	_, path, err := p.FindTask(number)
	if err != nil {
		return err
	}
	for depth := 1; depth < len(path); depth++ {
		task := p.taskAtPath(path[:depth])
		if task != nil {
			task.Collapsed = false
			task.UpdatedAt = now
		}
	}
	p.UpdatedAt = now
	return nil
}

func (p *Project) FindTask(number string) (*Task, []int, error) {
	parts, err := parseNumber(number)
	if err != nil {
		return nil, nil, err
	}
	task := taskAtPath(p.Tasks, parts)
	if task == nil {
		return nil, nil, ErrInvalidNumber
	}
	return task, parts, nil
}

func (p *Project) VisibleItems() []TaskItem {
	var items []TaskItem
	walkVisible(&items, p.Tasks, nil, "", 0)
	return items
}

func (p *Project) AllItems() []TaskItem {
	var items []TaskItem
	walkAll(&items, p.Tasks, nil, "", 0)
	return items
}

func (p *Project) taskAtPath(path []int) *Task {
	return taskAtPath(p.Tasks, path)
}

func (p *Project) touchPath(path []int, now time.Time) {
	for depth := 1; depth <= len(path); depth++ {
		task := p.taskAtPath(path[:depth])
		if task != nil {
			task.UpdatedAt = now
		}
	}
}
