package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	marginH = 2
	marginV = 1
)

func contentWidth(m Model) int {
	w := m.Width - marginH*2
	if w < 40 {
		w = 40
	}
	return w
}

func contentHeight(m Model) int {
	h := m.Height - marginV*2
	if h < 20 {
		h = 20
	}
	return h
}

func renderHeader(m Model, subtitle string) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorPrimary)).
		Render("🪴 PLANT PAL")

	sub := SubtleStyle().Render(subtitle)
	gap := contentWidth(m) - lipgloss.Width(title) - lipgloss.Width(sub) - 2
	if gap < 1 {
		gap = 1
	}

	line := title + strings.Repeat(" ", gap) + sub

	return lipgloss.NewStyle().
		Width(contentWidth(m)).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colorSecondary)).
		Padding(0, 1).
		Render(line)
}

func renderFooter(m Model, hints string) string {
	return lipgloss.NewStyle().
		Width(contentWidth(m)).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colorSecondary)).
		Padding(0, 1).
		Foreground(lipgloss.Color(colorMuted)).
		Render(hints)
}

func frameScreen(m Model, subtitle, body, footer string) string {
	header := renderHeader(m, subtitle)
	foot := renderFooter(m, footer)

	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(foot)
	bodyH := contentHeight(m) - headerH - footerH
	if bodyH < 1 {
		bodyH = 1
	}

	bodyPanel := lipgloss.NewStyle().
		Width(contentWidth(m)).
		Height(bodyH).
		Padding(0, 1).
		Render(body)

	return lipgloss.NewStyle().Padding(marginV, marginH).Render(
		lipgloss.JoinVertical(lipgloss.Left, header, bodyPanel, foot),
	)
}

func fullScreenCenter(m Model, content string) string {
	w := contentWidth(m)
	h := contentHeight(m)
	panel := PanelStyle().Width(min(w-4, 72)).Padding(1, 2).Render(content)
	return lipgloss.NewStyle().
		Padding(marginV, marginH).
		Render(lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, panel))
}

func panelBox(content string, width int) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorSecondary)).
		Padding(1, 2).
		Width(width).
		Render(content)
}

func statLine(icon, label, bar, value string) string {
	return fmt.Sprintf("  %s %-10s %s  %s", icon, label, bar, value)
}
