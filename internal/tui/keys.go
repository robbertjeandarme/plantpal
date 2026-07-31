package tui

import tea "charm.land/bubbletea/v2"

func keyPressed(msg tea.KeyPressMsg, keys ...string) bool {
	s := msg.String()
	for _, k := range keys {
		if s == k {
			return true
		}
	}
	return false
}

func keyEnter(msg tea.KeyPressMsg) bool {
	k := msg.Key()
	return k.Code == tea.KeyEnter || keyPressed(msg, "enter", "return")
}

func keyEsc(msg tea.KeyPressMsg) bool {
	k := msg.Key()
	return k.Code == tea.KeyEscape || k.Code == tea.KeyEsc || keyPressed(msg, "esc", "escape")
}

func keyLeft(msg tea.KeyPressMsg) bool {
	k := msg.Key()
	return k.Code == tea.KeyLeft || keyPressed(msg, "left", "h")
}

func keyRight(msg tea.KeyPressMsg) bool {
	k := msg.Key()
	return k.Code == tea.KeyRight || keyPressed(msg, "right", "l")
}

func keyRemove(msg tea.KeyPressMsg) bool {
	return keyPressed(msg, "d", "D", "delete", "backspace", "backspace2")
}

func keyPlant(msg tea.KeyPressMsg) bool {
	return keyPressed(msg, "p", "P")
}
