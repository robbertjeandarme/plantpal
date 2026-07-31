package tui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/robbert/plantpal/internal/sim"
)

func press(m *Model, key tea.Key) *Model {
	next, _ := m.Update(tea.KeyPressMsg(key))
	return next.(*Model)
}

func keyText(s string) tea.Key {
	if len(s) == 0 {
		return tea.Key{}
	}
	return tea.Key{Text: s, Code: rune(s[0])}
}

func enterKey() tea.Key { return tea.Key{Code: tea.KeyEnter} }

func newTestModel() *Model {
	m := NewModel(sim.NewGameState(time.Now(), 1), nil, 0, false)
	m.Width, m.Height = 100, 30
	m.Screen = screenShelf
	return &m
}

func TestRemoveAndPlantFlow(t *testing.T) {
	m := newTestModel()
	if m.plantAtSlot(0) == nil {
		t.Fatal("expected starter plant")
	}

	m = press(m, keyText("d"))
	if m.Screen != screenSeedPicker {
		t.Fatalf("after d: screen=%d want seed picker (%d)", m.Screen, screenSeedPicker)
	}
	if m.PlantingSlot != 0 {
		t.Fatalf("planting slot=%d", m.PlantingSlot)
	}

	m.SeedChoice = 2
	m = press(m, enterKey())
	if m.Screen != screenShelf {
		t.Fatalf("after plant: screen=%d", m.Screen)
	}
	p := m.plantAtSlot(0)
	if p == nil || p.SpeciesID != sim.SpeciesSnakePlant {
		t.Fatalf("planted %+v", p)
	}
}

func TestEnterEmptySlotOpensPicker(t *testing.T) {
	m := newTestModel()
	m.Selected = 1

	m = press(m, enterKey())
	if m.Screen != screenSeedPicker {
		t.Fatalf("screen=%d want seed picker", m.Screen)
	}
	if m.PlantingSlot != 1 {
		t.Fatalf("slot=%d", m.PlantingSlot)
	}

	m = press(m, enterKey())
	p := m.plantAtSlot(1)
	if p == nil || p.SpeciesID != sim.SpeciesCactus {
		t.Fatalf("slot 1 plant %+v", p)
	}
}

func TestPlantKeyOnEmptySlot(t *testing.T) {
	m := newTestModel()
	m.Selected = 2
	m = press(m, keyText("p"))
	if m.Screen != screenSeedPicker || m.PlantingSlot != 2 {
		t.Fatalf("screen=%d slot=%d", m.Screen, m.PlantingSlot)
	}
}

func TestSeedPickerRenders(t *testing.T) {
	m := newTestModel()
	m.Screen = screenSeedPicker
	m.PlantingSlot = 0
	out := renderSeedPicker(*m)
	if !strings.Contains(stripANSI(out), "Cactus") || !strings.Contains(stripANSI(out), "Plant a seed") {
		t.Fatalf("seed picker missing content:\n%s", stripANSI(out))
	}
}

func TestRenameFlow(t *testing.T) {
	m := newTestModel()
	m = press(m, enterKey()) // detail
	m = press(m, keyText("r"))
	if m.Screen != screenRename {
		t.Fatalf("screen=%d want rename", m.Screen)
	}
	if m.RenameBuffer != "Gerald" {
		t.Fatalf("buffer=%q", m.RenameBuffer)
	}

	// Clear and type new name
	m.RenameBuffer = ""
	for _, ch := range "Basil" {
		m = press(m, tea.Key{Text: string(ch), Code: ch})
	}
	m = press(m, enterKey())

	if m.Screen != screenDetail {
		t.Fatalf("screen=%d want detail", m.Screen)
	}
	if m.plantAtSlot(0).Name != "Basil" {
		t.Fatalf("name=%q", m.plantAtSlot(0).Name)
	}
}

func TestDemoModeBlocksPlanting(t *testing.T) {
	m := NewModel(sim.DemoGarden(time.Now()), nil, 0, true)
	m.Selected = 4
	m.Width, m.Height = 100, 30
	mp := &m
	mp = press(mp, enterKey())
	if mp.Screen == screenSeedPicker {
		t.Fatal("demo should not open seed picker")
	}
}
