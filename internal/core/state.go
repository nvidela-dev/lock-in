package core

import (
	"fmt"
	"strings"
	"time"
)

func NewState(nextID IDGenerator, clock Clock) State {
	now := clock()
	project := Project{
		ID:        nextID("project"),
		Name:      "Inbox",
		CreatedAt: now,
		UpdatedAt: now,
	}
	return State{
		Version:       CurrentVersion,
		ActiveProject: project.ID,
		Projects:      []Project{project},
	}
}

func (s *State) EnsureValid(nextID IDGenerator, clock Clock) {
	if s.Version == 0 {
		s.Version = CurrentVersion
	}
	if len(s.Projects) == 0 {
		*s = NewState(nextID, clock)
		return
	}
	activeExists := false
	for i := range s.Projects {
		project := &s.Projects[i]
		project.Name = strings.TrimSpace(project.Name)
		if project.Name == "" {
			project.Name = fmt.Sprintf("Project %d", i+1)
		}
		if project.ID == "" {
			project.ID = nextID("project")
		}
		if project.CreatedAt.IsZero() {
			project.CreatedAt = clock()
		}
		if project.UpdatedAt.IsZero() {
			project.UpdatedAt = project.CreatedAt
		}
		normalizeTasks(project.Tasks, nextID, clock)
		if project.ID == s.ActiveProject {
			activeExists = true
		}
	}
	if !activeExists {
		s.ActiveProject = s.Projects[0].ID
	}
}

func (s *State) ActiveProjectRef() (*Project, error) {
	for i := range s.Projects {
		if s.Projects[i].ID == s.ActiveProject {
			return &s.Projects[i], nil
		}
	}
	return nil, ErrNoActiveProject
}

func (s *State) ActiveProjectIndex() int {
	for i := range s.Projects {
		if s.Projects[i].ID == s.ActiveProject {
			return i
		}
	}
	return -1
}

func (s *State) SwitchProject(displayIndex int) error {
	if displayIndex < 1 || displayIndex > len(s.Projects) {
		return fmt.Errorf("project %d does not exist", displayIndex)
	}
	s.ActiveProject = s.Projects[displayIndex-1].ID
	return nil
}

func (s *State) CreateProject(name, id string, now time.Time) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrEmptyName
	}
	if s.projectNameExists(name, "") {
		return ErrDuplicateProject
	}
	project := Project{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.Projects = append(s.Projects, project)
	s.ActiveProject = id
	return nil
}

func (s *State) RenameActiveProject(name string, now time.Time) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrEmptyName
	}
	project, err := s.ActiveProjectRef()
	if err != nil {
		return err
	}
	if s.projectNameExists(name, project.ID) {
		return ErrDuplicateProject
	}
	project.Name = name
	project.UpdatedAt = now
	return nil
}

func (s *State) DeleteActiveProject() error {
	if len(s.Projects) <= 1 {
		return ErrLastProject
	}
	idx := s.ActiveProjectIndex()
	if idx < 0 {
		return ErrNoActiveProject
	}
	s.Projects = append(s.Projects[:idx], s.Projects[idx+1:]...)
	if idx >= len(s.Projects) {
		idx = len(s.Projects) - 1
	}
	s.ActiveProject = s.Projects[idx].ID
	return nil
}

func (s *State) projectNameExists(name, exceptID string) bool {
	for _, project := range s.Projects {
		if project.ID == exceptID {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(project.Name), name) {
			return true
		}
	}
	return false
}

func normalizeTasks(tasks []Task, nextID IDGenerator, clock Clock) {
	for i := range tasks {
		task := &tasks[i]
		if task.ID == "" {
			task.ID = nextID("task")
		}
		task.Title = strings.TrimSpace(task.Title)
		if task.Title == "" {
			task.Title = "Untitled"
		}
		task.Status = NormalizeStatus(task.Status)
		if task.CreatedAt.IsZero() {
			task.CreatedAt = clock()
		}
		if task.UpdatedAt.IsZero() {
			task.UpdatedAt = task.CreatedAt
		}
		normalizeTasks(task.Subtasks, nextID, clock)
	}
}
