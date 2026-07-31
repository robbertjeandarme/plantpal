package tui

import "charm.land/lipgloss/v2"

// Colors are tuned for contrast on a dark terminal: mid-to-high lightness
// so nothing disappears into the background.
var (
	colorPrimary   = "#8BC34A"
	colorSecondary = "#6B8E5A"
	colorAccent    = "#FFCA6A"
	colorDanger    = "#FF8A80"
	colorMuted     = "#9FB3AE"
	colorWater     = "#5FC9F8"
	colorTrack     = "#44544A"
	colorText      = "#EDF6ED"
	colorInk       = "#12210F"
)

func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorPrimary))
}

func PanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorSecondary)).
		Padding(0, 1)
}

func SubtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
}

func TextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
}

func DangerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorDanger)).Bold(true)
}

func AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
}

func DemoBadgeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorInk)).
		Background(lipgloss.Color(colorAccent)).
		Bold(true).
		Padding(0, 1)
}

// HealthBar shifts green -> amber -> red as the plant declines, so the bar
// itself signals trouble without reading the number.
func HealthBar(value float64, width int) string {
	fill := colorPrimary
	switch {
	case value < 30:
		fill = colorDanger
	case value < 60:
		fill = colorAccent
	}
	return bar(value, width, fill)
}

func MoistureBar(value float64, width int) string {
	fill := colorWater
	if value < 20 {
		fill = colorDanger
	}
	return bar(value, width, fill)
}

func bar(value float64, width int, fill string) string {
	if width < 4 {
		width = 10
	}
	n := int(value / 100 * float64(width))
	if n < 0 {
		n = 0
	}
	if n > width {
		n = width
	}
	filled := lipgloss.NewStyle().Foreground(lipgloss.Color(fill)).
		Render(repeat("█", n))
	track := lipgloss.NewStyle().Foreground(lipgloss.Color(colorTrack)).
		Render(repeat("─", width-n))
	return filled + track
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
