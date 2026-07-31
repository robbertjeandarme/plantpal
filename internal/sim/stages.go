package sim

import "fmt"

type Stage int

const (
	StageSeed Stage = iota
	StageSprout
	StageYoung
	StageMature
	StageBlooming
	StageWilting
	StageDead
)

func (s Stage) String() string {
	switch s {
	case StageSeed:
		return "Seed"
	case StageSprout:
		return "Sprout"
	case StageYoung:
		return "Young"
	case StageMature:
		return "Mature"
	case StageBlooming:
		return "Blooming"
	case StageWilting:
		return "Wilting"
	case StageDead:
		return "Dead"
	default:
		return fmt.Sprintf("Stage(%d)", int(s))
	}
}

func (s Stage) Emoji() string {
	switch s {
	case StageSeed:
		return "🌰"
	case StageSprout:
		return "🌱"
	case StageYoung:
		return "🌿"
	case StageMature:
		return "🪴"
	case StageBlooming:
		return "🌸"
	case StageWilting:
		return "🥴"
	case StageDead:
		return "🥀"
	default:
		return "🌿"
	}
}

func (s Stage) Index() int {
	if s == StageWilting {
		return int(StageMature)
	}
	if s == StageBlooming {
		return int(StageMature) + 1
	}
	if s == StageDead {
		return -1
	}
	return int(s)
}
