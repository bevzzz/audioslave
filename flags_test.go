package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	for _, tt := range []struct {
		Args         []string
		ExpectedConf *Config
	}{
		{[]string{""}, &Config{MinVolume: 40}}, // default MinVolume 30
		{[]string{"--min-volume=20"}, &Config{MinVolume: 20}},
	} {
		t.Run(strings.Join(tt.Args, " "), func(t *testing.T) {
			conf := parseFlags("test", tt.Args)
			if !reflect.DeepEqual(conf, tt.ExpectedConf) {
				t.Errorf("got %+v, want %+v", conf, tt.ExpectedConf)
			}
		})
	}
}
