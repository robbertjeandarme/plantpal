# Plant Pal

A slim CLI plant tracker — keep an eye on health and water.

## Run

```bash
go run .
go run . --demo    # preview all growth stages (planting disabled)
```

## Controls

| Key | Action |
|-----|--------|
| `←` `→` | Browse plants |
| `enter` / `p` | Inspect plant, or plant seed in empty slot |
| `w` | Water |
| `r` | Rename plant |
| `d` | Remove plant (opens seed picker) |
| `esc` | Back / cancel |
| `q` | Quit & save |

### Planting

1. Select an **empty slot** (or remove a plant with `d`)
2. Press **`enter`** or **`p`**
3. Pick a seed with `←` `→` or `1` `2` `3`
4. Press **`enter`** to plant

Save: `savegame.json` in the directory you run the game from (project root when using `go run .`)
