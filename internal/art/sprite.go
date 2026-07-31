package art

import "strings"

// Sprite holds plant art plus the metadata needed to color it. Roles are
// derived from glyph shape and row position rather than a parallel mask,
// which keeps art and color permanently in sync.
type Sprite struct {
	Art string
	// PotRows counts trailing rows that make up the pot and its soil.
	PotRows int
	// Body is the role used for foliage glyphs: RoleLeaf for leafy plants,
	// RoleSpike for cacti.
	Body Role
}

const (
	stemGlyphs   = "│"
	flowerGlyphs = "❀✿❁✽●"
	soilGlyphs   = "▒░"
)

func (s Sprite) lines() []string {
	return strings.Split(strings.Trim(s.Art, "\n"), "\n")
}

// droopedLines sags the plant by sinking the foliage toward the soil: the
// topmost rows are shed and replaced by blank rows above. Total height and
// the pot's position stay fixed, so a wilting plant visibly sinks in its
// pot without the whole sprite jumping around between frames.
func (s Sprite) droopedLines(droop float64) []string {
	lines := s.lines()
	sink := int(droop)
	if sink < 1 {
		return lines
	}
	foliage := len(lines) - s.PotRows
	if maxSink := foliage - 1; sink > maxSink {
		sink = maxSink
	}
	if sink < 1 {
		return lines
	}
	out := make([]string, 0, len(lines))
	for i := 0; i < sink; i++ {
		out = append(out, "")
	}
	return append(out, lines[sink:]...)
}

func (s Sprite) roleAt(row, potStart int, ch rune) Role {
	if ch == ' ' {
		return RoleSpace
	}
	if row >= potStart {
		if strings.ContainsRune(soilGlyphs, ch) {
			return RoleSoil
		}
		return RolePot
	}
	if strings.ContainsRune(flowerGlyphs, ch) {
		return RoleFlower
	}
	if strings.ContainsRune(stemGlyphs, ch) {
		return RoleStem
	}
	if s.Body == 0 {
		return RoleLeaf
	}
	return s.Body
}

// Pots. Every plant gets one, sized to the plant it holds.

const potSmall = `╭─────╮
│▒▒▒▒▒│
 ╲___╱`

const potMedium = `╭───────╮
│▒▒▒▒▒▒▒│
 ╲_____╱`

const potLarge = `╭─────────╮
│▒▒▒▒▒▒▒▒▒│
 ╲_______╱`

// potLargeOffset is potLarge nudged one column right, for foliage whose
// crown sits on column 6 rather than column 5.
const potLargeOffset = ` ╭─────────╮
 │▒▒▒▒▒▒▒▒▒│
  ╲_______╱`

func EmptyPot() Sprite {
	return Sprite{
		Art: `
       
       
` + potMedium,
		PotRows: 3,
	}
}
