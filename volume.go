package main

import (
	"math"
	"time"

	itchyny "github.com/itchyny/volume-go"
)

// VolumeController can read and modify the output volume.
type VolumeController interface {
	GetVolume() int
	SetVolume(v int)
}

// ItchynyVolumeController implements VolumeController.
// Wraps around `itchyny` package for manipulating volume on the computer.
type ItchynyVolumeController struct{}

// GetVolume reads the current volume.
func (vc *ItchynyVolumeController) GetVolume() int {
	// TODO: check error
	v, _ := itchyny.GetVolume()
	return v
}

// SetVolume sets the volume.
func (vc *ItchynyVolumeController) SetVolume(v int) {
	itchyny.SetVolume(v)
}

// Output holds the information necessary for adjusting
// the output levels based on the user's typing speed.
type Output struct {
	strokes       []float64
	interval      time.Duration
	initialVolume int
	controller    VolumeController
	reduceBy      func(float64) float64
}

// NewOutput creates an Output object with the appropriate strokes buffer
// and a function for calculating the level of output reduction.
func NewOutput(window, interval time.Duration, averageCPM float64, vc VolumeController) *Output {
	strokes := make([]float64, window/interval)
	r := getExponentialDecayFunc(averageCPM)
	return &Output{
		strokes:       strokes,
		interval:      interval,
		initialVolume: vc.GetVolume(),
		controller:    vc,
		reduceBy:      r,
	}
}

// Adjust recalculates the average typing speed and adjusts the output accordingly.
func (o *Output) Adjust(nStrokes int) {
	nStrokesPerMinute := float64(nStrokes) * float64(time.Minute/o.interval)
	o.strokes = append(o.strokes[1:], nStrokesPerMinute) // push new value

	averageStrokes := mean(o.strokes...)
	reduceVolumeBy := o.reduceBy(averageStrokes) / 100
	newVolume := int(float64(o.initialVolume) * (1 - reduceVolumeBy))

	if newVolume < minVolume {
		newVolume = minVolume
	}

	o.controller.SetVolume(newVolume)
}

// Reset restores the original output level.
func (o *Output) Reset() {
	o.controller.SetVolume(o.initialVolume)
}

// mean calculates the average value of its arguments
func mean(numbers ...float64) float64 {
	var total float64
	for _, n := range numbers {
		total += n
	}
	return total / float64(len(numbers))
}

// getExponentialDecayFunc derives a formula for the exponential equation from
// 2 points at which we want the output level to drop by 10% and 100% respectively.
// It then creates a function that can estimate the result for any x.
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

	return func(x float64) float64 {
		if x == 0 {
			return 0
		}
		return a * math.Pow(b, x)
	}
}
