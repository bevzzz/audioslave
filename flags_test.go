package main

import (
	"flag"
	"os"
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
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = []string{"cmd"}
			for _, arg := range tt.Args {
				os.Args = append(os.Args, arg)
			}
			conf := parseCommand()
			if !reflect.DeepEqual(conf, tt.ExpectedConf) {
				t.Errorf("got %+v, want %+v", conf, tt.ExpectedConf)
			}
		})
	}
}
