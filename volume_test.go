package main

import (
	"testing"
	"time"
)

func TestAdjust(t *testing.T) {

	t.Run("lower boundary is set", func(t *testing.T) {

		initialVolume := 87
		volume, controller := createVolumeAndController(t, initialVolume)

		volume.Adjust(50) // would imply volume 0

		if controller.Volume < minVolume {
			t.Errorf("got %d, expected minimum %d", controller.Volume, minVolume)
		}
	})

	t.Run("sets the volume based on the average char/min", func(t *testing.T) {

		initialVolume := 100
		volume, controller := createVolumeAndController(t, initialVolume)

		// Typing speed in each interval
		charactersTyped := []int{11, 7, 1, 9, 7, 0}
		// We expect the volume to decay exponentially as speed increases
		want := []int{96, 92, 92, 83, 71, 88}

		for i, typed := range charactersTyped {

			volume.Adjust(typed)

			if controller.Volume != want[i] {
				t.Fatalf(
					"got %d, want %d, given typed %v, and initial volume %d",
					controller.Volume, want[i], charactersTyped[:i+1], initialVolume,
				)
			}
		}
	})
}

func TestReset(t *testing.T) {
	initialVolume := 70
	volume, controller := createVolumeAndController(t, initialVolume)

	// Typing speed in each interval
	for _, typed := range []int{3, 3, 6} {
		volume.Adjust(typed)
	}

	volume.Reset()

	if controller.Volume != initialVolume {
		t.Errorf("got %d, want %d", controller.Volume, initialVolume)
	}
}

type spyVolumeController struct {
	Volume int
}

func (s *spyVolumeController) SetVolume(v int) {
	s.Volume = v
}

func (s *spyVolumeController) GetVolume() int {
	return s.Volume
}

func createVolumeAndController(t testing.TB, initialVolume int) (v *Volume, sc *spyVolumeController) {
	t.Helper()

	interval := 2 * time.Second
	window := 10 * time.Second

	sc = &spyVolumeController{initialVolume}
	v = NewVolume(window, interval, 200, sc)

	return
}
