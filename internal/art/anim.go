package art

import "github.com/charmbracelet/harmonica"

// PlantAnim drives the slow, physical responses of a plant: sagging as its
// health falls and perking up after a drink. Plants are deliberately still
// otherwise — no idle motion.
type PlantAnim struct {
	DroopSpring harmonica.Spring
	PerkSpring  harmonica.Spring
	Droop       float64
	DroopVel    float64
	Perk        float64
	PerkVel     float64
	PerkTarget  float64
	Watering    bool
	WaterTimer  float64
}

func NewPlantAnim() *PlantAnim {
	return &PlantAnim{
		DroopSpring: harmonica.NewSpring(harmonica.FPS(12), 8.0, 0.5),
		PerkSpring:  harmonica.NewSpring(harmonica.FPS(12), 10.0, 0.6),
	}
}

func (a *PlantAnim) Update(dt float64, health float64) {
	droopTarget := (100-health)/100*2.5 + a.Perk
	a.Droop, a.DroopVel = a.DroopSpring.Update(a.Droop, a.DroopVel, droopTarget)
	a.Perk, a.PerkVel = a.PerkSpring.Update(a.Perk, a.PerkVel, a.PerkTarget)

	if a.Watering {
		a.WaterTimer -= dt
		if a.WaterTimer <= 0 {
			a.Watering = false
		}
	}
}

func (a *PlantAnim) TriggerWater() {
	a.Watering = true
	a.WaterTimer = 1.0
	a.PerkTarget = -1.5
}

func (a *PlantAnim) OnWaterComplete() {
	a.PerkTarget = 0
}

// WaterOverlay is a droplet line shown above the plant while it drinks.
func WaterOverlay(timer float64) string {
	if timer <= 0 {
		return ""
	}
	if timer > 0.6 {
		return "💧"
	}
	return " 💧"
}
