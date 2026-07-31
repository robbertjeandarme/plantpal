package art

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/robbert/plantpal/internal/sim"
)

var stageColors = map[sim.Stage]string{
	sim.StageSeed:     "#B0885F",
	sim.StageSprout:   "#8BC34A",
	sim.StageYoung:    "#7CB342",
	sim.StageMature:   "#66BB6A",
	sim.StageBlooming: "#F06292",
	sim.StageWilting:  "#EF9A9A",
	sim.StageDead:     "#8D7C74",
}

func StageColor(s sim.Stage) string {
	if c, ok := stageColors[s]; ok {
		return c
	}
	return "#7CB342"
}

// RenderPortrait draws the plant large, inside a framed "photo" with a
// stage caption underneath.
func RenderPortrait(p sim.Plant, health, moisture, droop float64) string {
	sprite := SpriteDetail(p)
	body := RenderSpriteDroop(sprite, health, moisture, droop)
	stage := p.EffectiveStage()

	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(StageColor(stage))).
		Padding(1, 3).
		Align(lipgloss.Center).
		Render(body)

	caption := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#12210F")).
		Background(lipgloss.Color(StageColor(stage))).
		Bold(true).
		Padding(0, 2).
		Render(stage.Emoji() + "  " + stage.String())

	return lipgloss.JoinVertical(lipgloss.Center, frame, caption)
}

// StageBadge is the compact stage chip used on the shelf.
func StageBadge(stage sim.Stage) string { return StageBadgeFit(stage, 999) }

// StageBadgeFit drops the stage name when there isn't room for it, so a
// narrow card shows a bare emoji rather than wrapping onto a second line.
func StageBadgeFit(stage sim.Stage, width int) string {
	label := stage.Emoji() + " " + stage.String()
	if lipgloss.Width(label) > width {
		label = stage.Emoji()
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(StageColor(stage))).
		Bold(true).
		Render(label)
}

// StageTimeline shows growth progress as a dot trail.
func StageTimeline(p sim.Plant) string {
	stages := []sim.Stage{
		sim.StageSeed, sim.StageSprout, sim.StageYoung,
		sim.StageMature, sim.StageBlooming,
	}
	current := p.Stage
	if p.EffectiveStage() == sim.StageWilting {
		current = sim.StageMature
	}

	done := lipgloss.NewStyle().Foreground(lipgloss.Color("#66BB6A"))
	now := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCA6A")).Bold(true)
	todo := lipgloss.NewStyle().Foreground(lipgloss.Color("#7B8F9E"))

	var parts []string
	for _, s := range stages {
		label := stageShort(s)
		switch {
		case s == current:
			parts = append(parts, now.Render("◉ "+label))
		case int(s) < int(current):
			parts = append(parts, done.Render("● "+label))
		default:
			parts = append(parts, todo.Render("○ "+label))
		}
	}
	return strings.Join(parts, "  ")
}

func stageShort(s sim.Stage) string {
	switch s {
	case sim.StageSeed:
		return "seed"
	case sim.StageSprout:
		return "sprout"
	case sim.StageYoung:
		return "young"
	case sim.StageMature:
		return "mature"
	case sim.StageBlooming:
		return "bloom"
	default:
		return "?"
	}
}
