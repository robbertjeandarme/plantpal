# Plant Pal

A slim CLI plant tracker — keep an eye on health and water.

## Download a release

Pre-built binaries are on the [Releases](https://github.com/robbertjeandarme/plantpal/releases) page. Pick the asset that matches your OS and CPU:

| Platform | Binary |
|----------|--------|
| macOS (Apple Silicon) | `plantpal-darwin-arm64` |
| macOS (Intel) | `plantpal-darwin-amd64` |
| Linux (x86_64) | `plantpal-linux-amd64` |
| Linux (ARM64) | `plantpal-linux-arm64` |
| Windows (x86_64) | `plantpal-windows-amd64.exe` |
| Windows (ARM64) | `plantpal-windows-arm64.exe` |

### macOS / Linux

```bash
chmod +x plantpal-darwin-arm64   # or whichever binary you downloaded
./plantpal-darwin-arm64
```

On macOS, the first launch may be blocked because the binary is unsigned. Right-click the file in Finder → **Open**, or clear the quarantine flag:

```bash
xattr -d com.apple.quarantine plantpal-darwin-arm64
```

### Windows

Download the `.exe`, then run it from PowerShell or Command Prompt:

```powershell
.\plantpal-windows-amd64.exe
```

Save data is written to `savegame.json` in the directory you run the binary from.

## Build from source

Requires [Go](https://go.dev/dl/) 1.26 or later.

### macOS / Linux

```bash
go build -o plantpal .
./plantpal
```

### Windows

PowerShell or Command Prompt:

```powershell
go build -o plantpal.exe .
.\plantpal.exe
```

Cross-compile all platforms into `dist/`:

```bash
./build-all.sh
```

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
