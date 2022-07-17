package config

import (
	"flag"
	"time"
)

// TODO: Algorithm method to use, linear, exponential, logorithm
// TODO: Give the user 5 options on detail grade for the speed of changing

type Config struct {
	Path             string
	MinVolume        int
	MaxVolume        int
	AverageCpm       int
	Interval, Window time.Duration
	Verbose          bool
}

// ParseCommand - parses the command from the cli
func ParseCommand() *Config {
	config := &Config{
		Path:       "config.json",
		MinVolume:  20,
		MaxVolume:  100,
		AverageCpm: 0,
		Interval:   time.Second,
		Window:     10 * time.Second,
		Verbose:    false,
	}

	minVolume := flag.Int("min-volume", config.MinVolume, "set the minimum volume")
	maxVolume := flag.Int("max-volume", config.MaxVolume, "set the maximum volume")
	averageCpm := flag.Int("average-cpm", config.AverageCpm, "set your average typing speed")
	interval := flag.Duration("interval", config.Interval, "change output every N seconds")
	window := flag.Duration("window", config.Window, "change output based on last N values")
	path := flag.String("path", config.Path, "config path to safe and load")
	verbose := flag.Bool("verbose", config.Verbose, "verbose logs")

	flag.Parse()

	config.Path = *path
	config.MaxVolume = *maxVolume
	config.MinVolume = *minVolume
	config.AverageCpm = *averageCpm
	config.Interval = *interval
	config.Window = *window
	config.Verbose = *verbose

	return config
}
