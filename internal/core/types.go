package core

import (
	"errors"
	"time"
)

const (
	StatusReady      Status = "ready"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

const CurrentVersion = 1

var (
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrDuplicateProject = errors.New("project name already exists")
	ErrLastProject      = errors.New("cannot delete the last project")
	ErrNoActiveProject  = errors.New("active project not found")
	ErrInvalidNumber    = errors.New("task number not found")
)

type Status string

func (s Status) Label() string {
	switch s {
	case StatusReady:
		return "Ready"
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	default:
		return string(s)
	}
}

func NormalizeStatus(status Status) Status {
	switch status {
	case StatusReady, StatusInProgress, StatusDone:
		return status
	default:
		return StatusReady
	}
}

type State struct {
	Version       int       `json:"version"`
	ActiveProject string    `json:"active_project"`
	Projects      []Project `json:"projects"`
}

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tasks     []Task    `json:"tasks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	Collapsed bool      `json:"collapsed"`
	Subtasks  []Task    `json:"subtasks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskItem struct {
	Number string
	Depth  int
	Task   *Task
	Path   []int
}

type IDGenerator func(prefix string) string
type Clock func() time.Time
