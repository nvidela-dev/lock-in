package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"lock-in/internal/core"
)

func TestLoadMissingFileInitializesState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	store := NewJSONStore(path)

	state, err := store.Load(testID, testClock)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Projects) != 1 || state.Projects[0].Name != "Inbox" {
		t.Fatalf("state projects = %#v, want default Inbox", state.Projects)
	}
}

func TestSaveAndLoadState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.json")
	store := NewJSONStore(path)
	state := core.NewState(testID, testClock)
	project, err := state.ActiveProjectRef()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := project.AddTask("Task 1", "task_1", testClock()); err != nil {
		t.Fatal(err)
	}

	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.Load(testID, testClock)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Projects[0].Tasks) != 1 || loaded.Projects[0].Tasks[0].Title != "Task 1" {
		t.Fatalf("loaded tasks = %#v, want saved task", loaded.Projects[0].Tasks)
	}
}

func TestLoadCorruptedFileReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := NewJSONStore(path).Load(testID, testClock)
	if err == nil {
		t.Fatal("expected corrupted file error")
	}
}

func TestSaveUsesAtomicReplacementPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	store := NewJSONStore(path)

	state := core.NewState(testID, testClock)
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() != "state.json" {
			t.Fatalf("unexpected temp file left behind: %s", entry.Name())
		}
	}
	if _, err := store.Load(testID, testClock); err != nil {
		t.Fatalf("saved file should load as JSON: %v", err)
	}
}

func testID(prefix string) string {
	return prefix + "_1"
}

func testClock() time.Time {
	return time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
}
