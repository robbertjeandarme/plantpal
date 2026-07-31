package tui

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/robbert/plantpal/internal/sim"
)

func demoModel(w, h int) Model {
	now := time.Now()
	m := NewModel(sim.DemoGarden(now), nil, 0, true)
	m.Width, m.Height = w, h
	m.Screen = screenShelf
	return m
}

// TestShelfCardsAreUniform is the guarantee that the shelf reads as an even
// grid: every card, occupied or empty, renders at identical dimensions.
func TestShelfCardsAreUniform(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {100, 30}, {140, 40}, {160, 50}} {
		m := demoModel(size[0], size[1])
		// Leave one slot empty so empty and occupied cards are compared.
		m.Game.Plants = m.Game.Plants[:len(m.Game.Plants)-1]

		var wantW, wantH int
		for slot := 0; slot < m.Game.ShelfCapacity; slot++ {
			card := renderShelfSlot(m, slot, 18)
			gotW, gotH := lipgloss.Width(card), lipgloss.Height(card)
			if slot == 0 {
				wantW, wantH = gotW, gotH
				continue
			}
			if gotW != wantW || gotH != wantH {
				t.Errorf("term %dx%d slot %d: card is %dx%d, want %dx%d",
					size[0], size[1], slot, gotW, gotH, wantW, wantH)
			}
		}
	}
}

// TestShelfCardHeightIsStageIndependent checks a seedling and a blooming
// plant produce the same card, so growth never reflows the shelf.
func TestShelfCardHeightIsStageIndependent(t *testing.T) {
	m := demoModel(120, 40)
	stages := []sim.Stage{
		sim.StageSeed, sim.StageSprout, sim.StageYoung,
		sim.StageMature, sim.StageBlooming,
	}
	species := []sim.SpeciesID{
		sim.SpeciesCactus, sim.SpeciesSpiderPlant, sim.SpeciesSnakePlant,
	}

	var wantW, wantH int
	for _, sp := range species {
		for _, st := range stages {
			m.Game.Plants[0].SpeciesID = sp
			m.Game.Plants[0].Stage = st
			card := renderShelfSlot(m, 0, 18)
			gotW, gotH := lipgloss.Width(card), lipgloss.Height(card)
			if wantW == 0 {
				wantW, wantH = gotW, gotH
				continue
			}
			if gotW != wantW || gotH != wantH {
				t.Errorf("%s/%s: card is %dx%d, want %dx%d", sp, st, gotW, gotH, wantW, wantH)
			}
		}
	}
}

// TestPotsShareABaseline verifies pots land on the same row across cards,
// which is what makes the row look like a real shelf.
func TestPotsShareABaseline(t *testing.T) {
	m := demoModel(120, 40)
	baseline := -1
	for _, st := range []sim.Stage{sim.StageSeed, sim.StageYoung, sim.StageBlooming} {
		m.Game.Plants[0].Stage = st
		lines := strings.Split(renderShelfSlot(m, 0, 18), "\n")

		row := -1
		for i, l := range lines {
			if strings.Contains(stripANSI(l), "╲_") {
				row = i
			}
		}
		if row < 0 {
			t.Fatalf("%s: no pot base found", st)
		}
		if baseline == -1 {
			baseline = row
			continue
		}
		if row != baseline {
			t.Errorf("%s: pot base on row %d, want %d", st, row, baseline)
		}
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// TestDetailHidesGrowthProgress keeps the surprise of growth intact.
func TestDetailHidesGrowthProgress(t *testing.T) {
	m := demoModel(140, 44)
	m.Game.Plants[0].Growth = 47
	out := renderDetail(m, &m.Game.Plants[0])
	for _, banned := range []string{"next stage", "47%", "growth"} {
		if strings.Contains(strings.ToLower(out), strings.ToLower(banned)) {
			t.Errorf("detail view still shows %q", banned)
		}
	}
}

// TestShelfDoesNotWrap checks the whole row of cards fits on one line. If a
// card overflows, lipgloss wraps it onto the next line and the shelf breaks
// apart, so every card must share a single top-border row.
func TestShelfDoesNotWrap(t *testing.T) {
	for _, w := range []int{80, 90, 100, 120, 160, 200} {
		m := demoModel(w, 34)
		n := m.Game.ShelfCapacity

		var joinRows, joins int
		for _, line := range strings.Split(renderShelf(m), "\n") {
			if c := strings.Count(stripANSI(line), "╮╭"); c > 0 {
				joinRows++
				joins += c
			}
			if got := lipgloss.Width(line); got > w {
				t.Errorf("term %d: line %d wide", w, got)
			}
		}
		if joinRows != 1 || joins != n-1 {
			t.Errorf("term %d: cards wrapped — top border spans %d rows with %d joins (want 1 row, %d)",
				w, joinRows, joins, n-1)
		}
	}
}

// TestDetailStatsDoNotWrap keeps each stat on its own line; a wrapped bar
// or timeline makes the panel look broken.
func TestDetailStatsDoNotWrap(t *testing.T) {
	for _, w := range []int{80, 100, 120, 160, 200} {
		m := demoModel(w, 40)
		m.Screen = screenDetail
		for slot := range m.Game.Plants {
			out := renderDetail(m, &m.Game.Plants[slot])
			for _, frag := range []string{"❤ health", "💧 water", " seed", "watered "} {
				if n := countLinesWith(out, frag); n != 1 {
					t.Errorf("term %d plant %d: %q appears on %d lines, want 1",
						w, slot, frag, n)
				}
			}
			for _, line := range strings.Split(out, "\n") {
				if got := lipgloss.Width(line); got > w {
					t.Errorf("term %d: detail line %d wide", w, got)
					break
				}
			}
		}
	}
}

func countLinesWith(s, sub string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(stripANSI(line), sub) {
			n++
		}
	}
	return n
}
