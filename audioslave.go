package audioslave

import (
	"context"
	"fmt"
	"github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/config"
	"log"
)

type AudioSlave struct {
	KeystrokeCounter keyboard.KeystrokeCounter
	VolumeController volume.VolumeController
	Config           config.Application
}

// Start - starts the audioslave
func (s AudioSlave) Start(ctx context.Context) error {
	countStrokes := s.KeystrokeCounter.Count(keyboard.NewDefaultTicker(s.Config.Config.Interval))
	output := volume.NewOutput(s.Config.Config.Window, s.Config.Config.Interval,
		s.Config.Config.AverageCpm, s.VolumeController, s.Config.Config.MinVolume)
	s.HandleConfig()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
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
func (s AudioSlave) HandleConfig() {
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
func (s AudioSlave) Stop() {
	s.KeystrokeCounter.Stop()

}
