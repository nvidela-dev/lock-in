package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"lock-in/internal/core"
)

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if m.manualVisible {
		switch key {
		case "esc", "q", "?":
			m.manualVisible = false
			return m, nil
		default:
			return m, nil
		}
	}
	if m.prompt != promptNone {
		return m.handlePromptKey(msg)
	}

	m.clearMessages()
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "?":
		m.manualVisible = true
	case "a":
		m.startPrompt(promptAddTask, "Task title")
	case "s":
		if m.selectedNumber == "" {
			m.setError("No task selected")
			break
		}
		m.pendingNumber = m.selectedNumber
		m.startPrompt(promptAddSubtaskSelected, fmt.Sprintf("Subtask title for #%s", m.selectedNumber))
	case "G":
		m.startPrompt(promptGotoNumber, "Go to task number")
	case "g":
		m.startPrompt(promptGotoProject, "Go to project number")
	case "D":
		m.startStatusPrompt(promptDoneNumber, core.StatusDone, "Mark done: task number")
	case "R":
		m.startStatusPrompt(promptReadyNumber, core.StatusReady, "Mark ready: task number")
	case "P":
		m.startStatusPrompt(promptProgressNumber, core.StatusInProgress, "Mark in progress: task number")
	case "S":
		m.startPrompt(promptSubtaskParent, "Parent task number")
	case "j", "down":
		m.moveSelection(1)
	case "k", "up":
		m.moveSelection(-1)
	case "[":
		m.moveProject(-1)
	case "]":
		m.moveProject(1)
	case "h", "left":
		m.collapseSelected()
	case "l", "right":
		m.expandSelected()
	case "d":
		m.setSelectedStatus(core.StatusDone)
	case "r":
		m.setSelectedStatus(core.StatusReady)
	case "p":
		m.setSelectedStatus(core.StatusInProgress)
	case "X":
		m.startDeleteTaskPrompt()
	case "c":
		m.startPrompt(promptCreateProject, "Project name")
	case "e":
		m.startPrompt(promptRenameProject, "New project name")
	case "x":
		m.startPrompt(promptDeleteProject, fmt.Sprintf("Delete project %s? y/n", m.activeProjectName()))
	default:
		m.setStatus("Press ? for commands")
	}
	return m, nil
}

func (m Model) handlePromptKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.cancelPrompt("Canceled")
		return m, nil
	case "enter":
		m.submitPrompt()
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m *Model) startStatusPrompt(kind promptKind, status core.Status, label string) {
	m.pendingStatus = status
	m.startPrompt(kind, label)
}

func (m *Model) startPrompt(kind promptKind, label string) {
	m.clearMessages()
	m.prompt = kind
	m.promptLabel = label
	m.input = textinput.New()
	m.input.Prompt = ""
	m.input.SetWidth(maxInt(m.width-20, 20))
	_ = m.input.Focus()
}

func (m *Model) cancelPrompt(message string) {
	m.prompt = promptNone
	m.promptLabel = ""
	m.pendingNumber = ""
	m.pendingStatus = ""
	m.input = textinput.New()
	m.setStatus(message)
}

func (m *Model) submitPrompt() {
	value := strings.TrimSpace(m.input.Value())
	switch m.prompt {
	case promptAddTask:
		m.addTask(value)
	case promptAddSubtaskSelected:
		m.addSubtask(m.pendingNumber, value)
	case promptGotoNumber:
		m.gotoNumber(value)
	case promptGotoProject:
		m.gotoProject(value)
	case promptDoneNumber, promptReadyNumber, promptProgressNumber:
		m.setNumberStatus(value, m.pendingStatus)
	case promptDoneCascade:
		m.confirmDoneCascade(value)
	case promptSubtaskParent:
		m.captureSubtaskParent(value)
	case promptSubtaskTitle:
		m.addSubtask(m.pendingNumber, value)
	case promptCreateProject:
		m.createProject(value)
	case promptRenameProject:
		m.renameProject(value)
	case promptDeleteTask:
		m.deleteTask(value)
	case promptDeleteProject:
		m.deleteProject(value)
	}
}
