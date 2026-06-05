package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"lock-in/internal/core"
	"lock-in/internal/manual"
)

func (m Model) render() string {
	if m.manualVisible {
		return m.renderManual()
	}
	header := m.renderHeader()
	project, err := m.state.ActiveProjectRef()
	var body string
	if err != nil {
		body = errorStyle.Render(err.Error())
	} else {
		body = m.renderTaskList(*project)
	}
	footer := m.renderFooter()
	projectBar := m.renderProjectBar()

	contentHeight := maxInt(m.height-lipgloss.Height(header)-lipgloss.Height(footer)-lipgloss.Height(projectBar), 1)
	body = fitLines(body, contentHeight)
	body = lipgloss.NewStyle().Height(contentHeight).Render(body)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer, projectBar)
}

func (m Model) renderHeader() string {
	projectName := m.activeProjectName()
	return titleStyle.Width(m.width).Render(fmt.Sprintf("Lock-In | %s", projectName))
}

func (m Model) renderTaskList(project core.Project) string {
	items := project.VisibleItems()
	if len(items) == 0 {
		return footerStyle.Render("No tasks yet. Press a to add one, c to create a project, or ? for help.")
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		prefix := strings.Repeat("   ", item.Depth)
		if item.Depth > 0 {
			prefix += "|_ "
		}
		collapseMarker := " "
		if len(item.Task.Subtasks) > 0 {
			if item.Task.Collapsed {
				collapseMarker = "+"
			} else {
				collapseMarker = "-"
			}
		}
		status := statusStyleFor(item.Task.Status).Render(item.Task.Status.Label())
		line := fmt.Sprintf("%s%s # %s | %s | %s", prefix, collapseMarker, item.Number, item.Task.Title, status)
		if item.Number == m.selectedNumber {
			line = selectedStyle.Width(m.width).Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderFooter() string {
	if m.prompt != promptNone {
		if m.errorMessage != "" {
			return errorStyle.Width(m.width).Render(fmt.Sprintf("%s | %s: %s", m.errorMessage, m.promptLabel, m.input.View()))
		}
		return footerStyle.Width(m.width).Render(fmt.Sprintf("%s: %s", m.promptLabel, m.input.View()))
	}
	if m.errorMessage != "" {
		return errorStyle.Width(m.width).Render(m.errorMessage)
	}
	if m.statusMessage != "" {
		return statusStyle.Width(m.width).Render(m.statusMessage)
	}
	return footerStyle.Width(m.width).Render("a task | s subtask | j/k tasks | [/] projects | d/r/p status | X confirm-delete | ? manual")
}

func (m Model) renderProjectBar() string {
	parts := make([]string, 0, len(m.state.Projects))
	active := m.state.ActiveProject
	for i, project := range m.state.Projects {
		label := fmt.Sprintf("%d: %s", i+1, project.Name)
		if project.ID == active {
			parts = append(parts, projectActive.Render(label))
		} else {
			parts = append(parts, projectInactive.Render(label))
		}
	}
	bar := "[" + strings.Join(parts, " |") + "]"
	return lipgloss.NewStyle().Width(m.width).Render(bar)
}

func (m Model) renderManual() string {
	width := maxInt(m.width-6, 40)
	height := maxInt(m.height-4, 8)
	content := fitLines(manual.Text, height-2)
	box := manualBorderStyle.Width(width).Height(height).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func statusStyleFor(status core.Status) lipgloss.Style {
	switch status {
	case core.StatusDone:
		return doneStyle
	case core.StatusInProgress:
		return progressStyle
	default:
		return readyStyle
	}
}

func fitLines(text string, maxLines int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return text
	}
	if maxLines <= 0 {
		return ""
	}
	if maxLines == 1 {
		return "..."
	}
	out := append([]string{}, lines[:maxLines-1]...)
	out = append(out, "...")
	return strings.Join(out, "\n")
}
