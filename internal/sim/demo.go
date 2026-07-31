package sim

import "time"

func DemoGarden(now time.Time) *GameState {
	g := &GameState{
		SchemaVersion:   CurrentSchemaVersion,
		Coins:           42,
		ShelfCapacity:   5,
		UnlockedSpecies: StarterSpecies(),
		LastSavedAt:     now,
		LastPlayedAt:    now,
		TimeScale:       0,
		Seed:            99,
	}
	day := now.Format("2006-01-02")

	g.Plants = []Plant{
		{ID: formatID(99, 0), Name: "Spike", SpeciesID: SpeciesCactus, Stage: StageSeed,
			Health: 88, Moisture: 22, Temperament: TemperamentStoic,
			LastWatered: now.Add(-5 * 24 * time.Hour), PlantedAt: now.Add(-2 * 24 * time.Hour), LastCareDay: day},
		{ID: formatID(99, 1), Name: "Susan", SpeciesID: SpeciesSpiderPlant, Stage: StageSprout,
			Health: 92, Moisture: 58, Temperament: TemperamentCheerful,
			LastWatered: now.Add(-6 * time.Hour), PlantedAt: now.Add(-3 * 24 * time.Hour), LastCareDay: day},
		{ID: formatID(99, 2), Name: "Severus", SpeciesID: SpeciesSnakePlant, Stage: StageYoung,
			Health: 78, Moisture: 35, Temperament: TemperamentStoic,
			LastWatered: now.Add(-2 * 24 * time.Hour), PlantedAt: now.Add(-5 * 24 * time.Hour), LastCareDay: day},
		{ID: formatID(99, 3), Name: "Gerald", SpeciesID: SpeciesSpiderPlant, Stage: StageMature,
			Health: 25, Moisture: 8, Temperament: TemperamentDramatic,
			LastWatered: now.Add(-3 * 24 * time.Hour), PlantedAt: now.Add(-10 * 24 * time.Hour), LastCareDay: day},
		{ID: formatID(99, 4), Name: "Bartholomew", SpeciesID: SpeciesCactus, Stage: StageBlooming,
			Health: 95, Moisture: 28, Temperament: TemperamentCheerful,
			LastWatered: now.Add(-12 * time.Hour), PlantedAt: now.Add(-14 * 24 * time.Hour), LastCareDay: day},
	}
	return g
}
