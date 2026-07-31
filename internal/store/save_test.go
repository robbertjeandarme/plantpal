package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/robbert/plantpal/internal/sim"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	now := time.Now()
	g := sim.NewGameState(now, 123)
	g.Coins = 99
	g.Plants[0].Name = "TestPlant"

	if err := Save(g, now); err != nil {
		t.Fatal(err)
	}

	loaded, existed, err := Load(now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if !existed {
		t.Fatal("expected save to exist")
	}
	if loaded.Coins != 99 {
		t.Fatalf("coins=%d", loaded.Coins)
	}
	if loaded.Plants[0].Name != "TestPlant" {
		t.Fatalf("name=%q", loaded.Plants[0].Name)
	}

	path, _ := SavePath()
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "savegame.json" {
		t.Fatal(path)
	}
}
