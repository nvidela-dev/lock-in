package core

import (
	"errors"
	"testing"
	"time"
)

func TestTaskNumberingVisibilityAndAncestorExpansion(t *testing.T) {
	state := testState()
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	now := testTime()

	if _, err := project.AddTask("Task 1", "task_1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 2", "task_2", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 3", "task_3", now); err != nil {
		t.Fatal(err)
	}
	if number, err := project.AddSubtask("3", "Subtask 1", "task_31", now); err != nil {
		t.Fatal(err)
	} else if number != "3.1" {
		t.Fatalf("subtask number = %q, want 3.1", number)
	}
	if err := project.Collapse("3", now); err != nil {
		t.Fatal(err)
	}

	if numbers := itemNumbers(project.VisibleItems()); contains(numbers, "3.1") {
		t.Fatalf("collapsed visible numbers include 3.1: %v", numbers)
	}
	if numbers := itemNumbers(project.AllItems()); !contains(numbers, "3.1") {
		t.Fatalf("all numbers missing 3.1: %v", numbers)
	}

	if err := project.ExpandAncestors("3.1", now); err != nil {
		t.Fatal(err)
	}
	if numbers := itemNumbers(project.VisibleItems()); !contains(numbers, "3.1") {
		t.Fatalf("expanded visible numbers missing 3.1: %v", numbers)
	}
}

func TestStatusChangesByNumber(t *testing.T) {
	state := testState()
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	now := testTime()
	if _, err := project.AddTask("Task", "task_1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("1", "Subtask", "task_11", now); err != nil {
		t.Fatal(err)
	}

	if err := project.SetStatus("1.1", StatusInProgress, now); err != nil {
		t.Fatal(err)
	}
	task, _, err := project.FindTask("1.1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != StatusInProgress {
		t.Fatalf("status = %q, want %q", task.Status, StatusInProgress)
	}

	if err := project.SetStatus("9", StatusDone, now); !errors.Is(err, ErrInvalidNumber) {
		t.Fatalf("invalid status error = %v, want ErrInvalidNumber", err)
	}
}

func TestSetStatusCascadeMarksDescendants(t *testing.T) {
	state := testState()
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	now := testTime()
	if _, err := project.AddTask("Task", "task_1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("1", "Subtask", "task_11", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("1.1", "Nested", "task_111", now); err != nil {
		t.Fatal(err)
	}

	descendants, err := project.DescendantCount("1")
	if err != nil {
		t.Fatal(err)
	}
	if descendants != 2 {
		t.Fatalf("descendants = %d, want 2", descendants)
	}
	if err := project.SetStatusCascade("1", StatusDone, now); err != nil {
		t.Fatal(err)
	}
	for _, number := range []string{"1", "1.1", "1.1.1"} {
		task, _, err := project.FindTask(number)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status != StatusDone {
			t.Fatalf("#%s status = %q, want done", number, task.Status)
		}
	}
}

func TestDeleteTaskAndSubtaskByNumber(t *testing.T) {
	state := testState()
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	now := testTime()
	if _, err := project.AddTask("Task 1", "task_1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 2", "task_2", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("2", "Subtask 1", "task_21", now); err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddSubtask("2", "Subtask 2", "task_22", now); err != nil {
		t.Fatal(err)
	}

	if err := project.DeleteTask("2.1", now); err != nil {
		t.Fatal(err)
	}
	if _, _, err := project.FindTask("2.2"); !errors.Is(err, ErrInvalidNumber) {
		t.Fatalf("old 2.2 lookup error = %v, want ErrInvalidNumber after renumbering", err)
	}
	task, _, err := project.FindTask("2.1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "Subtask 2" {
		t.Fatalf("remaining subtask title = %q, want Subtask 2", task.Title)
	}

	if err := project.DeleteTask("1", now); err != nil {
		t.Fatal(err)
	}
	task, _, err = project.FindTask("1")
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "Task 2" {
		t.Fatalf("remaining root title = %q, want Task 2", task.Title)
	}
	if err := project.DeleteTask("9", now); !errors.Is(err, ErrInvalidNumber) {
		t.Fatalf("invalid delete error = %v, want ErrInvalidNumber", err)
	}
}

func TestProjectManagement(t *testing.T) {
	state := testState()
	now := testTime()

	if err := state.CreateProject("Work", "project_2", now); err != nil {
		t.Fatal(err)
	}
	if got := state.ActiveProjectIndex(); got != 1 {
		t.Fatalf("active project index = %d, want 1", got)
	}
	if err := state.CreateProject("work", "project_3", now); !errors.Is(err, ErrDuplicateProject) {
		t.Fatalf("duplicate create error = %v, want ErrDuplicateProject", err)
	}
	if err := state.RenameActiveProject("Inbox", now); !errors.Is(err, ErrDuplicateProject) {
		t.Fatalf("duplicate rename error = %v, want ErrDuplicateProject", err)
	}
	if err := state.RenameActiveProject("Client", now); err != nil {
		t.Fatal(err)
	}
	if err := state.SwitchProject(1); err != nil {
		t.Fatal(err)
	}
	if err := state.DeleteActiveProject(); err != nil {
		t.Fatal(err)
	}
	if len(state.Projects) != 1 || state.Projects[0].Name != "Client" {
		t.Fatalf("projects after delete = %#v, want only Client", state.Projects)
	}
	if err := state.DeleteActiveProject(); !errors.Is(err, ErrLastProject) {
		t.Fatalf("last delete error = %v, want ErrLastProject", err)
	}
}

func testState() State {
	return State{
		Version:       CurrentVersion,
		ActiveProject: "project_1",
		Projects: []Project{{
			ID:        "project_1",
			Name:      "Inbox",
			CreatedAt: testTime(),
			UpdatedAt: testTime(),
		}},
	}
}

func testTime() time.Time {
	return time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
}

func itemNumbers(items []TaskItem) []string {
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
