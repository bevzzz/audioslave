package volume

import (
	"testing"
	"time"
)

const testMinVolume = 30

func TestAdjust(t *testing.T) {

	t.Run("lower boundary is set", func(t *testing.T) {

		initialVolume := 87
		output, vc := createOutputAndVolumeController(t, initialVolume)

		output.Adjust(50) // would imply volume 0

		if vc.GetVolume() != testMinVolume {
			t.Errorf("got %d, expected minimum %d", vc.GetVolume(), testMinVolume)
		}
	})

	t.Run("sets the volume based on the average char/min", func(t *testing.T) {

		initialVolume := 100
		output, vc := createOutputAndVolumeController(t, initialVolume)

		// Typing speed in each interval
		charactersTyped := []int{11, 7, 1, 9, 7, 0}
		// We expect the volume to decay exponentially as speed increases
		want := []int{96, 92, 92, 83, 71, 88}

		for i, typed := range charactersTyped {

			output.Adjust(typed)

			if vc.GetVolume() != want[i] {
				t.Fatalf(
					"got %d, want %d, given typed %v, and initial volume %d",
					vc.GetVolume(), want[i], charactersTyped[:i+1], initialVolume,
				)
			}
		}
	})
}

func TestReset(t *testing.T) {
	initialVolume := 70
	output, vc := createOutputAndVolumeController(t, initialVolume)

	// Typing speed in each interval
	for _, typed := range []int{3, 3, 6} {
		output.Adjust(typed)
	}

	output.Reset()

	if vc.GetVolume() != initialVolume {
		t.Errorf("got %d, want %d", vc.GetVolume(), initialVolume)
	}
}

// spyVolumeController implements main.VolumeController interface for testing purposes.
type spyVolumeController struct {
	volume int
}

func (s *spyVolumeController) SetVolume(v int) {
	s.volume = v
}

func (s *spyVolumeController) GetVolume() int {
	return s.volume
}

// createOutputAndVolumeController creates spyVolumeController and Output objects with pre-set parameters.
func createOutputAndVolumeController(t testing.TB, initialVolume int) (v *Output, sc *spyVolumeController) {
	t.Helper()

	interval := 2 * time.Second
	window := 10 * time.Second
	averageTypingSpeed := 200

	sc = &spyVolumeController{initialVolume}
	v = NewOutput(window, interval, averageTypingSpeed, sc, testMinVolume)

	return
}
