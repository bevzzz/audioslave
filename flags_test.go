package main

import (
	"flag"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseFlags(t *testing.T) {
	for _, tt := range []struct {
		Args         []string
		ExpectedConf *Config
	}{
		{[]string{""}, &Config{
			MinVolume:  40,
			AverageCpm: 200,
			Interval:   time.Second,
			Window:     10 * time.Second,
		}}, // default MinVolume 30
		{
			[]string{
				"--min-volume=20",
				"--average-cpm", "100",
				"--interval=5s",
				"--window=60s",
			},
			&Config{
				MinVolume:  20,
				AverageCpm: 100,
				Interval:   5 * time.Second,
				Window:     60 * time.Second,
			}},
	} {
		t.Run(strings.Join(tt.Args, " "), func(t *testing.T) {
			conf, _, _ := parseFlags("test", tt.Args)
			if !reflect.DeepEqual(conf, tt.ExpectedConf) {
				t.Errorf("got %+v, want %+v", conf, tt.ExpectedConf)
			}
		})
	}

	t.Run("--help flag returns default dump", func(t *testing.T) {
		_, output, err := parseFlags("helpMe", []string{"--help"})
		if output == "" {
			t.Errorf("expected default help text, but didn't get one")
		}
		if err == nil || !wantsHelp(err) {
			t.Errorf("must return %q if --help flag is not defined", flag.ErrHelp)
		}
	})
}
