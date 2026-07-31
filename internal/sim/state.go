package sim

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

type Plant struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	SpeciesID        SpeciesID   `json:"species_id"`
	Stage            Stage       `json:"stage"`
	Health           float64     `json:"health"`
	Moisture         float64     `json:"moisture"`
	Growth           float64     `json:"growth"`
	Temperament      Temperament `json:"temperament"`
	LastWatered      time.Time   `json:"last_watered"`
	PlantedAt        time.Time   `json:"planted_at"`
	OverwaterHours   float64     `json:"overwater_hours"`
	LastCareDay      string      `json:"last_care_day"`
	HadGoodCareToday bool        `json:"had_good_care_today"`
	CareStreak       int         `json:"care_streak"`
}

type MemorialEntry struct {
	Name      string    `json:"name"`
	SpeciesID SpeciesID `json:"species_id"`
	DiedAt    time.Time `json:"died_at"`
	Cause     string    `json:"cause"`
}

type GameState struct {
	SchemaVersion   int             `json:"schema_version"`
	Coins           int             `json:"coins"`
	ShelfCapacity   int             `json:"shelf_capacity"`
	Plants          []Plant         `json:"plants"`
	Memorial        []MemorialEntry `json:"memorial"`
	UnlockedSpecies []SpeciesID     `json:"unlocked_species"`
	LastSavedAt     time.Time       `json:"last_saved_at"`
	LastPlayedAt    time.Time       `json:"last_played_at"`
	TimeScale       float64         `json:"time_scale"`
	Seed            int64           `json:"seed"`
}

const CurrentSchemaVersion = 2

func DefaultTimeScale() float64 { return 1.0 }

func NewGameState(now time.Time, seed int64) *GameState {
	if seed == 0 {
		var b [8]byte
		_, _ = rand.Read(b[:])
		seed = int64(binary.LittleEndian.Uint64(b[:]))
	}
	g := &GameState{
		SchemaVersion:   CurrentSchemaVersion,
		Coins:           30,
		ShelfCapacity:   5,
		UnlockedSpecies: StarterSpecies(),
		LastSavedAt:     now,
		LastPlayedAt:    now,
		TimeScale:       DefaultTimeScale(),
		Seed:            seed,
	}
	g.Plants = []Plant{NewStarterPlant(now, seed)}
	return g
}

func NewStarterPlant(now time.Time, seed int64) Plant {
	return Plant{
		ID:          newID(seed, 0),
		Name:        "Gerald",
		SpeciesID:   SpeciesSpiderPlant,
		Stage:       StageSeed,
		Health:      85,
		Moisture:    55,
		Growth:      0,
		Temperament: TemperamentCheerful,
		LastWatered: now,
		PlantedAt:   now,
		LastCareDay: now.Format("2006-01-02"),
	}
}

func newID(seed int64, n int) string {
	return formatID(seed, n)
}

func formatID(seed int64, n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	v := uint64(seed) + uint64(n)*7919
	for i := range b {
		b[i] = chars[v%uint64(len(chars))]
		v /= uint64(len(chars))
	}
	return string(b)
}

func (g *GameState) PlantByID(id string) (*Plant, int) {
	for i := range g.Plants {
		if g.Plants[i].ID == id {
			return &g.Plants[i], i
		}
	}
	return nil, -1
}

func (g *GameState) UsedSlots() int {
	n := 0
	for _, p := range g.Plants {
		if p.Stage != StageDead {
			n++
		}
	}
	return n
}

func (g *GameState) PlantAtSlot(slot int) *Plant {
	if slot < 0 || slot >= g.ShelfCapacity {
		return nil
	}
	if slot >= len(g.Plants) {
		return nil
	}
	p := &g.Plants[slot]
	if !p.IsAlive() {
		return nil
	}
	return p
}

func (g *GameState) FirstAlivePlant() *Plant {
	for i := range g.Plants {
		if g.Plants[i].IsAlive() {
			return &g.Plants[i]
		}
	}
	return nil
}

func (g *GameState) EmptySlotIndex() int {
	for i := 0; i < g.ShelfCapacity; i++ {
		if g.IsSlotEmpty(i) {
			return i
		}
	}
	return -1
}

func (g *GameState) IsSlotEmpty(slot int) bool {
	if slot < 0 || slot >= g.ShelfCapacity {
		return false
	}
	if slot >= len(g.Plants) {
		return true
	}
	return !g.Plants[slot].IsAlive()
}

func (g *GameState) HasSpecies(id SpeciesID) bool {
	for _, s := range g.UnlockedSpecies {
		if s == id {
			return true
		}
	}
	return false
}

func (p *Plant) Species() Species {
	s, ok := GetSpecies(p.SpeciesID)
	if !ok {
		return speciesCatalog[SpeciesSpiderPlant]
	}
	return s
}

func (p *Plant) IsAlive() bool {
	return p.Stage != StageDead
}

func (p *Plant) EffectiveStage() Stage {
	if !p.IsAlive() {
		return StageDead
	}
	if p.Health < 30 && p.Stage != StageSeed {
		return StageWilting
	}
	return p.Stage
}

func scaleDuration(d time.Duration, scale float64) time.Duration {
	if scale <= 0 {
		scale = 1
	}
	return time.Duration(float64(d) * scale)
}
