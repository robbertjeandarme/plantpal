package sim

import "time"

type EventKind int

const (
	EventStageUp EventKind = iota
	EventStageDown
	EventBloomed
	EventWilting
	EventDied
	EventThirsty
	EventOverwatered
	EventCareStreak
	EventRecovered
)

type Event struct {
	At      time.Time `json:"at"`
	PlantID string    `json:"plant_id"`
	Kind    EventKind `json:"kind"`
	Text    string    `json:"text"`
}

func (k EventKind) String() string {
	switch k {
	case EventStageUp:
		return "stage_up"
	case EventBloomed:
		return "bloomed"
	case EventWilting:
		return "wilting"
	case EventDied:
		return "died"
	case EventThirsty:
		return "thirsty"
	case EventOverwatered:
		return "overwatered"
	case EventCareStreak:
		return "care_streak"
	case EventRecovered:
		return "recovered"
	default:
		return "other"
	}
}
