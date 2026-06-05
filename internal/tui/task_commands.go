package tui

import (
	"fmt"
	"strings"

	"lock-in/internal/core"
)

func (m *Model) addTask(title string) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	number, err := project.AddTask(title, m.nextID("task"), m.now())
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = number
	m.finishMutation(fmt.Sprintf("Added task #%s", number))
}

func (m *Model) addSubtask(parentNumber, title string) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	number, err := project.AddSubtask(parentNumber, title, m.nextID("task"), m.now())
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = number
	m.finishMutation(fmt.Sprintf("Added subtask #%s", number))
}

func (m *Model) gotoNumber(number string) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	now := m.now()
	if err := project.ExpandAncestors(number, now); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = strings.TrimSpace(number)
	m.finishMutation(fmt.Sprintf("Moved to #%s", m.selectedNumber))
}

func (m *Model) setNumberStatus(number string, status core.Status) {
	if status == core.StatusDone {
		m.setNumberDone(number)
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if err := project.SetStatus(number, status, m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = strings.TrimSpace(number)
	m.finishMutation(fmt.Sprintf("Marked #%s %s", m.selectedNumber, status.Label()))
}

func (m *Model) setNumberDone(number string) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	number = strings.TrimSpace(number)
	descendants, err := project.DescendantCount(number)
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if descendants > 0 {
		m.pendingNumber = number
		m.startPrompt(promptDoneCascade, fmt.Sprintf("Mark #%s and %d children done? y/n", number, descendants))
		return
	}
	if err := project.SetStatus(number, core.StatusDone, m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = number
	m.finishMutation(fmt.Sprintf("Marked #%s %s", m.selectedNumber, core.StatusDone.Label()))
}

func (m *Model) captureSubtaskParent(number string) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if _, _, err := project.FindTask(number); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.pendingNumber = strings.TrimSpace(number)
	m.startPrompt(promptSubtaskTitle, fmt.Sprintf("Subtask title for #%s", m.pendingNumber))
}

func (m *Model) collapseSelected() {
	if m.selectedNumber == "" {
		m.setError("No task selected")
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.setError(err.Error())
		return
	}
	if err := project.Collapse(m.selectedNumber, m.now()); err != nil {
		m.setError(err.Error())
		return
	}
	m.finishMutation(fmt.Sprintf("Collapsed #%s", m.selectedNumber))
}

func (m *Model) expandSelected() {
	if m.selectedNumber == "" {
		m.setError("No task selected")
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.setError(err.Error())
		return
	}
	if err := project.Expand(m.selectedNumber, m.now()); err != nil {
		m.setError(err.Error())
		return
	}
	m.finishMutation(fmt.Sprintf("Expanded #%s", m.selectedNumber))
}

func (m *Model) setSelectedStatus(status core.Status) {
	if m.selectedNumber == "" {
		m.setError("No task selected")
		return
	}
	if status == core.StatusDone {
		m.setNumberDone(m.selectedNumber)
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.setError(err.Error())
		return
	}
	if err := project.SetStatus(m.selectedNumber, status, m.now()); err != nil {
		m.setError(err.Error())
		return
	}
	m.finishMutation(fmt.Sprintf("Marked #%s %s", m.selectedNumber, status.Label()))
}

func (m *Model) confirmDoneCascade(value string) {
	confirmed, err := parseConfirmation(value)
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if !confirmed {
		m.cancelPrompt("Done canceled")
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if err := project.SetStatusCascade(m.pendingNumber, core.StatusDone, m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = m.pendingNumber
	m.finishMutation(fmt.Sprintf("Marked #%s and children %s", m.selectedNumber, core.StatusDone.Label()))
}

func (m *Model) startDeleteTaskPrompt() {
	if m.selectedNumber == "" {
		m.setError("No task selected")
		return
	}
	m.pendingNumber = m.selectedNumber
	m.startPrompt(promptDeleteTask, fmt.Sprintf("Delete #%s? y/n", m.pendingNumber))
}

func (m *Model) deleteTask(value string) {
	confirmed, err := parseConfirmation(value)
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	if !confirmed {
		m.cancelPrompt("Delete canceled")
		return
	}
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.finishPromptWithError(err)
		return
	}
	visibleItems := project.VisibleItems()
	currentIndex := visibleIndexForNumber(visibleItems, m.pendingNumber)
	deletedNumber := m.pendingNumber
	if err := project.DeleteTask(deletedNumber, m.now()); err != nil {
		m.finishPromptWithError(err)
		return
	}
	m.selectedNumber = selectionAfterDelete(project.VisibleItems(), currentIndex)
	m.finishMutation(fmt.Sprintf("Deleted #%s", deletedNumber))
}
