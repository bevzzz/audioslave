package main

import (
	"fmt"
	"time"

	itchyny "github.com/itchyny/volume-go"
)


type VolumeController interface {
	GetVolume() int
	SetVolume(v int)
}

type ItchynyVolumeController struct {}

func (vc *ItchynyVolumeController) GetVolume() int {
	// TODO: check error
	v, _ := itchyny.GetVolume()
	fmt.Println("Current volume is", v)
	return v
}

func (vc *ItchynyVolumeController) SetVolume(v int) {
	fmt.Println("Setting volume to", v)
	itchyny.SetVolume(v)
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

// TODO: the volume should increase faster then it decreases
// TODO: volume should not increase/decrease linearly: the volume should change faster as the speed approaches `maxStrokesPerMinute`

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
