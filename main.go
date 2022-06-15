package main

import (
	"fmt"
	"time"
)

const (
	typingSpeedInterval = 1  // average typing speed will be recalculated every N seconds
	typingSpeedWindow   = 10 // calculate average over the past M seconds
	minVolume           = 50
)

func main() {

	kc := NewKeystrokeCounter()
	defer func() {
		kc.Stop()
	}()

	countStrokes := kc.Count(NewDefaultTicker(typingSpeedInterval * time.Second))

	// TODO: make averageStrokesPerMinute a command-line option
	averageStrokesPerMinute := 300.0
	vc := &ItchynyVolumeController{}

	// TODO: consider hiding "time.Seconds" behind Output
	output := NewOutput(typingSpeedWindow*time.Second, typingSpeedInterval*time.Second, averageStrokesPerMinute, vc)

	for {
		n, ok := <-countStrokes
		if !ok {
			fmt.Println("Got interrupted")
			output.Reset()
			break
		}
		output.Adjust(n)
	}
}
