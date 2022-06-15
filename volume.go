package main

import (
	"math"
	"time"

	itchyny "github.com/itchyny/volume-go"
)

type VolumeController interface {
	GetVolume() int
	SetVolume(v int)
}

type ItchynyVolumeController struct{}

func (vc *ItchynyVolumeController) GetVolume() int {
	// TODO: check error
	v, _ := itchyny.GetVolume()
	return v
}

func (vc *ItchynyVolumeController) SetVolume(v int) {
	itchyny.SetVolume(v)
}

type Volume struct {
	strokes       []float64
	interval      time.Duration
	initialVolume int
	controller    VolumeController
	reduceBy      func(float64) float64
}

func NewVolume(window, interval time.Duration, averageCPM float64, vc VolumeController) *Volume {
	strokes := make([]float64, window/interval)
	r := getExponentialDecayFunc(averageCPM)
	return &Volume{
		strokes:       strokes,
		interval:      interval,
		initialVolume: vc.GetVolume(),
		controller:    vc,
		reduceBy:      r,
	}
}

func (v *Volume) Adjust(nStrokes int) {
	nStrokesPerMinute := float64(nStrokes) * float64(time.Minute/v.interval)
	v.strokes = append(v.strokes[1:], nStrokesPerMinute) // push new value

	averageStrokes := mean(v.strokes...)
	reduceVolumeBy := v.reduceBy(averageStrokes) / 100
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

func getExponentialDecayFunc(averageCPM float64) func(float64) float64 {
	type point struct {
		x float64
		y float64
	}

	p1 := point{x: averageCPM / 1.5, y: 10}  // point at which volume should be reduced by 10%
	p2 := point{x: averageCPM * 1.5, y: 100} // point at which volume should be reduced by 100%

	// Solve f(x) = ab^x for a and b
	b := math.Pow(p2.y/p1.y, 1/(p2.x-p1.x))
	a := p1.y * math.Pow(b, -p1.x)

	return func(i float64) float64 {
		if i == 0 {
			return 0
		}
		return a * math.Pow(b, i)
	}
}
