package main

import (
	"bytes"
	"flag"
	"time"
)

type Config struct {
	MinVolume        int
	AverageCpm       int
	Interval, Window time.Duration
}

func parseFlags(programName string, args []string) (*Config, string, error) {
	var conf Config
	var buf bytes.Buffer

	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	flags.SetOutput(&buf)

	flags.IntVar(&conf.MinVolume, "min-volume", 40, "set the minimum volume")
	flags.IntVar(&conf.AverageCpm, "average-cpm", 200, "set your average typing speed")
	flags.DurationVar(&conf.Interval, "interval", time.Second, "change output every N seconds")
	flags.DurationVar(&conf.Window, "window", 10*time.Second, "change output based on last N values")

	err := flags.Parse(args)
	return &conf, buf.String(), err
}

func wantsHelp(err error) bool {
	return err == flag.ErrHelp
}
