package audioslave

import (
	"context"
	"fmt"
	"github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/config"
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

// Stop - stops the audioslave
func (s AudioSlave) Stop() {
	s.KeystrokeCounter.Stop()

}
