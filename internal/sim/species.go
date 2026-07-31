package sim

import "time"

type SpeciesID string

const (
	SpeciesCactus      SpeciesID = "cactus"
	SpeciesSpiderPlant SpeciesID = "spider_plant"
	SpeciesSnakePlant  SpeciesID = "snake_plant"
)

type Species struct {
	ID                 SpeciesID
	Name               string
	WaterEvery         time.Duration
	MoistureMin        float64
	MoistureMax        float64
	GrowthDuration     time.Duration
	Cost               int
	BaseDrainPerHour   float64
	OverwaterSensitive bool
}

var speciesCatalog = map[SpeciesID]Species{
	SpeciesCactus: {
		ID: SpeciesCactus, Name: "Cactus",
		WaterEvery: 7 * 24 * time.Hour, MoistureMin: 15, MoistureMax: 35,
		GrowthDuration: 5 * 24 * time.Hour, Cost: 5,
		BaseDrainPerHour: 0.35, OverwaterSensitive: true,
	},
	SpeciesSpiderPlant: {
		ID: SpeciesSpiderPlant, Name: "Spider Plant",
		WaterEvery: 3 * 24 * time.Hour, MoistureMin: 40, MoistureMax: 70,
		GrowthDuration: 7 * 24 * time.Hour, Cost: 10,
		BaseDrainPerHour: 0.8, OverwaterSensitive: false,
	},
	SpeciesSnakePlant: {
		ID: SpeciesSnakePlant, Name: "Snake Plant",
		WaterEvery: 5 * 24 * time.Hour, MoistureMin: 20, MoistureMax: 50,
		GrowthDuration: 6 * 24 * time.Hour, Cost: 12,
		BaseDrainPerHour: 0.5, OverwaterSensitive: false,
	},
}

func GetSpecies(id SpeciesID) (Species, bool) {
	s, ok := speciesCatalog[id]
	return s, ok
}

func StarterSpecies() []SpeciesID {
	return []SpeciesID{SpeciesCactus, SpeciesSpiderPlant, SpeciesSnakePlant}
}

func AllSpecies() []Species {
	out := make([]Species, 0, len(speciesCatalog))
	for _, id := range StarterSpecies() {
		out = append(out, speciesCatalog[id])
	}
	return out
}

func (s Species) InSweetSpot(moisture float64) bool {
	return moisture >= s.MoistureMin && moisture <= s.MoistureMax
}

func (s Species) NeedsWater(moisture float64) bool {
	return moisture < s.MoistureMin
}
