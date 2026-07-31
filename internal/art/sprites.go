package art

import (
	"fmt"

	"github.com/robbert/plantpal/internal/sim"
)

// SpriteShelf returns the compact sprite used on the shelf.
func SpriteShelf(p sim.Plant) Sprite { return lookup(shelfSprites, p) }

// SpriteDetail returns the large hero sprite used on the detail screen.
func SpriteDetail(p sim.Plant) Sprite { return lookup(detailSprites, p) }

func SpriteFor(p sim.Plant) Sprite { return SpriteShelf(p) }

func lookup(catalog map[string]Sprite, p sim.Plant) Sprite {
	key := fmt.Sprintf("%s_%d", p.SpeciesID, stageKey(p.EffectiveStage()))
	if s, ok := catalog[key]; ok {
		return s
	}
	return catalog[fmt.Sprintf("%s_0", p.SpeciesID)]
}

func stageKey(s sim.Stage) int {
	switch s {
	case sim.StageSeed:
		return 0
	case sim.StageSprout:
		return 1
	case sim.StageYoung:
		return 2
	case sim.StageMature, sim.StageWilting:
		return 3
	case sim.StageBlooming:
		return 4
	default:
		return 0
	}
}

// A freshly sown pot looks the same whatever is buried in it.
func seedSprite(body Role) Sprite {
	return Sprite{Art: `
   ·
` + potSmall, PotRows: 3, Body: body}
}

// ---------------------------------------------------------------------------
// Shelf sprites: small, readable at a glance.
// ---------------------------------------------------------------------------

var shelfSprites = map[string]Sprite{
	// Cactus: a squat column of blocky flesh that sprouts arms.
	"cactus_0": seedSprite(RoleSpike),
	"cactus_1": {Art: `
   ▄
   █
` + potSmall, PotRows: 3, Body: RoleSpike},
	"cactus_2": {Art: `
  ▄█▄
  ███
  ███
` + potSmall, PotRows: 3, Body: RoleSpike},
	"cactus_3": {Art: `
   ▄█▄
 ▄ ███ ▄
 █▄███▄█
   ███
` + potMedium, PotRows: 3, Body: RoleSpike},
	"cactus_4": {Art: `
 ❀ ▄█▄ ❀
 ▄ ███ ▄
 █▄███▄█
   ███
` + potMedium, PotRows: 3, Body: RoleSpike},

	// Spider plant: strap leaves arching out and cascading over the rim.
	"spider_plant_0": seedSprite(RoleLeaf),
	"spider_plant_1": {Art: `
  ╲│╱
   │
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"spider_plant_2": {Art: `
 ╲╲│╱╱
  ╲│╱
   │
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"spider_plant_3": {Art: `
╰╲╲ │ ╱╱╯
  ╲╲│╱╱
   ╲│╱
    │
` + potMedium, PotRows: 3, Body: RoleLeaf},
	"spider_plant_4": {Art: `
❀╲╲ │ ╱╱❀
  ╲╲│╱╱
   ╲│╱
    │
` + potMedium, PotRows: 3, Body: RoleLeaf},

	// Snake plant: stiff vertical swords, never arching.
	"snake_plant_0": seedSprite(RoleLeaf),
	"snake_plant_1": {Art: `
   ║
   ║
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"snake_plant_2": {Art: `
  ║ ║
  ║ ║
  ║║║
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"snake_plant_3": {Art: `
  ║ ║ ║
 ║║ ║ ║║
 ║║ ║ ║║
 ║║║║║║║
` + potMedium, PotRows: 3, Body: RoleLeaf},
	"snake_plant_4": {Art: `
  ❀ ║ ❀
 ║║ ║ ║║
 ║║ ║ ║║
 ║║║║║║║
` + potMedium, PotRows: 3, Body: RoleLeaf},
}

// ---------------------------------------------------------------------------
// Detail sprites: large hero art for the inspect screen.
// ---------------------------------------------------------------------------

var detailSprites = map[string]Sprite{
	"cactus_0": seedSprite(RoleSpike),
	"cactus_1": {Art: `
   ▄
   █
   █
` + potSmall, PotRows: 3, Body: RoleSpike},
	"cactus_2": {Art: `
   ▄█▄
   ███
   ███
   ███
` + potMedium, PotRows: 3, Body: RoleSpike},
	"cactus_3": {Art: `
    ▄█▄
 ▄  ███  ▄
 █  ███  █
 █  ███  █
 █▄▄███▄▄█
    ███
    ███
` + potLarge, PotRows: 3, Body: RoleSpike},
	"cactus_4": {Art: `
  ❀  ❀  ❀
    ▄█▄
 ▄  ███  ▄
 █  ███  █
 █  ███  █
 █▄▄███▄▄█
    ███
    ███
` + potLarge, PotRows: 3, Body: RoleSpike},

	"spider_plant_0": seedSprite(RoleLeaf),
	"spider_plant_1": {Art: `
  ╲│╱
   │
   │
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"spider_plant_2": {Art: `
  ╲╲│╱╱
   ╲│╱
    │
` + potMedium, PotRows: 3, Body: RoleLeaf},
	"spider_plant_3": {Art: `
╰╲╲       ╱╱╯
 ╲╲╲     ╱╱╱
  ╲╲╲   ╱╱╱
   ╲╲╲ ╱╱╱
    ╲╲│╱╱
     ╲│╱
      │
` + potLargeOffset, PotRows: 3, Body: RoleLeaf},
	"spider_plant_4": {Art: `
❀╲╲       ╱╱❀
 ╲╲╲     ╱╱╱
  ╲╲╲   ╱╱╱
   ╲╲╲ ╱╱╱
    ╲╲│╱╱
     ╲│╱
      │
` + potLargeOffset, PotRows: 3, Body: RoleLeaf},

	"snake_plant_0": seedSprite(RoleLeaf),
	"snake_plant_1": {Art: `
   ║
   ║
   ║
` + potSmall, PotRows: 3, Body: RoleLeaf},
	"snake_plant_2": {Art: `
  ║ ║ ║
  ║ ║ ║
  ║║║║║
` + potMedium, PotRows: 3, Body: RoleLeaf},
	"snake_plant_3": {Art: `
   ║  ║  ║
  ║║  ║  ║║
  ║║  ║  ║║
  ║║  ║  ║║
  ║║║ ║ ║║║
` + potLargeOffset, PotRows: 3, Body: RoleLeaf},
	"snake_plant_4": {Art: `
      ❀
   ║  ║  ║
  ║║  ║  ║║
  ║║  ║  ║║
  ║║║ ║ ║║║
` + potLargeOffset, PotRows: 3, Body: RoleLeaf},
}
