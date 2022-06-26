package main

import "flag"

type Config struct {
	MinVolume int
}

func parseFlags(programName string, args []string) *Config {
	var conf Config
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	flags.IntVar(&conf.MinVolume, "min-volume", 40, "volume to set at maximum typing speed")
	flags.Parse(args)
	return &conf
}
