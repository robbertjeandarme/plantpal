package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/robbert/plantpal/internal/sim"
	"github.com/robbert/plantpal/internal/store"
	"github.com/robbert/plantpal/internal/tui"
)

func main() {
	timeScale := flag.Float64("time-scale", 1.0, "simulation speed multiplier (1440 = 1 real minute equals 1 sim day)")
	seed := flag.Int64("seed", 0, "random seed for new games")
	reset := flag.Bool("reset", false, "delete save and start fresh")
	demo := flag.Bool("demo", false, "UI preview with showcase plants at every growth stage")
	flag.Parse()

	now := time.Now()

	if *demo {
		game := sim.DemoGarden(now)
		if err := tui.Run(game, nil, 0, true); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *reset {
		if err := store.DeleteSave(); err != nil {
			fmt.Fprintf(os.Stderr, "reset failed: %v\n", err)
			os.Exit(1)
		}
	}

	game, existed, err := store.Load(now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load failed: %v\n", err)
		os.Exit(1)
	}
	if !existed && *seed != 0 {
		game = sim.NewGameState(now, *seed)
	}
	if *timeScale > 0 {
		game.TimeScale = *timeScale
	}

	away := now.Sub(game.LastPlayedAt)
	events := sim.Advance(game, now)

	if err := tui.Run(game, events, away, false); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := store.Save(game, time.Now()); err != nil {
		fmt.Fprintf(os.Stderr, "save failed: %v\n", err)
		os.Exit(1)
	}
}
