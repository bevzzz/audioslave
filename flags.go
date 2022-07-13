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
func parseCommand() *Config {
	config := &Config{
		MinVolume:  40,
		AverageCpm: 200,
		Interval:   time.Second,
		Window:     10 * time.Second,
	}

	minVolume := flag.Int("min-volume", config.MinVolume, "set the minimum volume")
	averageCpm := flag.Int("average-cpm", config.AverageCpm, "set your average typing speed")
	interval := flag.Duration("interval", config.Interval, "change output every N seconds")
	window := flag.Duration("window", config.Window, "change output based on last N values")

	flag.Parse()

	config.MinVolume = *minVolume
	config.AverageCpm = *averageCpm
	config.Interval = *interval
	config.Window = *window

	return config
}
