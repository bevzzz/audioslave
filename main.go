package main

import (
	"fmt"
	"time"
)

const (
	typingSpeedInterval = 2
	typingSpeedWindow = 10
	maxStrokesPerMinute float64 = 200
	minVolume int = 30
)

func main() {

	kb := NewKeyboard()
	defer func() {
		kb.Close()
	}()

	countStrokes := kb.Strokes(NewDefaultTicker(typingSpeedInterval*time.Second))

	vc := &ItchynyVolumeController{}
	volume := NewVolume(typingSpeedWindow*time.Second, typingSpeedInterval*time.Second, vc)

	for {
		n, ok := <-countStrokes
		if !ok {
			fmt.Println("Got interrupted")
			volume.Reset()
			break
		}
		fmt.Printf("You've pressed %d keys in the past %v\n", n, typingSpeedInterval*time.Second)

		volume.Adjust(n)
	}
}

