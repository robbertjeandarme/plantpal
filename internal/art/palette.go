package art

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
)

type Role rune

const (
	RoleSpace  Role = ' '
	RoleLeaf   Role = 'L'
	RoleStem   Role = 'S'
	RolePot    Role = 'P'
	RoleSoil   Role = 'D'
	RoleFlower Role = 'F'
	RoleSpike  Role = 'X'
)

type Palette struct {
	Leaf   colorful.Color
	Stem   colorful.Color
	Pot    colorful.Color
	PotRim colorful.Color
	Soil   colorful.Color
	Flower colorful.Color
	Spike  colorful.Color
}

// PaletteFor builds a high-contrast palette. Hue is in degrees (0-360);
// lightness stays in the 0.38-0.70 band so every element reads clearly
// against a dark terminal background.
func PaletteFor(health, moisture float64) Palette {
	h := clamp01(health / 100)
	m := clamp01(moisture / 100)

	// Foliage shifts green -> yellow -> brown as health drops.
	leaf := colorful.Hsl(
		lerp(28, 128, h),
		lerp(0.55, 0.68, h),
		lerp(0.44, 0.55, h),
	)

	// Cactus flesh is a cooler blue-green so it never reads as a leafy plant.
	spike := colorful.Hsl(
		lerp(32, 158, h),
		lerp(0.45, 0.52, h),
		lerp(0.42, 0.52, h),
	)

	// Wet soil is dark and rich, dry soil is pale and dusty.
	soil := colorful.Hsl(
		lerp(38, 24, m),
		lerp(0.28, 0.48, m),
		lerp(0.56, 0.30, m),
	)

	return Palette{
		Leaf:   leaf,
		Stem:   colorful.Hsl(112, 0.48, lerp(0.34, 0.42, h)),
		Pot:    colorful.Hsl(16, 0.58, 0.47),
		PotRim: colorful.Hsl(14, 0.62, 0.58),
		Soil:   soil,
		Flower: colorful.Hsl(336, 0.78, 0.70),
		Spike:  spike,
	}
}

func (pal Palette) StyleFor(r Role) lipgloss.Style {
	var c colorful.Color
	switch r {
	case RoleStem:
		c = pal.Stem
	case RolePot:
		c = pal.Pot
	case RoleSoil:
		c = pal.Soil
	case RoleFlower:
		c = pal.Flower
	case RoleSpike:
		c = pal.Spike
	case RoleSpace:
		return lipgloss.NewStyle()
	default:
		c = pal.Leaf
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex()))
}

// RenderSprite colors a sprite by deriving each glyph's role from its
// position and shape, so art and color can never drift out of alignment.
func RenderSprite(sp Sprite, health, moisture float64) string {
	return RenderSpriteDroop(sp, health, moisture, 0)
}

// RenderSpriteDroop additionally sags the foliage. Every output line is
// padded to the same visible width so that centering the block in a
// container shifts all rows equally and the art stays aligned.
func RenderSpriteDroop(sp Sprite, health, moisture, droop float64) string {
	pal := PaletteFor(health, moisture)
	lines := sp.droopedLines(droop)
	potStart := len(lines) - sp.PotRows

	// Measure the undrooped art so the block keeps a constant width even
	// when droop sheds the widest row; otherwise a centered plant would
	// jitter sideways between animation frames.
	width := 0
	for _, line := range sp.lines() {
		if n := len([]rune(line)); n > width {
			width = n
		}
	}

	var out strings.Builder
	for y, line := range lines {
		visible := 0
		for _, ch := range line {
			visible++
			role := sp.roleAt(y, potStart, ch)
			if role == RoleSpace {
				out.WriteRune(ch)
				continue
			}
			out.WriteString(pal.StyleFor(role).Render(string(ch)))
		}
		if pad := width - visible; pad > 0 {
			out.WriteString(strings.Repeat(" ", pad))
		}
		out.WriteRune('\n')
	}
	return strings.TrimRight(out.String(), "\n")
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
