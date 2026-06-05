# Lock-In Architecture

This document explains how the project is organized, why the main boundaries exist, and where to make changes safely.

## High-Level Shape

Lock-In uses a small layered design:

- CLI layer: starts the app and handles non-interactive commands.
- App/UI layer: owns terminal state, input handling, rendering, prompts, and persistence timing.
- Domain layer: owns projects, tasks, task numbering, tree operations, statuses, and validation.
- Repository layer: loads and saves state as local JSON.

The app is intentionally split by responsibility instead of by framework concept alone. Bubble Tea is important, but it should not leak into the domain package. Task rules should be testable without running a terminal.

## Package Map

```text
cmd/lock-in
  main.go                 CLI entrypoint

internal/core
  types.go                Domain types, statuses, and shared errors
  state.go                App state and project lifecycle operations
  tasks.go                Task and subtask operations
  tree.go                 Task numbering, traversal, and recursive helpers
  id.go                   ID and clock helpers

internal/storage
  json_store.go           JSON repository for loading and saving state

internal/manual
  manual.go               CLI and in-app command text

internal/tui
  model.go                Bubble Tea model fields and constructor
  lifecycle.go            Init, Update, View, mutation persistence, messages
  input.go                Key dispatch and prompt dispatch
  task_commands.go        Task command handlers
  project_commands.go     Project command handlers
  selection.go            Visible task selection behavior
  confirmation.go         y/n confirmation parsing
  render.go               Screen rendering
  styles.go               Lip Gloss styles
```

Tests follow the same shape: core behavior is tested under `internal/core`, storage behavior under `internal/storage`, and UI command flows under `internal/tui`.

## Architectural Decisions

### `internal/core` Has No TUI Dependency

The core package does not import Bubble Tea, Bubbles, or Lip Gloss. It only knows about state, projects, tasks, statuses, numbering, and tree behavior.

This keeps the task rules independent from terminal rendering. For example, `Project.SetStatusCascade` can be tested directly without simulating keypresses.

### Display Numbers Are Computed, Not Stored

Task numbers like `1`, `2`, and `3.1` are display positions inside the active project. They are computed by walking the current task tree.

Stored tasks use stable IDs. Display numbers can change when tasks are added or deleted. This keeps JSON state simpler and avoids having to rewrite stored numbers after every tree change.

### Commands Mutate Through Domain Methods

The TUI does not directly edit task fields except for selection and pending prompt state. Command handlers call domain methods such as:

- `Project.AddTask`
- `Project.AddSubtask`
- `Project.SetStatus`
- `Project.SetStatusCascade`
- `Project.DeleteTask`
- `State.CreateProject`
- `State.DeleteActiveProject`

This keeps validation and tree mutation rules in one place.

### Persistence Happens After Successful Mutations

The TUI calls `finishMutation` after a command changes state. That method:

1. clears prompt state,
2. ensures selection still points at a visible task,
3. saves through the storage interface,
4. shows a status message.

Canceled prompts and validation errors do not save.

### Storage Is Behind a Small Interface

The TUI depends on this interface:

```go
type Store interface {
    Save(core.State) error
}
```

The concrete implementation is `storage.JSONStore`. Tests use an in-memory fake store. This keeps UI command tests fast and avoids writing real files during command-flow tests.

### JSON Writes Are Atomic

`JSONStore.Save` writes to a temporary file in the target directory, syncs it, closes it, and then renames it over the state file.

That avoids leaving a partially written JSON file if the process fails during save.

### Prompts Are Explicit State

Prompt flows are modeled with `promptKind` in `internal/tui/model.go`.

Single-key commands can start a prompt. When Enter is pressed, `submitPrompt` routes the current input to the right command handler based on `promptKind`.

This is why multi-step commands such as `S`, `D`, `G`, `g`, `X`, and `x` stay readable.

### Destructive Actions Use `y/n`

Task deletion and project deletion require `y/n` confirmation.

Marking a task with children as `Done` also requires confirmation because it cascades to descendants. Leaf tasks are marked done immediately.

The parser for this behavior lives in `internal/tui/confirmation.go`.

### Selection Uses Visible Items

Selection is based on the list returned by `Project.VisibleItems`. Collapsed descendants are hidden from normal movement.

Commands that target hidden tasks by number can still find them through `Project.FindTask`, which uses the full tree. `G` expands ancestors before selecting a hidden descendant.

## How To Navigate Changes

### Add Or Change A Key Command

Start in `internal/tui/input.go`.

Add the key to `handleKey`, then route it to a small command handler. Put task-specific behavior in `task_commands.go` and project-specific behavior in `project_commands.go`.

If the command needs input, add a `promptKind` in `model.go` and handle it in `submitPrompt`.

### Add A Task Rule

Start in `internal/core/tasks.go`.

If the rule needs recursive traversal or numbering, place helper functions in `internal/core/tree.go`.

Add core tests before wiring the rule into the TUI. This makes it clear whether a bug is in the domain behavior or the terminal command flow.

### Add A Project Rule

Start in `internal/core/state.go`.

Project lifecycle operations belong on `State`, because `State` owns the project list and active project.

Then call that method from `internal/tui/project_commands.go`.

### Change Persistence

Start in `internal/storage/json_store.go`.

Keep the `Store` interface small. The TUI should not know whether state is saved as JSON, SQLite, or something else.

If storage behavior changes, update `internal/storage/json_store_test.go`.

### Change Rendering

Start in `internal/tui/render.go`.

Rendering should read model state and format it. It should not mutate state. Styles belong in `internal/tui/styles.go`.

If a rendering change affects visible text, update `internal/tui/render_test.go`.

### Change The Command Manual

Update `internal/manual/manual.go`.

The same manual text is used by the CLI command and the in-app `?` view, so command docs stay in one place.

## Maintenance Rules

- Keep `internal/core` independent from terminal packages.
- Prefer adding behavior in the domain first, then calling it from the TUI.
- Keep command handlers small and named after the user action.
- Keep prompt-based workflows explicit with `promptKind`.
- Save only after successful mutations.
- Do not save after canceling, invalid input, or failed confirmation.
- Keep tests near the behavior they describe.
- Run `gofmt -w internal cmd` before committing.
- Run `go test ./...` before pushing.

## Current Validation Commands

```sh
gofmt -w internal cmd
go test ./...
go run ./cmd/lock-in man
```
