package audioslave

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/algorithms"
	"github.com/bevzzz/audioslave/pkg/config"
	"log"
)

type AudioSlave struct {
	KeystrokeCounter keyboard.KeystrokeCounter
	VolumeController volume.VolumeController
	Config           *config.Application
	ReloadConfig     chan bool
	PauseApplication chan bool
}

// Start - starts the audioslave
func (s *AudioSlave) Start(ctx context.Context) error {
	s.HandleConfig()
	s.PauseApplication = make(chan bool)
	s.ReloadConfig = make(chan bool)
	countStrokes := s.KeystrokeCounter.Count(keyboard.NewDefaultTicker(s.Config.Config.Interval))
	output := volume.NewOutput(s.Config.Config.Window, s.Config.Config.Interval,
		s.Config.Config.AverageCpm, s.VolumeController, s.Config.Config.MinVolume)
	for {
		pause := false
		select {
		case <-ctx.Done():
			return nil
		case pause = <-s.PauseApplication:
			// pause application
		case <-s.ReloadConfig:
			// reload config
			s.HandleConfig()
			countStrokes = s.KeystrokeCounter.Count(keyboard.NewDefaultTicker(s.Config.Config.Interval))
			output = volume.NewOutput(s.Config.Config.Window, s.Config.Config.Interval,
				s.Config.Config.AverageCpm, s.VolumeController, s.Config.Config.MinVolume)
		default:
		}
		if !pause {
			n, ok := <-countStrokes
			if !ok {
				output.Reset()
				return fmt.Errorf("channel closed")
			}
			output.Adjust(n)
		}
	}
}

// HandleConfig - handles the reading and saving of the config
func (s *AudioSlave) HandleConfig() {
	err := s.Config.Read()
	if err != nil && s.Config.Config.Verbose {
		log.Println("No config found")
	}
	if s.Config.Config.Path != "" {
		err := s.Config.Write()
		if err != nil && s.Config.Config.Verbose {
			log.Println("Config could not be saved")
		}
	} else if s.Config.Config.Verbose {
		log.Println("No config is saved")
	}
	configData, err := s.Config.ToJson()
	if err == nil && s.Config.Config.Verbose {
		log.Printf("Starting audioslave with config:\n%s\n", string(configData))
	}
}

// Stop - stops the audioslave
func (s *AudioSlave) Stop() {
	s.KeystrokeCounter.Stop()

}

// ChangeAlg - changes algorithm
func (s *AudioSlave) ChangeAlg(name string, data any, increase bool, reduce bool) error {
	newAlgo := algorithms.AlgorithmByName(name)
	dataRaw, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataRaw, &newAlgo)
	if err != nil {
		return err
	}
	if increase {
		s.Config.IncreaseAlg = newAlgo
	}
	if reduce {
		s.Config.ReduceAlg = newAlgo
	}
	s.Config.Write()
	return nil
}

func (s *AudioSlave) Pause() {
	s.PauseApplication <- true
}

func (s *AudioSlave) Resume() {
	s.PauseApplication <- false
}
