package config

import "github.com/bevzzz/audioslave/pkg/algorithms"

// Application - application config
type Application struct {
	Config      Config
	ReduceAlg   algorithms.Algorithm
	IncreaseAlg algorithms.Algorithm
}
