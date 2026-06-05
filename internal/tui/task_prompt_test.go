package tui

import (
	"testing"

	"lock-in/internal/core"
)

func TestPromptGotoExpandsHiddenSubtask(t *testing.T) {
	state := stateWithTasks(t)
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "G")
	model = typeText(model, "3.1")
	model = press(model, "enter")

	if model.SelectedNumber() != "3.1" {
		t.Fatalf("selected number = %q, want 3.1", model.SelectedNumber())
	}
	modelState := model.State()
	project, err := modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if numbers := itemNumbers(project.VisibleItems()); !contains(numbers, "3.1") {
		t.Fatalf("visible numbers = %v, want 3.1 after jump", numbers)
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
}

func TestPromptSubtaskFlow(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", fixedClock()); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "S")
	model = typeText(model, "1")
	model = press(model, "enter")
	model = typeText(model, "Subtask 1")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	task, _, err := project.FindTask("1.1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "Subtask 1" {
		t.Fatalf("subtask title = %q, want Subtask 1", task.Title)
	}
	if model.SelectedNumber() != "1.1" {
		t.Fatalf("selected number = %q, want 1.1", model.SelectedNumber())
	}
}

func TestPromptCancelDoesNotMutate(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "a")
	model = typeText(model, "Draft task")
	model = press(model, "esc")

	modelState := model.State()
	project, err := modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 0 {
		t.Fatalf("tasks = %#v, want none after cancel", project.Tasks)
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0 after cancel", store.saves)
	}
}

func TestInvalidInputDoesNotMutate(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", fixedClock()); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "D")
	model = typeText(model, "99")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	task, _, err := project.FindTask("1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != core.StatusReady {
		t.Fatalf("status = %q, want ready", task.Status)
	}
	if model.ErrorMessage() == "" {
		t.Fatal("expected error message for invalid task number")
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}
