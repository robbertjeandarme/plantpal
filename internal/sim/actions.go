package sim

import (
	"fmt"
	"strings"
	"time"
)

const maxPlantNameLen = 24

type ActionResult struct {
	OK      bool
	Message string
	PlantID string
}

func Water(g *GameState, plantID string, now time.Time) ActionResult {
	p, _ := g.PlantByID(plantID)
	if p == nil || !p.IsAlive() {
		return ActionResult{OK: false, Message: "Plant not found."}
	}
	species := p.Species()
	add := 35.0
	if species.OverwaterSensitive {
		add = 25.0
	}
	p.Moisture = clamp(p.Moisture+add, 0, 100)
	p.LastWatered = now
	if species.InSweetSpot(p.Moisture) {
		return ActionResult{OK: true, Message: fmt.Sprintf("%s watered — looking good.", p.Name), PlantID: p.ID}
	}
	if p.Moisture > species.MoistureMax {
		return ActionResult{OK: true, Message: fmt.Sprintf("%s watered — easy on the water next time.", p.Name), PlantID: p.ID}
	}
	return ActionResult{OK: true, Message: fmt.Sprintf("%s watered.", p.Name), PlantID: p.ID}
}

// RemovePlant clears a shelf slot so a new seed can be planted there.
// Unlike natural death this does not add a memorial entry.
func RemovePlant(g *GameState, slot int) ActionResult {
	if slot < 0 || slot >= g.ShelfCapacity {
		return ActionResult{OK: false, Message: "Invalid slot."}
	}
	if g.IsSlotEmpty(slot) {
		return ActionResult{OK: false, Message: "Slot is already empty."}
	}
	name := g.Plants[slot].Name
	id := g.Plants[slot].ID
	// Blank tombstone — no name, so the death screen knows this was manual.
	g.Plants[slot] = Plant{Stage: StageDead}
	return ActionResult{OK: true, Message: fmt.Sprintf("%s removed.", name), PlantID: id}
}

// PlantSeed starts a new plant in the given shelf slot.
func PlantSeed(g *GameState, slot int, speciesID SpeciesID, now time.Time) ActionResult {
	if slot < 0 || slot >= g.ShelfCapacity {
		return ActionResult{OK: false, Message: "Invalid slot."}
	}
	if !g.IsSlotEmpty(slot) {
		return ActionResult{OK: false, Message: "Slot is occupied."}
	}
	species, ok := GetSpecies(speciesID)
	if !ok {
		return ActionResult{OK: false, Message: "Unknown species."}
	}
	if !g.HasSpecies(speciesID) {
		return ActionResult{OK: false, Message: "Species not available."}
	}

	plant := Plant{
		ID:          formatID(g.Seed, slot*997+int(now.Unix()%100000)),
		Name:        species.Name,
		SpeciesID:   speciesID,
		Stage:       StageSeed,
		Health:      90,
		Moisture:    (species.MoistureMin + species.MoistureMax) / 2,
		Temperament: Temperament(int(g.Seed+int64(slot)) % 4),
		LastWatered: now,
		PlantedAt:   now,
		LastCareDay: now.Format("2006-01-02"),
	}

	if slot < len(g.Plants) {
		g.Plants[slot] = plant
	} else {
		for len(g.Plants) < slot {
			g.Plants = append(g.Plants, Plant{Stage: StageDead})
		}
		g.Plants = append(g.Plants, plant)
	}
	return ActionResult{OK: true, Message: fmt.Sprintf("%s planted!", species.Name), PlantID: plant.ID}
}

// RenamePlant sets a new display name for a living plant.
func RenamePlant(g *GameState, plantID, name string) ActionResult {
	name = strings.TrimSpace(name)
	if name == "" {
		return ActionResult{OK: false, Message: "Name can't be empty."}
	}
	if len([]rune(name)) > maxPlantNameLen {
		return ActionResult{OK: false, Message: fmt.Sprintf("Name too long (max %d).", maxPlantNameLen)}
	}
	p, _ := g.PlantByID(plantID)
	if p == nil || !p.IsAlive() {
		return ActionResult{OK: false, Message: "Plant not found."}
	}
	p.Name = name
	return ActionResult{OK: true, Message: fmt.Sprintf("Renamed to %s.", name), PlantID: p.ID}
}

func BuyPlant(g *GameState, speciesID SpeciesID, name string, now time.Time) ActionResult {
	species, ok := GetSpecies(speciesID)
	if !ok {
		return ActionResult{OK: false, Message: "Unknown species."}
	}
	slot := g.EmptySlotIndex()
	if slot < 0 {
		return ActionResult{OK: false, Message: "No empty slots."}
	}
	if g.Coins < species.Cost {
		return ActionResult{OK: false, Message: fmt.Sprintf("Need %d coins.", species.Cost)}
	}
	g.Coins -= species.Cost
	res := PlantSeed(g, slot, speciesID, now)
	if res.OK && name != "" && name != species.Name {
		if p := g.PlantAtSlot(slot); p != nil {
			p.Name = name
		}
	}
	return res
}
