package sim

import (
	"fmt"
	"time"
)

const (
	rootRotThresholdHours = 12.0
	overwaterLevel        = 90.0
	criticalMoisture      = 5.0
	lowMoisture           = 20.0
)

func Advance(g *GameState, until time.Time) []Event {
	if until.Before(g.LastPlayedAt) {
		return nil
	}
	if g.TimeScale == 0 {
		g.LastPlayedAt = until
		return nil
	}

	elapsed := until.Sub(g.LastPlayedAt)
	step := time.Minute
	if elapsed > 48*time.Hour {
		step = 15 * time.Minute
	}

	var events []Event
	t := g.LastPlayedAt
	for t.Before(until) {
		next := t.Add(step)
		if next.After(until) {
			next = until
		}
		simDt := scaleDuration(next.Sub(t), g.TimeScale)
		events = append(events, tick(g, t, simDt)...)
		t = next
	}
	g.LastPlayedAt = until
	return events
}

func tick(g *GameState, now time.Time, dt time.Duration) []Event {
	var events []Event
	hours := dt.Hours()
	dayKey := now.Format("2006-01-02")

	for i := range g.Plants {
		p := &g.Plants[i]
		if !p.IsAlive() {
			continue
		}
		species := p.Species()
		prevStage := p.Stage
		wasWilting := p.Health < 30

		p.Moisture = clamp(p.Moisture-species.BaseDrainPerHour*hours, 0, 100)

		if p.Moisture > overwaterLevel {
			p.OverwaterHours += hours
		} else {
			p.OverwaterHours = 0
		}

		inSweet := species.InSweetSpot(p.Moisture)
		if inSweet {
			p.HadGoodCareToday = true
		}

		switch {
		case p.Moisture < criticalMoisture:
			p.Health -= 2.5 * hours
		case p.Moisture < lowMoisture:
			p.Health -= 0.4 * hours
		case p.OverwaterHours >= rootRotThresholdHours && species.OverwaterSensitive:
			p.Health -= 1.2 * hours
		case p.OverwaterHours >= rootRotThresholdHours*1.5:
			p.Health -= 0.6 * hours
		case inSweet && p.Health < 100:
			p.Health += 0.15 * hours
		}

		if inSweet && p.Stage != StageBlooming {
			growthRate := 100.0 / species.GrowthDuration.Hours()
			p.Growth += growthRate * hours
			if p.Growth >= 100 {
				p.Growth = 0
				advanceStage(p, now, &events)
			}
		}

		if p.LastCareDay != dayKey && dayKey != "" {
			if p.HadGoodCareToday {
				p.CareStreak++
				if p.Stage == StageMature && p.CareStreak >= 3 {
					p.Stage = StageBlooming
					events = append(events, Event{At: now, PlantID: p.ID, Kind: EventBloomed,
						Text: fmt.Sprintf("%s is blooming!", p.Name)})
				}
			} else if p.LastCareDay != "" {
				p.CareStreak = 0
			}
			p.HadGoodCareToday = inSweet
			p.LastCareDay = dayKey
		}

		p.Health = clamp(p.Health, 0, 100)

		if p.Health < 30 && !wasWilting && p.Stage != StageSeed {
			events = append(events, Event{At: now, PlantID: p.ID, Kind: EventWilting,
				Text: fmt.Sprintf("%s is wilting.", p.Name)})
		}
		if wasWilting && p.Health >= 35 {
			events = append(events, Event{At: now, PlantID: p.ID, Kind: EventRecovered,
				Text: fmt.Sprintf("%s is looking better.", p.Name)})
		}
		if p.Stage != prevStage && p.Stage > prevStage {
			events = append(events, Event{At: now, PlantID: p.ID, Kind: EventStageUp,
				Text: fmt.Sprintf("%s grew to %s.", p.Name, p.Stage.String())})
		}
		if p.Health <= 0 {
			killPlant(g, p, now, deathCause(p), &events)
		}
	}
	return events
}

func advanceStage(p *Plant, now time.Time, events *[]Event) {
	switch p.Stage {
	case StageSeed:
		p.Stage = StageSprout
	case StageSprout:
		p.Stage = StageYoung
	case StageYoung:
		p.Stage = StageMature
	}
}

func deathCause(p *Plant) string {
	if p.Moisture < criticalMoisture {
		return "dehydration"
	}
	if p.OverwaterHours >= rootRotThresholdHours {
		return "root rot"
	}
	return "neglect"
}

func killPlant(g *GameState, p *Plant, now time.Time, cause string, events *[]Event) {
	p.Stage = StageDead
	p.Health = 0
	g.Memorial = append(g.Memorial, MemorialEntry{
		Name: p.Name, SpeciesID: p.SpeciesID, DiedAt: now, Cause: cause,
	})
	*events = append(*events, Event{At: now, PlantID: p.ID, Kind: EventDied,
		Text: fmt.Sprintf("%s has passed away from %s.", p.Name, cause)})
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func MoistureLabel(m float64) string {
	switch {
	case m < 20:
		return "dry"
	case m < 50:
		return "low"
	case m < 80:
		return "ok"
	default:
		return "wet"
	}
}

func WaterStatus(p *Plant) string {
	species := p.Species()
	switch {
	case p.Moisture < criticalMoisture:
		return "critical"
	case species.NeedsWater(p.Moisture):
		return "thirsty"
	case p.Moisture > overwaterLevel && species.OverwaterSensitive:
		return "too wet"
	case species.InSweetSpot(p.Moisture):
		return "happy"
	default:
		return "ok"
	}
}
