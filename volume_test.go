package main

import "testing"

func TestNewVolume(t *testing.T) {
	v := NewVolume(1)
	expectVolumeOn(t, v)
}

func TestAdjustVolume(t *testing.T) {

	t.Run("toggles the volume on and off", func(t *testing.T) {
		v := NewVolume(1)

		v.Adjust(1)
		expectVolumeOff(t, v)

		v.Adjust(1)
		expectVolumeOn(t, v)
	})

	t.Run("does not adjust below sensitivity threshold", func(t *testing.T) {
		v := NewVolume(10)
		v.Adjust(7)
		expectVolumeOn(t, v)
	})
}

func expectVolumeOn(t *testing.T, v *Volume) {
	if !v.On() {
		t.Error("expected the sound to be on")
	}
}

func expectVolumeOff(t *testing.T, v *Volume) {
	if v.On() {
		t.Error("expected the sound to be off")
	}
}
