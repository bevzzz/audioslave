package main

import (
	"flag"
	"time"
)

type Config struct {
	MinVolume        int
	AverageCpm       int
	Interval, Window time.Duration
}

// parseCommand - parses the command from the cli
func parseCommand() Config {
	config := &Config{
		MinVolume:  40,
		AverageCpm: 200,
		Interval:   time.Second,
		Window:     10 * time.Second,
	}

	config.MinVolume = *flag.Int("min-volume", config.MinVolume, "set the minimum volume")
	config.AverageCpm = *flag.Int("average-cpm", config.AverageCpm, "set your average typing speed")
	config.Interval = *flag.Duration("interval", config.Interval, "change output every N seconds")
	config.Window = *flag.Duration("window", config.Window, "change output based on last N values")

	flag.Parse()
	return *config
}
