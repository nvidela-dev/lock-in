package tui

import (
	"time"

	"charm.land/bubbles/v2/textinput"

	"lock-in/internal/core"
)

type Store interface {
	Save(core.State) error
}

type Model struct {
	state  core.State
	store  Store
	nextID core.IDGenerator
	clock  core.Clock

	width  int
	height int

	selectedNumber string
	statusMessage  string
	errorMessage   string

	input         textinput.Model
	prompt        promptKind
	promptLabel   string
	pendingNumber string
	pendingStatus core.Status
	manualVisible bool
}

type promptKind int

const (
	promptNone promptKind = iota
	promptAddTask
	promptAddSubtaskSelected
	promptGotoNumber
	promptDoneNumber
	promptReadyNumber
	promptProgressNumber
	promptDoneCascade
	promptSubtaskParent
	promptSubtaskTitle
	promptGotoProject
	promptCreateProject
	promptRenameProject
	promptDeleteTask
	promptDeleteProject
)

func New(state core.State, store Store) Model {
	return NewWithDependencies(state, store, core.RandomID, core.SystemClock)
}

func NewWithDependencies(state core.State, store Store, nextID core.IDGenerator, clock core.Clock) Model {
	state.EnsureValid(nextID, clock)
	input := textinput.New()
	input.Prompt = ""
	input.SetWidth(40)
	model := Model{
		state:  state,
		store:  store,
		nextID: nextID,
		clock:  clock,
		width:  80,
		height: 24,
		input:  input,
	}
	model.ensureSelection()
	return model
}

func (m Model) State() core.State {
	return m.state
}

func (m Model) SelectedNumber() string {
	return m.selectedNumber
}

func (m Model) PromptActive() bool {
	return m.prompt != promptNone
}

func (m Model) StatusMessage() string {
	return m.statusMessage
}

func (m Model) ErrorMessage() string {
	return m.errorMessage
}

func (m Model) RenderPlain() string {
	return m.render()
}

func (m Model) now() time.Time {
	return m.clock()
}
