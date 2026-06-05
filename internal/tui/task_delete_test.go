package tui

import (
	"testing"

	"lock-in/internal/core"
)

func TestDeleteSelectedTaskWithX(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 2", "task_2", fixedClock()); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "X")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 2 {
		t.Fatalf("tasks after X = %#v, want no deletion before confirmation", project.Tasks)
	}
	if !model.PromptActive() {
		t.Fatal("expected delete confirmation prompt")
	}
	if store.saves != 0 {
		t.Fatalf("saves after X = %d, want 0 before confirmation", store.saves)
	}

	model = typeText(model, "y")
	model = press(model, "enter")

	modelState = model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 1 || project.Tasks[0].Title != "Task 2" {
		t.Fatalf("tasks = %#v, want only Task 2", project.Tasks)
	}
	if model.SelectedNumber() != "1" {
		t.Fatalf("selected number = %q, want renumbered 1", model.SelectedNumber())
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
}

func TestDeleteSelectedSubtaskWithX(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("1", "Subtask 1", "task_11", fixedClock()); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("1", "Subtask 2", "task_12", fixedClock()); err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)
	model = press(model, "j")

	model = press(model, "X")
	model = typeText(model, "y")
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
	if task.Title != "Subtask 2" {
		t.Fatalf("remaining subtask title = %q, want Subtask 2", task.Title)
	}
	if model.SelectedNumber() != "1.1" {
		t.Fatalf("selected number = %q, want remaining subtask 1.1", model.SelectedNumber())
	}
}

func TestDeleteLastTaskClearsSelection(t *testing.T) {
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

	model = press(model, "X")
	model = typeText(model, "y")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 0 {
		t.Fatalf("tasks = %#v, want none", project.Tasks)
	}
	if model.SelectedNumber() != "" {
		t.Fatalf("selected number = %q, want empty", model.SelectedNumber())
	}
}

func TestDeleteTaskConfirmationCanCancelWithN(t *testing.T) {
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

	model = press(model, "X")
	model = typeText(model, "n")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 1 {
		t.Fatalf("tasks = %#v, want no deletion", project.Tasks)
	}
	if model.StatusMessage() != "Delete canceled" {
		t.Fatalf("status message = %q, want Delete canceled", model.StatusMessage())
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}

func TestDeleteTaskConfirmationRejectsInvalidInput(t *testing.T) {
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

	model = press(model, "X")
	model = typeText(model, "maybe")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 1 {
		t.Fatalf("tasks = %#v, want no deletion", project.Tasks)
	}
	if model.ErrorMessage() == "" {
		t.Fatal("expected confirmation error message")
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}
