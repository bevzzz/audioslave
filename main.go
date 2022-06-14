package main

import (
	"fmt"
	"time"
)

const (
	typingSpeedInterval = 1
	typingSpeedWindow = 10
	minVolume = 50
)


func main() {

	kb := NewKeyboard()
	defer func() {
		kb.Close()
	}()

	countStrokes := kb.Strokes(NewDefaultTicker(typingSpeedInterval*time.Second))

	averageStrokesPerMinute := 300.0
	vc := &ItchynyVolumeController{}
	volume := NewVolume(typingSpeedWindow*time.Second, typingSpeedInterval*time.Second, averageStrokesPerMinute, vc)

	for {
		n, ok := <-countStrokes
		if !ok {
			fmt.Println("Got interrupted")
			volume.Reset()
			break
		}
		//fmt.Printf("You've pressed %d keys in the past %v\n", n, typingSpeedInterval*time.Second)

		volume.Adjust(n)
	}
}

