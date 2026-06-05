package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"lock-in/internal/core"
)

type JSONStore struct {
	path string
}

func NewJSONStore(path string) JSONStore {
	return JSONStore{path: path}
}

func DefaultPath() (string, error) {
	if override := os.Getenv("LOCK_IN_DATA"); override != "" {
		return override, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "lock-in", "state.json"), nil
}

func (s JSONStore) Path() string {
	return s.path
}

func (s JSONStore) Load(nextID core.IDGenerator, clock core.Clock) (core.State, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core.NewState(nextID, clock), nil
		}
		return core.State{}, err
	}
	var state core.State
	if err := json.Unmarshal(data, &state); err != nil {
		return core.State{}, fmt.Errorf("load state: %w", err)
	}
	state.EnsureValid(nextID, clock)
	return state, nil
}

func (s JSONStore) Save(state core.State) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".state-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, s.path); err != nil {
		return err
	}
	cleanup = false
	return nil
}
