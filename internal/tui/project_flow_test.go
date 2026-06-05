package tui

import (
	"testing"

	"lock-in/internal/core"
)

func TestDeleteProjectConfirmationUsesYN(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	if err := state.CreateProject("Work", "project_work", fixedClock()); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "x")
	model = typeText(model, "n")
	model = press(model, "enter")

	if len(model.State().Projects) != 2 {
		t.Fatalf("projects after n = %#v, want no deletion", model.State().Projects)
	}
	if store.saves != 0 {
		t.Fatalf("saves after n = %d, want 0", store.saves)
	}

	model = press(model, "x")
	model = typeText(model, "y")
	model = press(model, "enter")

	if len(model.State().Projects) != 1 {
		t.Fatalf("projects after y = %#v, want one project", model.State().Projects)
	}
	if store.saves != 1 {
		t.Fatalf("saves after y = %d, want 1", store.saves)
	}
}

func TestProjectBracketSwitchPersistsActiveProject(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	if err := state.CreateProject("Work", "project_work", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if err := state.SwitchProject(1); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "]")

	modelState := model.State()
	if got := modelState.ActiveProjectIndex(); got != 1 {
		t.Fatalf("active project index = %d, want 1", got)
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
}

func TestProjectPromptSwitchPersistsActiveProject(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	if err := state.CreateProject("Work", "project_work", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if err := state.SwitchProject(1); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "g")
	model = typeText(model, "2")
	model = press(model, "enter")

	modelState := model.State()
	if got := modelState.ActiveProjectIndex(); got != 1 {
		t.Fatalf("active project index = %d, want 1", got)
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
}

func TestNumberKeyDoesNotSwitchProject(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	if err := state.CreateProject("Work", "project_work", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if err := state.SwitchProject(1); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "2")

	modelState := model.State()
	if got := modelState.ActiveProjectIndex(); got != 0 {
		t.Fatalf("active project index = %d, want 0", got)
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}
