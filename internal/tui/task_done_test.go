package tui

import (
	"testing"

	"lock-in/internal/core"
)

func TestDoneCascadeRequiresConfirmation(t *testing.T) {
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
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "d")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	parent, _, err := project.FindTask("1")
	if err != nil {
		t.Fatal(err)
	}
	child, _, err := project.FindTask("1.1")
	if err != nil {
		t.Fatal(err)
	}
	if parent.Status != core.StatusReady || child.Status != core.StatusReady {
		t.Fatalf("statuses after d = parent %q child %q, want unchanged", parent.Status, child.Status)
	}
	if !model.PromptActive() {
		t.Fatal("expected done cascade confirmation prompt")
	}
	if store.saves != 0 {
		t.Fatalf("saves after d = %d, want 0 before confirmation", store.saves)
	}

	model = typeText(model, "y")
	model = press(model, "enter")

	modelState = model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	for _, number := range []string{"1", "1.1"} {
		task, _, err := project.FindTask(number)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status != core.StatusDone {
			t.Fatalf("#%s status = %q, want done", number, task.Status)
		}
	}
	if store.saves != 1 {
		t.Fatalf("saves after y = %d, want 1", store.saves)
	}
}

func TestDoneCascadeCanCancelWithN(t *testing.T) {
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
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "d")
	model = typeText(model, "n")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	for _, number := range []string{"1", "1.1"} {
		task, _, err := project.FindTask(number)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status != core.StatusReady {
			t.Fatalf("#%s status = %q, want ready", number, task.Status)
		}
	}
	if model.StatusMessage() != "Done canceled" {
		t.Fatalf("status message = %q, want Done canceled", model.StatusMessage())
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}

func TestDoneLeafDoesNotPrompt(t *testing.T) {
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

	model = press(model, "d")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	task, _, err := project.FindTask("1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != core.StatusDone {
		t.Fatalf("status = %q, want done", task.Status)
	}
	if model.PromptActive() {
		t.Fatal("leaf done should not open confirmation prompt")
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
}

func TestPromptDoneCascadeUsesDFlow(t *testing.T) {
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
	store := &memoryStore{}
	model := NewWithDependencies(state, store, sequenceIDs(), fixedClock)

	model = press(model, "D")
	model = typeText(model, "1")
	model = press(model, "enter")
	if !model.PromptActive() {
		t.Fatal("expected cascade confirmation prompt after D 1")
	}
	model = typeText(model, "y")
	model = press(model, "enter")

	modelState := model.State()
	project, err = modelState.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	child, _, err := project.FindTask("1.1")
	if err != nil {
		t.Fatal(err)
	}
	if child.Status != core.StatusDone {
		t.Fatalf("child status = %q, want done", child.Status)
	}
}
