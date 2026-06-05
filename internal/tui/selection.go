package tui

import "lock-in/internal/core"

func (m *Model) moveSelection(delta int) {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.setError(err.Error())
		return
	}
	items := project.VisibleItems()
	if len(items) == 0 {
		m.selectedNumber = ""
		m.setStatus("No tasks in this project")
		return
	}
	current := m.selectedVisibleIndex(items)
	next := current + delta
	if next < 0 {
		next = 0
	}
	if next >= len(items) {
		next = len(items) - 1
	}
	m.selectedNumber = items[next].Number
}

func (m *Model) ensureSelection() {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		m.selectedNumber = ""
		return
	}
	items := project.VisibleItems()
	if len(items) == 0 {
		m.selectedNumber = ""
		return
	}
	if m.selectedNumber != "" {
		for _, item := range items {
			if item.Number == m.selectedNumber {
				return
			}
		}
	}
	m.selectedNumber = items[0].Number
}

func (m Model) selectedVisibleIndex(items []core.TaskItem) int {
	return visibleIndexForNumber(items, m.selectedNumber)
}

func visibleIndexForNumber(items []core.TaskItem, number string) int {
	for i, item := range items {
		if item.Number == number {
			return i
		}
	}
	return 0
}

func selectionAfterDelete(items []core.TaskItem, deletedIndex int) string {
	if len(items) == 0 {
		return ""
	}
	if deletedIndex >= len(items) {
		deletedIndex = len(items) - 1
	}
	if deletedIndex < 0 {
		deletedIndex = 0
	}
	return items[deletedIndex].Number
}
