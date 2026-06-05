package tui

import (
	"regexp"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"lock-in/internal/core"
)

type memoryStore struct {
	saves int
	state core.State
}

func (s *memoryStore) Save(state core.State) error {
	s.saves++
	s.state = state
	return nil
}

func stateWithTasks(t *testing.T) core.State {
	t.Helper()
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{"Task 1", "Task 2", "Task 3"} {
		if _, err := project.AddTask(title, title, fixedClock()); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := project.AddSubtask("3", "Subtask 1", "task_31", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if err := project.Collapse("3", fixedClock()); err != nil {
		t.Fatal(err)
	}
	return state
}

func press(model Model, key string) Model {
	updated, _ := model.Update(keyMsg(key))
	return updated.(Model)
}

func typeText(model Model, text string) Model {
	for _, r := range text {
		model = press(model, string(r))
	}
	return model
}

func resize(model Model, width, height int) Model {
	updated, _ := model.Update(tea.WindowSizeMsg{Width: width, Height: height})
	return updated.(Model)
}

func keyMsg(key string) tea.KeyPressMsg {
	switch key {
	case "enter":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	case "esc":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc})
	case "down":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyDown})
	case "up":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyUp})
	case "left":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyLeft})
	case "right":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyRight})
	default:
		r := []rune(key)[0]
		return tea.KeyPressMsg(tea.Key{Code: r, Text: key})
	}
}

func sequenceIDs() core.IDGenerator {
	i := 0
	return func(prefix string) string {
		i++
		return prefix + "_" + string(rune('a'+i))
	}
}

func fixedClock() time.Time {
	return time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
}

func itemNumbers(items []core.TaskItem) []string {
	numbers := make([]string, 0, len(items))
	for _, item := range items {
		numbers = append(numbers, item.Number)
	}
	return numbers
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}
