package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"lock-in/internal/core"
	"lock-in/internal/manual"
	"lock-in/internal/storage"
	"lock-in/internal/tui"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "man":
			fmt.Print(manual.Text)
			return nil
		case "-h", "--help", "help":
			fmt.Print(manual.Usage)
			return nil
		default:
			return fmt.Errorf("unknown command %q\n\n%s", args[0], manual.Usage)
		}
	}

	path, err := storage.DefaultPath()
	if err != nil {
		return err
	}
	store := storage.NewJSONStore(path)
	state, err := store.Load(core.RandomID, core.SystemClock)
	if err != nil {
		return err
	}
	model := tui.New(state, store)
	program := tea.NewProgram(model)
	_, err = program.Run()
	return err
}
