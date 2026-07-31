package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robbert/plantpal/internal/sim"
)

const saveFileName = "savegame.json"

var ErrNotFound = errors.New("save file not found")

func SaveDir() (string, error) {
	return os.Getwd()
}

func SavePath() (string, error) {
	dir, err := SaveDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, saveFileName), nil
}

func Load(now time.Time) (*sim.GameState, bool, error) {
	path, err := SavePath()
	if err != nil {
		return nil, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return sim.NewGameState(now, 0), false, nil
		}
		return nil, false, err
	}
	var g sim.GameState
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, false, fmt.Errorf("parse save: %w", err)
	}
	if g.SchemaVersion == 0 {
		g.SchemaVersion = sim.CurrentSchemaVersion
	}
	if g.ShelfCapacity == 0 {
		g.ShelfCapacity = 5
	}
	if len(g.UnlockedSpecies) == 0 {
		g.UnlockedSpecies = sim.StarterSpecies()
	}
	if g.TimeScale == 0 {
		g.TimeScale = sim.DefaultTimeScale()
	}
	return &g, true, nil
}

func Save(g *sim.GameState, now time.Time) error {
	g.LastSavedAt = now
	dir, err := SaveDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, saveFileName)
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func DeleteSave() error {
	path, err := SavePath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
