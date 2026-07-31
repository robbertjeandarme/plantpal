package art

import (
	"strings"
	"testing"

	"github.com/robbert/plantpal/internal/sim"
)

var allStages = []sim.Stage{
	sim.StageSeed, sim.StageSprout, sim.StageYoung,
	sim.StageMature, sim.StageBlooming,
}

var allSpecies = []sim.SpeciesID{
	sim.SpeciesCactus, sim.SpeciesSpiderPlant, sim.SpeciesSnakePlant,
}

// TestEveryPlantHasAPot guards the rule that no plant is ever drawn without
// a pot beneath it.
func TestEveryPlantHasAPot(t *testing.T) {
	for _, sp := range allSpecies {
		for _, stage := range allStages {
			p := sim.Plant{SpeciesID: sp, Stage: stage, Health: 90, Moisture: 50}
			for name, sprite := range map[string]Sprite{
				"shelf":  SpriteShelf(p),
				"detail": SpriteDetail(p),
			} {
				if sprite.Art == "" {
					t.Fatalf("%s/%s %s: no art", sp, stage, name)
				}
				if sprite.PotRows < 3 {
					t.Errorf("%s/%s %s: PotRows=%d, want >=3", sp, stage, name, sprite.PotRows)
				}
				lines := sprite.lines()
				if len(lines) <= sprite.PotRows {
					t.Errorf("%s/%s %s: only pot, no plant", sp, stage, name)
				}
				pot := strings.Join(lines[len(lines)-sprite.PotRows:], "\n")
				if !strings.Contains(pot, "▒") {
					t.Errorf("%s/%s %s: pot has no soil:\n%s", sp, stage, name, pot)
				}
			}
		}
	}
}

// TestSpeciesSilhouettesAreDistinct verifies each species uses its own
// characteristic glyphs rather than sharing generic art.
func TestSpeciesSilhouettesAreDistinct(t *testing.T) {
	want := map[sim.SpeciesID]string{
		sim.SpeciesCactus:      "█", // blocky columnar flesh
		sim.SpeciesSpiderPlant: "╲", // arching strap leaves
		sim.SpeciesSnakePlant:  "║", // stiff upright blades
	}
	for sp, glyph := range want {
		for _, stage := range []sim.Stage{sim.StageYoung, sim.StageMature, sim.StageBlooming} {
			p := sim.Plant{SpeciesID: sp, Stage: stage, Health: 90, Moisture: 50}
			for name, sprite := range map[string]Sprite{
				"shelf":  SpriteShelf(p),
				"detail": SpriteDetail(p),
			} {
				lines := sprite.lines()
				foliage := strings.Join(lines[:len(lines)-sprite.PotRows], "\n")
				if !strings.Contains(foliage, glyph) {
					t.Errorf("%s/%s %s: missing signature glyph %q:\n%s",
						sp, stage, name, glyph, foliage)
				}
			}
		}
	}
}

// TestPaletteHasRealHues catches the integer-division bug that collapsed
// every colour to hue 0.
func TestPaletteHasRealHues(t *testing.T) {
	pal := PaletteFor(95, 60)

	leafH, _, _ := pal.Leaf.Hsl()
	if leafH < 90 || leafH > 160 {
		t.Errorf("healthy leaf hue = %.0f, want green (90-160)", leafH)
	}
	potH, _, _ := pal.Pot.Hsl()
	if potH > 40 {
		t.Errorf("pot hue = %.0f, want terracotta (<40)", potH)
	}
	if pal.Leaf.Hex() == pal.Pot.Hex() {
		t.Error("leaf and pot render identically")
	}

	sick, _, _ := PaletteFor(10, 50).Leaf.Hsl()
	if sick > 60 {
		t.Errorf("dying leaf hue = %.0f, want brown/amber (<60)", sick)
	}
}

// TestSoilRespondsToWater confirms wet soil renders darker than dry soil.
func TestSoilRespondsToWater(t *testing.T) {
	_, _, dry := PaletteFor(80, 5).Soil.Hsl()
	_, _, wet := PaletteFor(80, 95).Soil.Hsl()
	if wet >= dry {
		t.Errorf("wet soil lightness %.2f should be below dry %.2f", wet, dry)
	}
}

func TestPortraitRenders(t *testing.T) {
	for _, sp := range allSpecies {
		for _, st := range allStages {
			p := sim.Plant{Name: "X", SpeciesID: sp, Stage: st, Health: 88, Moisture: 55}
			if RenderPortrait(p, p.Health, p.Moisture, 0) == "" {
				t.Fatalf("%s/%s empty portrait", sp, st)
			}
		}
	}
}

// TestDroopKeepsPotStill checks a sagging plant loses height at the top
// while the sprite's overall height and pot placement stay fixed.
func TestDroopKeepsPotStill(t *testing.T) {
	sp := SpriteDetail(sim.Plant{SpeciesID: sim.SpeciesSpiderPlant, Stage: sim.StageMature})
	healthy := RenderSpriteDroop(sp, 95, 60, 0)
	wilted := RenderSpriteDroop(sp, 20, 5, 2)

	hl, wl := strings.Count(healthy, "\n"), strings.Count(wilted, "\n")
	if hl != wl {
		t.Errorf("height changed with droop: %d vs %d", hl, wl)
	}
	strip := func(s string) []string {
		var out []string
		for _, l := range strings.Split(s, "\n") {
			out = append(out, stripANSI(l))
		}
		return out
	}
	h, w := strip(healthy), strip(wilted)
	for i := 1; i <= 3; i++ {
		hp, wp := h[len(h)-i], w[len(w)-i]
		if hp != wp {
			t.Errorf("pot row -%d moved:\n healthy %q\n wilted  %q", i, hp, wp)
		}
	}
	if strings.TrimSpace(w[0]) != "" {
		t.Errorf("droop should clear the top row, got %q", w[0])
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
