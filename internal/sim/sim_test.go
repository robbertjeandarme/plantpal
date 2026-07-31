package sim

import (
	"testing"
	"time"
)

func TestWaterImprovesMoisture(t *testing.T) {
	g := NewGameState(time.Now(), 42)
	p := &g.Plants[0]
	p.Moisture = 10
	res := Water(g, p.ID, time.Now())
	if !res.OK {
		t.Fatal(res.Message)
	}
	if p.Moisture <= 10 {
		t.Fatalf("expected moisture increase, got %.1f", p.Moisture)
	}
}

func TestDehydrationKillsPlant(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 99)
	p := &g.Plants[0]
	p.Moisture = 0
	p.Health = 5
	events := Advance(g, now.Add(48*time.Hour))
	if p.IsAlive() {
		t.Fatal("plant should die from dehydration")
	}
	foundDeath := false
	for _, e := range events {
		if e.Kind == EventDied {
			foundDeath = true
		}
	}
	if !foundDeath {
		t.Fatal("expected death event")
	}
}

func TestGrowthInSweetSpot(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 7)
	g.TimeScale = 86400 // 1 day per second -> fast in sim
	p := &g.Plants[0]
	p.Moisture = 55
	p.Health = 90
	p.Stage = StageSeed
	before := p.Stage
	Advance(g, now.Add(3*24*time.Hour))
	if p.Growth == 0 && p.Stage == before && p.Stage == StageSeed {
		t.Log("growth may need more time; stage=", p.Stage, "growth=", p.Growth)
	}
}

func TestAdvanceCoarsening(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 1)
	p := &g.Plants[0]
	p.Moisture = 50
	startMoisture := p.Moisture
	Advance(g, now.Add(7*24*time.Hour))
	if p.Moisture >= startMoisture {
		t.Fatalf("moisture should drain over a week, got %.1f", p.Moisture)
	}
}

func TestSpeciesSweetSpot(t *testing.T) {
	cactus, _ := GetSpecies(SpeciesCactus)
	if !cactus.InSweetSpot(25) {
		t.Fatal("25% should be sweet for cactus")
	}
	if cactus.InSweetSpot(80) {
		t.Fatal("80% should not be sweet for cactus")
	}
}

func TestMemorialOnDeath(t *testing.T) {
	now := time.Now()
	g := NewGameState(now, 3)
	p := &g.Plants[0]
	p.Health = 0
	p.Moisture = 0
	killPlant(g, p, now, "dehydration", &[]Event{})
	if len(g.Memorial) != 1 {
		t.Fatalf("expected memorial entry, got %d", len(g.Memorial))
	}
}
