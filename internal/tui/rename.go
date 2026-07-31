package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/robbert/plantpal/internal/sim"
)

const maxRenameLen = 24

func (m *Model) handleRenameKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case keyEsc(msg):
		m.cancelRename()

	case keyEnter(msg):
		m.confirmRename()

	case keyRemove(msg):
		m.renameBackspace()

	default:
		m.renameAppend(msg)
	}
	return m, nil
}

func (m *Model) openRename() {
	p := m.currentPlant()
	if p == nil {
		m.Status = "No plant selected."
		return
	}
	m.RenamePlantID = p.ID
	m.RenameBuffer = p.Name
	m.RenameFrom = m.Screen
	m.Screen = screenRename
	m.Status = ""
}

func (m *Model) cancelRename() {
	returnTo := m.RenameFrom
	if returnTo == 0 {
		returnTo = screenDetail
	}
	m.RenameBuffer = ""
	m.RenamePlantID = ""
	m.RenameFrom = 0
	m.Screen = returnTo
	m.Status = ""
}

func (m *Model) confirmRename() {
	if m.RenamePlantID == "" {
		m.cancelRename()
		return
	}
	res := sim.RenamePlant(m.Game, m.RenamePlantID, m.RenameBuffer)
	returnTo := m.RenameFrom
	if returnTo == 0 {
		returnTo = screenDetail
	}
	m.RenameBuffer = ""
	m.RenamePlantID = ""
	m.RenameFrom = 0
	m.Screen = returnTo
	m.Status = res.Message
}

func (m *Model) renameBackspace() {
	r := []rune(m.RenameBuffer)
	if len(r) == 0 {
		return
	}
	m.RenameBuffer = string(r[:len(r)-1])
	m.Status = ""
}

func (m *Model) renameAppend(msg tea.KeyPressMsg) {
	ch := typedRune(msg)
	if ch == 0 {
		return
	}
	if len([]rune(m.RenameBuffer)) >= maxRenameLen {
		return
	}
	m.RenameBuffer += string(ch)
	m.Status = ""
}

func typedRune(msg tea.KeyPressMsg) rune {
	k := msg.Key()
	if k.Text != "" {
		for _, r := range k.Text {
			if r >= 32 && r != 127 {
				return r
			}
		}
	}
	if k.Code >= 32 && k.Code < 127 {
		return k.Code
	}
	return 0
}
