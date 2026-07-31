package sim

import (
	"strings"
	"testing"
	"time"
)

func TestRemoveAndPlant(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 1)

	if g.IsSlotEmpty(0) {
		t.Fatal("starter should occupy slot 0")
	}

	res := RemovePlant(g, 0)
	if !res.OK {
		t.Fatal(res.Message)
	}
	if !g.IsSlotEmpty(0) {
		t.Fatal("slot should be empty after remove")
	}
	if g.Plants[0].Name != "" {
		t.Fatal("removed slot should be a blank tombstone")
	}

	res = PlantSeed(g, 0, SpeciesCactus, now)
	if !res.OK {
		t.Fatal(res.Message)
	}
	p := g.PlantAtSlot(0)
	if p == nil || p.SpeciesID != SpeciesCactus || p.Stage != StageSeed {
		t.Fatalf("got %+v", p)
	}
}

func TestPlantSeedOccupiedSlot(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 2)
	res := PlantSeed(g, 0, SpeciesSnakePlant, now)
	if res.OK {
		t.Fatal("should not plant over existing plant")
	}
}

func TestPlantSeedEmptySlotAtEnd(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 3)
	RemovePlant(g, 0)

	res := PlantSeed(g, 2, SpeciesSpiderPlant, now)
	if !res.OK {
		t.Fatal(res.Message)
	}
	if len(g.Plants) != 3 {
		t.Fatalf("len=%d", len(g.Plants))
	}
	if g.Plants[2].SpeciesID != SpeciesSpiderPlant {
		t.Fatalf("species=%s", g.Plants[2].SpeciesID)
	}
}

func TestRenamePlant(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 1)
	id := g.Plants[0].ID

	res := RenamePlant(g, id, "  Basil  ")
	if !res.OK || g.Plants[0].Name != "Basil" {
		t.Fatalf("rename=%v name=%q", res, g.Plants[0].Name)
	}

	res = RenamePlant(g, id, "")
	if res.OK {
		t.Fatal("empty name should fail")
	}

	res = RenamePlant(g, id, strings.Repeat("x", maxPlantNameLen+1))
	if res.OK {
		t.Fatal("long name should fail")
	}
}

func TestRemoveEmptySlot(t *testing.T) {
	g := NewGameState(time.Now(), 4)
	RemovePlant(g, 0)
	res := RemovePlant(g, 0)
	if res.OK {
		t.Fatal("should not remove an already-empty slot")
	}
}
