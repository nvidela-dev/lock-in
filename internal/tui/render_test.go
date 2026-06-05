package tui

import (
	"strings"
	"testing"

	"lock-in/internal/core"
)

func TestRenderContainsTaskFooterProjectBarAndManual(t *testing.T) {
	state := core.NewState(sequenceIDs(), fixedClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", fixedClock()); err != nil {
		t.Fatal(err)
	}
	model := NewWithDependencies(state, nil, sequenceIDs(), fixedClock)
	model = resize(model, 80, 20)

	view := stripANSI(model.RenderPlain())
	for _, want := range []string{"Lock-In | Inbox", "# 1 | Task 1 | Ready", "1: Inbox"} {
		if !strings.Contains(view, want) {
			t.Fatalf("render missing %q:\n%s", want, view)
		}
	}

	model = press(model, "?")
	if view := stripANSI(model.RenderPlain()); !strings.Contains(view, "Lock-In Command Manual") {
		t.Fatalf("manual render missing title:\n%s", view)
	}
}
