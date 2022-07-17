package config

import (
	"encoding/json"
	"github.com/bevzzz/audioslave/pkg/algorithms"
	"os"
)

// Application - application config
type Application struct {
	Config      Config
	ReduceAlg   algorithms.Algorithm
	IncreaseAlg algorithms.Algorithm
}

func (a *Application) Read() error {
	data, err := os.ReadFile(a.Config.Path)
	if err != nil {
		return err
	}
	// unmarshal into generic
	var objMap map[string]*json.RawMessage
	err = json.Unmarshal(data, &objMap)
	if err != nil {
		return err
	}
	// define embed alg struct
	reduceAlgEmbed := struct {
		Type      string
		Algorithm *json.RawMessage
	}{}
	increaseAlgEmbed := struct {
		Type      string
		Algorithm *json.RawMessage
	}{}
	// get embed alg data
	err = json.Unmarshal(*objMap["ReduceAlg"], &reduceAlgEmbed)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objMap["IncreaseAlg"], &increaseAlgEmbed)
	if err != nil {
		return err
	}

	// get alg per name and unmarshal into
	reduceAlg := algorithms.AlgorithmByName(reduceAlgEmbed.Type)
	err = json.Unmarshal(*reduceAlgEmbed.Algorithm, reduceAlg)
	if err != nil {
		return err
	}
	increaseAlg := algorithms.AlgorithmByName(increaseAlgEmbed.Type)
	err = json.Unmarshal(*increaseAlgEmbed.Algorithm, increaseAlg)
	if err != nil {
		return err
	}
	// unmarshal config
	err = json.Unmarshal(*objMap["Config"], &a.Config)
	if err != nil {
		return err
	}
	a.ReduceAlg = reduceAlg
	a.IncreaseAlg = increaseAlg
	return nil
}

func (a *Application) Write() error {
	data, err := a.ToJson()
	if err != nil {
		return err
	}
	return os.WriteFile(a.Config.Path, data, 0644)
}

func (a *Application) ToJson() ([]byte, error) {
	// embed alg
	applicationEmbed := &Application{
		Config: a.Config,
		ReduceAlg: struct {
			algorithms.Algorithm
			Type string
		}{a.ReduceAlg, a.ReduceAlg.Name()},
		IncreaseAlg: struct {
			algorithms.Algorithm
			Type string
		}{a.IncreaseAlg, a.IncreaseAlg.Name()},
	}

	return json.MarshalIndent(applicationEmbed, "", " ")
}
