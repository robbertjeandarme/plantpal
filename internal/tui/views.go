package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/robbert/plantpal/internal/art"
	"github.com/robbert/plantpal/internal/sim"
)

type screen int

const (
	screenWelcome screen = iota
	screenShelf
	screenDetail
	screenDeath
	screenSeedPicker
	screenRename
)

type Model struct {
	Game          *sim.GameState
	Width         int
	Height        int
	Screen        screen
	Selected      int
	Events        []sim.Event
	AwayDuration  time.Duration
	Status        string
	Anims         map[string]*art.PlantAnim
	LastTick      time.Time
	DeathPlant    *sim.Plant
	DeathCause    string
	Quitting      bool
	DemoMode      bool
	PlantingSlot  int // shelf slot being planted into; -1 when idle
	SeedChoice    int // index into Game.UnlockedSpecies
	RenameBuffer  string
	RenamePlantID string
	RenameFrom    screen
}

func NewModel(g *sim.GameState, events []sim.Event, away time.Duration, demo bool) Model {
	anims := make(map[string]*art.PlantAnim)
	for _, p := range g.Plants {
		if p.IsAlive() {
			anims[p.ID] = art.NewPlantAnim()
		}
	}
	m := Model{
		Game: g, Events: events, AwayDuration: away,
		Anims: anims, LastTick: time.Now(), DemoMode: demo,
		PlantingSlot: -1,
	}
	if demo {
		m.Screen = screenShelf
	} else if len(events) > 0 && away > time.Minute {
		m.Screen = screenWelcome
	} else {
		m.Screen = screenShelf
	}
	if !demo {
		for _, e := range events {
			if e.Kind == sim.EventDied {
				if p, _ := g.PlantByID(e.PlantID); p != nil {
					cp := *p
					m.DeathPlant = &cp
					m.DeathCause = "unknown"
					if len(g.Memorial) > 0 {
						m.DeathCause = g.Memorial[len(g.Memorial)-1].Cause
					}
					m.Screen = screenDeath
					break
				}
			}
		}
	}
	return m
}

func (m Model) selectedPlant() *sim.Plant       { return m.Game.PlantAtSlot(m.Selected) }
func (m Model) plantAtSlot(slot int) *sim.Plant { return m.Game.PlantAtSlot(slot) }

const (
	footerShelf      = "← → select   enter/p plant   w water   r rename   d remove   q quit"
	footerShelfEmpty = "← → select   enter/p plant   q quit"
	footerDetail     = "← → switch   w water   r rename   d remove   esc back   q quit"
	footerSeedPicker = "← → choose   enter plant   esc cancel   q quit"
	footerRename     = "enter save   esc cancel   backspace delete"

	// Shelf cards are a fixed size so the row reads as one tidy grid
	// regardless of what each pot happens to be holding.
	shelfArtHeight  = 8 // tallest sprite (7 rows) plus a row for the droplet
	shelfCardHeight = 16
	shelfSlotMinW   = 14
	shelfSlotMaxW   = 22
)

func shelfSubtitle(m Model) string {
	return fmt.Sprintf("%d plant%s",
		m.Game.UsedSlots(), plural(m.Game.UsedSlots()))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func renderWelcome(m Model) string {
	var b strings.Builder
	b.WriteString(TitleStyle().Render("Welcome back"))
	b.WriteString("\n\n")
	if m.AwayDuration > 0 {
		b.WriteString(fmt.Sprintf("Away for %s\n\n", formatDuration(m.AwayDuration)))
	}
	if len(m.Events) == 0 {
		b.WriteString(SubtleStyle().Render("Your plants are fine."))
	} else {
		for _, e := range m.Events {
			if len(b.String()) > 0 && !strings.HasSuffix(b.String(), "\n\n") {
				b.WriteString("\n")
			}
			line := e.Text
			switch e.Kind {
			case sim.EventDied:
				line = DangerStyle().Render(line)
			case sim.EventWilting, sim.EventThirsty:
				line = AccentStyle().Render(line)
			}
			b.WriteString("• " + line + "\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(SubtleStyle().Render("any key"))
	return fullScreenCenter(m, b.String())
}

func renderShelf(m Model) string {
	n := m.Game.ShelfCapacity
	// frameScreen pads the body by one column on each side.
	w := contentWidth(m) - 2

	// lipgloss Width() covers the border, so slots tile at exactly slotW.
	slotW := w / n
	if slotW < shelfSlotMinW {
		slotW = shelfSlotMinW
	}
	if slotW > shelfSlotMaxW {
		slotW = shelfSlotMaxW
	}

	cells := make([]string, n)
	for i := 0; i < n; i++ {
		cells[i] = renderShelfSlot(m, i, slotW)
	}

	shelfH := max(10, contentHeight(m)-4)
	body := lipgloss.Place(w, shelfH, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Top, cells...))

	if m.Status != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "",
			AccentStyle().Render(m.Status))
	}

	return frameScreen(m, shelfSubtitle(m), body, footerForShelf(m))
}

func footerForShelf(m Model) string {
	if m.DemoMode {
		return footerShelf
	}
	if m.plantAtSlot(m.Selected) == nil {
		return footerShelfEmpty
	}
	return footerShelf
}

func renderShelfSlot(m Model, slot, width int) string {
	selected := slot == m.Selected
	p := m.plantAtSlot(slot)

	border := lipgloss.Color(colorSecondary)
	if selected {
		border = lipgloss.Color(colorAccent)
	}

	innerW := width - 4 // width covers the border (2) and padding (2)
	// The 💧 icon plus its space costs 3 columns; anything wider than the
	// card wraps onto a second line and knocks the whole shelf out of step.
	barW := min(innerW-3, 8)
	if barW < 4 {
		barW = 4
	}

	var rows []string
	if p != nil {
		droop := 0.0
		watering := ""
		if a, ok := m.Anims[p.ID]; ok {
			droop = a.Droop
			watering = art.WaterOverlay(a.WaterTimer)
		}
		plant := art.RenderSpriteDroop(art.SpriteShelf(*p), p.Health, p.Moisture, droop)
		if watering != "" {
			plant = watering + "\n" + plant
		}

		status := sim.WaterStatus(p)
		statusStyle := SubtleStyle()
		switch status {
		case "critical":
			statusStyle = DangerStyle()
		case "thirsty", "too wet":
			statusStyle = AccentStyle()
		}

		name := TextStyle()
		if selected {
			name = name.Bold(true)
		}

		rows = []string{
			shelfArtBox(plant, innerW),
			"",
			name.Render(truncate(p.Name, innerW)),
			art.StageBadgeFit(p.EffectiveStage(), innerW),
			SubtleStyle().Render(truncate(p.Species().Name, innerW)),
			"",
			fmt.Sprintf("❤ %s", HealthBar(p.Health, barW)),
			fmt.Sprintf("💧 %s", MoistureBar(p.Moisture, barW)),
			statusStyle.Render(status),
		}
	} else {
		rows = []string{
			shelfArtBox(art.RenderSprite(art.EmptyPot(), 50, 25), innerW),
			"",
			SubtleStyle().Render("empty"),
		}
		if selected && !m.DemoMode {
			rows = append(rows, AccentStyle().Render("enter / p"))
		}
	}

	// Pad to a fixed row count rather than leaning on Height(), so an empty
	// pot occupies exactly as much space as a blooming plant.
	for lipgloss.Height(strings.Join(rows, "\n")) < shelfCardHeight {
		rows = append(rows, "")
	}

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Render(lipgloss.JoinVertical(lipgloss.Center, rows...))
}

// shelfArtBox pins the art to a fixed-height box anchored at the bottom, so
// every pot on the shelf rests on the same line no matter how tall the
// plant above it has grown.
func shelfArtBox(plant string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(shelfArtHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Bottom).
		Render(plant)
}

func truncate(s string, width int) string {
	if width < 2 || lipgloss.Width(s) <= width {
		return s
	}
	r := []rune(s)
	if len(r) > width-1 {
		r = r[:width-1]
	}
	return string(r) + "…"
}

func renderDetail(m Model, p *sim.Plant) string {
	if p == nil {
		body := SubtleStyle().Render("Empty slot.")
		if !m.DemoMode {
			body = lipgloss.JoinVertical(lipgloss.Center,
				body, "", AccentStyle().Render("enter to plant a seed"))
		}
		return frameScreen(m, shelfSubtitle(m), body, footerShelfEmpty)
	}

	w := contentWidth(m)
	species := p.Species()
	stage := p.EffectiveStage()

	droop := 0.0
	if a, ok := m.Anims[p.ID]; ok {
		droop = a.Droop
	}

	portrait := art.RenderPortrait(*p, p.Health, p.Moisture, droop)
	if a, ok := m.Anims[p.ID]; ok && a.Watering {
		portrait = art.WaterOverlay(a.WaterTimer) + "\n" + portrait
	}

	// The portrait is narrow, so give the stats panel an even share rather
	// than letting its rows wrap.
	leftW := w / 2
	rightW := w - leftW - 2

	// statLine spends 18 columns on its icon, label and gutters; reserve a
	// further 12 for the trailing value.
	barW := min(18, rightW-6-18-12)
	if barW < 6 {
		barW = 6
	}

	left := lipgloss.NewStyle().
		Width(leftW).
		Height(max(14, contentHeight(m)-12)).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(portrait)

	var stats strings.Builder
	stats.WriteString(TitleStyle().Render(p.Name))
	stats.WriteString("\n")
	stats.WriteString(SubtleStyle().Render(species.Name))
	stats.WriteString("\n\n")
	stats.WriteString(statLine("❤", "health", HealthBar(p.Health, barW), fmt.Sprintf("%.0f%%", p.Health)))
	stats.WriteString("\n")
	stats.WriteString(statLine("💧", "water", MoistureBar(p.Moisture, barW),
		fmt.Sprintf("%s · %.0f%%", sim.MoistureLabel(p.Moisture), p.Moisture)))
	stats.WriteString("\n\n")

	status := sim.WaterStatus(p)
	statusStyle := TextStyle()
	switch status {
	case "critical":
		statusStyle = DangerStyle()
	case "thirsty", "too wet":
		statusStyle = AccentStyle()
	}
	stats.WriteString(SubtleStyle().Render("status  ") + statusStyle.Render(status))
	if !p.LastWatered.IsZero() {
		stats.WriteString("\n")
		stats.WriteString(SubtleStyle().Render("watered " + formatDuration(time.Since(p.LastWatered))))
	}
	stats.WriteString("\n\n")
	stats.WriteString(art.StageTimeline(*p))

	right := lipgloss.NewStyle().
		Width(rightW).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorSecondary)).
		Render(stats.String())

	mainH := max(16, contentHeight(m)-8)
	main := lipgloss.NewStyle().Height(mainH).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, left, right),
	)

	quote := lipgloss.NewStyle().
		Width(w).
		Italic(true).
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 2).
		Align(lipgloss.Center).
		Render(`"` + sim.PlantDialogue(p) + `"`)

	body := lipgloss.JoinVertical(lipgloss.Left, main, "", quote)
	if m.Status != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", AccentStyle().Render(m.Status))
	}

	subtitle := fmt.Sprintf("%s · %s", stage.String(), species.Name)
	return frameScreen(m, subtitle, body, footerDetail)
}

func renderSeedPicker(m Model) string {
	species := m.Game.UnlockedSpecies
	if len(species) == 0 {
		return frameScreen(m, "choose a seed", SubtleStyle().Render("No seeds available."),
			footerSeedPicker)
	}

	w := contentWidth(m) - 2
	n := len(species)
	slotW := w / n
	if slotW < 16 {
		slotW = 16
	}
	if slotW > 28 {
		slotW = 28
	}

	cells := make([]string, n)
	for i, id := range species {
		cells[i] = renderSeedOption(m, id, i, slotW)
	}

	subtitle := "choose a seed"
	bodyH := max(12, contentHeight(m)-6)
	title := TitleStyle().Render("Plant a seed") + "\n" + SubtleStyle().Render(fmt.Sprintf("slot %d", m.PlantingSlot+1))
	body := lipgloss.JoinVertical(lipgloss.Center,
		title, "",
		lipgloss.Place(w, bodyH, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(lipgloss.Top, cells...)),
	)

	if m.Status != "" {
		body = lipgloss.JoinVertical(lipgloss.Center, body, "",
			AccentStyle().Render(m.Status))
	}
	return frameScreen(m, subtitle, body, footerSeedPicker)
}

func renderRename(m Model) string {
	p, _ := m.Game.PlantByID(m.RenamePlantID)
	speciesLabel := "plant"
	if p != nil {
		speciesLabel = p.Species().Name
	}

	fieldW := min(contentWidth(m)-8, 40)
	input := m.RenameBuffer
	if input == "" {
		input = SubtleStyle().Render("…")
	}
	cursor := AccentStyle().Render("▌")
	line := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorAccent)).
		Padding(0, 1).
		Width(fieldW).
		Render(input + cursor)

	body := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle().Render("Rename plant"),
		"",
		SubtleStyle().Render(speciesLabel),
		"",
		line,
		"",
		SubtleStyle().Render(fmt.Sprintf("%d/%d characters", len([]rune(m.RenameBuffer)), maxRenameLen)),
	)
	if m.Status != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", AccentStyle().Render(m.Status))
	}

	return frameScreen(m, "give your plant a name", body, footerRename)
}

func renderSeedOption(m Model, id sim.SpeciesID, idx, width int) string {
	selected := idx == m.SeedChoice
	species, _ := sim.GetSpecies(id)

	border := lipgloss.Color(colorSecondary)
	if selected {
		border = lipgloss.Color(colorAccent)
	}

	innerW := width - 4
	preview := sim.Plant{SpeciesID: id, Stage: sim.StageSeed, Health: 90, Moisture: 55}
	artStr := art.RenderSprite(art.SpriteShelf(preview), 90, 55)

	waterDays := int(species.WaterEvery.Hours() / 24)
	if waterDays < 1 {
		waterDays = 1
	}

	nameStyle := TextStyle()
	if selected {
		nameStyle = nameStyle.Bold(true)
	}

	rows := []string{
		shelfArtBox(artStr, innerW),
		"",
		nameStyle.Render(truncate(species.Name, innerW)),
		SubtleStyle().Render(fmt.Sprintf("water ~%dd", waterDays)),
	}
	if selected {
		rows = append(rows, AccentStyle().Render("▸ selected"))
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(shelfCardHeight).
		Align(lipgloss.Center).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Render(lipgloss.JoinVertical(lipgloss.Center, rows...))
}

func renderDeath(m Model) string {
	p := m.DeathPlant
	if p == nil {
		return ""
	}
	msg := fmt.Sprintf("%s didn't make it.\n\nCause: %s\n\nany key to continue",
		p.Name, m.DeathCause)
	return fullScreenCenter(m, DangerStyle().Render(msg))
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
