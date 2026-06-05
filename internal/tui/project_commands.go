package tui

import (
	"fmt"
	"strconv"
	"strings"
)

func (m *Model) gotoProject(number string) {
	index, err := strconv.Atoi(strings.TrimSpace(number))
	if err != nil || index < 1 {
		m.finishPromptWithError(fmt.Errorf("project number not found"))
		return
	}
	if err := m.state.SwitchProject(index); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.ensureSelection()
	m.finishMutation(fmt.Sprintf("Switched to %s", m.activeProjectName()))
}

func (m *Model) createProject(name string) {
	if err := m.state.CreateProject(name, m.nextID("project"), m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.ensureSelection()
	m.finishMutation(fmt.Sprintf("Created project %s", m.activeProjectName()))
}

func (m *Model) renameProject(name string) {
	if err := m.state.RenameActiveProject(name, m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.finishMutation(fmt.Sprintf("Renamed project to %s", m.activeProjectName()))
}

func (m *Model) deleteProject(value string) {
	confirmed, err := parseConfirmation(value)
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if !confirmed {
		m.cancelPrompt("Delete canceled")
		return
	}
	name := m.activeProjectName()
	if err := m.state.DeleteActiveProject(); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.ensureSelection()
	m.finishMutation(fmt.Sprintf("Deleted project %s", name))
}

func (m *Model) moveProject(delta int) {
	if len(m.state.Projects) == 0 {
		m.setError("No projects")
		return
	}
	current := m.state.ActiveProjectIndex()
	if current < 0 {
		m.setError("Active project not found")
		return
	}
	next := (current + delta + len(m.state.Projects)) % len(m.state.Projects)
	if err := m.state.SwitchProject(next + 1); err != nil {
		m.setError(err.Error())
		return
	}
	m.ensureSelection()
	m.finishMutation(fmt.Sprintf("Switched to %s", m.activeProjectName()))
}
