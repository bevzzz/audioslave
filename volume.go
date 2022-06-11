package main

import (
	"time"
)

const maxStrokesPerMinute float64 = 200
const minVolume int = 30

type VolumeController interface {
	SetVolume(v int)
	GetVolume() int
}

type Volume struct {
	strokes []float64
	interval time.Duration
	initialVolume int
	controller VolumeController
}

func NewVolume(window, interval time.Duration, vc VolumeController) *Volume {
	strokes := make([]float64, window / interval)
	return &Volume{
		strokes: strokes,
		interval: interval,
		initialVolume: vc.GetVolume(),
		controller: vc,
	}
}

func (v *Volume) Adjust(nStrokes int) {
	nStrokesPerMinute := float64(nStrokes) * float64(time.Minute / v.interval)
	v.strokes = append(v.strokes[1:], nStrokesPerMinute) // push new value
	averageStrokes := mean(v.strokes...)
	reduceVolumeBy := averageStrokes / maxStrokesPerMinute
	newVolume := int(float64(v.initialVolume) * (1 - reduceVolumeBy))
	if newVolume < minVolume {
		newVolume = minVolume
	}

	v.controller.SetVolume(newVolume)
}

func (v *Volume) Reset() {
	v.controller.SetVolume(v.initialVolume)
}

func mean(numbers ...float64) float64 {
	var total float64
	for _, n := range numbers {
		total += n
	}
	return total / float64(len(numbers))
}
