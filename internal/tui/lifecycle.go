package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = maxInt(msg.Width, 40)
		m.height = maxInt(msg.Height, 10)
		m.input.SetWidth(maxInt(m.width-20, 20))
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	default:
		return m, nil
	}
}

func (m Model) View() tea.View {
	view := tea.NewView(m.render())
	view.AltScreen = true
	return view
}

func (m *Model) finishMutation(message string) {
	m.prompt = promptNone
	m.promptLabel = ""
	m.pendingNumber = ""
	m.pendingStatus = ""
	m.input = textinput.New()
	m.ensureSelection()
	if err := m.persist(); err != nil {
		m.setError("Save failed: " + err.Error())
		return
	}
	m.setStatus(message)
}

func (m *Model) finishPromptWithError(err error) {
	m.setError(err.Error())
}

func (m *Model) persist() error {
	if m.store == nil {
		return nil
	}
	return m.store.Save(m.state)
}

func (m *Model) clearMessages() {
	m.statusMessage = ""
	m.errorMessage = ""
}

func (m *Model) setStatus(message string) {
	m.statusMessage = message
	m.errorMessage = ""
}

func (m *Model) setError(message string) {
	m.errorMessage = message
	m.statusMessage = ""
}

func (m Model) activeProjectName() string {
	project, err := m.state.ActiveProjectRef()
	if err != nil {
		return "unknown"
	}
	return project.Name
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
