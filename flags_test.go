package main

import (
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
			conf := parseFlags("test", tt.Args)
			if !reflect.DeepEqual(conf, tt.ExpectedConf) {
				t.Errorf("got %+v, want %+v", conf, tt.ExpectedConf)
			}
		})
	}
}
