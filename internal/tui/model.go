package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/robbert/plantpal/internal/art"
	"github.com/robbert/plantpal/internal/sim"
	"github.com/robbert/plantpal/internal/store"
)

const tickInterval = time.Second / 12

type animTickMsg time.Time

func Run(initial *sim.GameState, events []sim.Event, away time.Duration, demo bool) error {
	m := NewModel(initial, events, away, demo)
	p := tea.NewProgram(&m)
	_, err := p.Run()
	return err
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg { return animTickMsg(t) })
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case animTickMsg:
		now := time.Time(msg)
		dt := now.Sub(m.LastTick).Seconds()
		if dt <= 0 || dt > 1 {
			dt = 1.0 / 12.0
		}
		m.LastTick = now
		for i := range m.Game.Plants {
			p := &m.Game.Plants[i]
			a, ok := m.Anims[p.ID]
			if !ok {
				a = art.NewPlantAnim()
				m.Anims[p.ID] = a
			}
			if p.IsAlive() {
				a.Update(dt, p.Health)
				if a.Watering && a.WaterTimer <= 0.3 {
					a.OnWaterComplete()
				}
			}
		}
		return m, tea.Tick(tickInterval, func(t time.Time) tea.Msg { return animTickMsg(t) })

	case tea.KeyPressMsg:
		if m.Screen == screenWelcome || m.Screen == screenDeath {
			m.DeathPlant = nil
			m.Screen = screenShelf
			return m, nil
		}

		if m.Screen == screenRename {
			if keyPressed(msg, "ctrl+c", "q") {
				if !m.DemoMode {
					_ = store.Save(m.Game, time.Now())
				}
				m.Quitting = true
				return m, tea.Quit
			}
			return m.handleRenameKey(msg)
		}

		switch {
		case keyPressed(msg, "ctrl+c", "q"):
			if !m.DemoMode {
				_ = store.Save(m.Game, time.Now())
			}
			m.Quitting = true
			return m, tea.Quit

		case keyEsc(msg):
			switch m.Screen {
			case screenDetail:
				m.Screen = screenShelf
				m.Status = ""
			case screenSeedPicker:
				m.cancelSeedPicker()
			}

		case keyPressed(msg, "r", "R"):
			if m.selectedPlant() != nil {
				m.openRename()
			}

		case keyPressed(msg, "down", "j"):
			if m.Screen == screenDetail {
				m.Screen = screenShelf
				m.Status = ""
			}

		case keyLeft(msg):
			if m.Screen == screenSeedPicker {
				m.moveSeedChoice(-1)
			} else {
				m.moveSelection(-1)
			}

		case keyRight(msg):
			if m.Screen == screenSeedPicker {
				m.moveSeedChoice(1)
			} else {
				m.moveSelection(1)
			}

		case keyPressed(msg, "1", "2", "3", "4", "5"):
			if m.Screen == screenSeedPicker {
				idx := int(msg.String()[0]-'0') - 1
				if idx >= 0 && idx < len(m.Game.UnlockedSpecies) {
					m.SeedChoice = idx
				}
			} else {
				m.Selected = int(msg.String()[0]-'0') - 1
				m.Status = ""
			}

		case keyRemove(msg):
			m.tryRemove()

		case keyPlant(msg):
			m.tryPlant()

		case keyPressed(msg, "w"):
			if p := m.currentPlant(); p != nil {
				res := sim.Water(m.Game, p.ID, time.Now())
				m.Status = res.Message
				if res.OK {
					if a, ok := m.Anims[p.ID]; ok {
						a.TriggerWater()
					}
				}
			}

		case keyEnter(msg), keyPressed(msg, "up", "k"):
			switch m.Screen {
			case screenSeedPicker:
				m.confirmSeed()
			default:
				if m.selectedPlant() != nil {
					m.Screen = screenDetail
					m.Status = ""
				} else {
					m.tryPlant()
				}
			}
		}

		if !m.DemoMode {
			m.checkDeaths()
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) tryPlant() {
	if m.DemoMode {
		m.Status = "Planting disabled in demo mode."
		return
	}
	if m.Screen == screenSeedPicker {
		m.confirmSeed()
		return
	}
	m.openSeedPicker(m.Selected)
}

func (m *Model) tryRemove() {
	if m.DemoMode {
		m.Status = "Removing disabled in demo mode."
		return
	}
	m.removeSelected()
}

func (m *Model) openSeedPicker(slot int) {
	if slot < 0 || slot >= m.Game.ShelfCapacity {
		m.Status = "Invalid slot."
		return
	}
	if !m.Game.IsSlotEmpty(slot) {
		m.Status = "Remove the plant first (d)."
		return
	}
	m.PlantingSlot = slot
	m.SeedChoice = 0
	m.Screen = screenSeedPicker
	m.Status = ""
}

func (m *Model) cancelSeedPicker() {
	m.PlantingSlot = -1
	m.Screen = screenShelf
	m.Status = ""
}

func (m *Model) confirmSeed() {
	if m.DemoMode {
		m.Status = "Planting disabled in demo mode."
		return
	}
	if m.PlantingSlot < 0 || len(m.Game.UnlockedSpecies) == 0 {
		return
	}
	if m.SeedChoice < 0 {
		m.SeedChoice = 0
	}
	if m.SeedChoice >= len(m.Game.UnlockedSpecies) {
		m.SeedChoice = len(m.Game.UnlockedSpecies) - 1
	}

	speciesID := m.Game.UnlockedSpecies[m.SeedChoice]
	res := sim.PlantSeed(m.Game, m.PlantingSlot, speciesID, time.Now())
	m.Status = res.Message
	if !res.OK {
		return
	}

	m.Anims[res.PlantID] = art.NewPlantAnim()
	m.Selected = m.PlantingSlot
	m.PlantingSlot = -1
	m.Screen = screenShelf
}

func (m *Model) removeSelected() {
	slot := m.Selected
	if slot < 0 || slot >= m.Game.ShelfCapacity {
		return
	}
	if m.Game.IsSlotEmpty(slot) {
		m.openSeedPicker(slot)
		return
	}

	res := sim.RemovePlant(m.Game, slot)
	m.Status = res.Message
	if res.OK {
		delete(m.Anims, res.PlantID)
		m.openSeedPicker(slot)
	}
}

func (m *Model) moveSeedChoice(delta int) {
	n := len(m.Game.UnlockedSpecies)
	if n == 0 {
		return
	}
	m.SeedChoice = (m.SeedChoice + delta + n) % n
	m.Status = ""
}

func (m *Model) moveSelection(delta int) {
	n := m.Game.ShelfCapacity
	m.Selected = (m.Selected + delta + n) % n
	m.Status = ""
	if m.Screen == screenDetail {
		for i := 0; i < n && m.selectedPlant() == nil; i++ {
			m.Selected = (m.Selected + delta + n) % n
		}
	}
}

func (m *Model) View() tea.View {
	if m.Width == 0 {
		m.Width = 80
	}
	if m.Height == 0 {
		m.Height = 24
	}

	var content string
	switch m.Screen {
	case screenWelcome:
		content = renderWelcome(*m)
	case screenDetail:
		content = renderDetail(*m, m.selectedPlant())
	case screenDeath:
		content = renderDeath(*m)
	case screenSeedPicker:
		content = renderSeedPicker(*m)
	case screenRename:
		content = renderRename(*m)
	default:
		content = renderShelf(*m)
	}

	if m.DemoMode {
		content = DemoBadgeStyle().Render(" demo ") + "\n" + content
	}
	if m.Quitting {
		content += "\n" + SubtleStyle().Render("bye!")
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m *Model) currentPlant() *sim.Plant {
	if m.Screen == screenDetail {
		return m.selectedPlant()
	}
	if p := m.Game.PlantAtSlot(m.Selected); p != nil {
		return p
	}
	return m.Game.FirstAlivePlant()
}

func (m *Model) checkDeaths() {
	for _, p := range m.Game.Plants {
		if p.Stage == sim.StageDead && p.Name != "" && m.DeathPlant == nil {
			cp := p
			m.DeathPlant = &cp
			m.DeathCause = "neglect"
			if len(m.Game.Memorial) > 0 {
				m.DeathCause = m.Game.Memorial[len(m.Game.Memorial)-1].Cause
			}
			m.Screen = screenDeath
			return
		}
	}
}
