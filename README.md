# Lock-In

Lock-In is a local terminal task manager written in Go. It provides a keyboard-driven TUI for managing projects, tasks, nested subtasks, and task status.

## Run

Start the TUI from the repository root:

```sh
go run ./cmd/lock-in
```

Print the command manual:

```sh
go run ./cmd/lock-in man
```

Print CLI help:

```sh
go run ./cmd/lock-in --help
```

The app stores data at `os.UserConfigDir()/lock-in/state.json` by default. Use `LOCK_IN_DATA` to point at a custom state file:

```sh
LOCK_IN_DATA=/tmp/lock-in.json go run ./cmd/lock-in
```

If the local zsh helper is installed, the app can also be started with:

```sh
lockin
```

## Packages

Lock-In uses:

- `charm.land/bubbletea/v2`: TUI application loop, messages, keyboard input, and rendering lifecycle.
- `charm.land/bubbles/v2/textinput`: prompt input for commands that need text or numbers.
- `charm.land/lipgloss/v2`: terminal styling and layout.
- Go standard library packages for JSON persistence, paths, timestamps, and IDs.

## Feature Overview

- Multiple projects with a tmux-like project bar.
- Hierarchical tasks with nested subtasks.
- Display numbers such as `1`, `2`, and `3.1` for task navigation.
- Statuses: `Ready`, `In Progress`, and `Done`.
- Vim-style task movement with `j` and `k`.
- Project movement with `[` and `]`.
- Prompted jumps with `G` for tasks and `g` for projects.
- Expand/collapse task children with `h` and `l`.
- Done cascade confirmation for tasks with children.
- `y/n` confirmation for task and project deletion.
- Local JSON persistence after each successful mutation.
- In-app manual with `?`.
