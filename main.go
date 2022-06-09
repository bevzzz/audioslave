package main

import (
	"fmt"
	"time"
)

func main() {

	kb := NewKeyboard()
	defer func() {
		_ = kb.Close()
	}()

	const tickInterval = 2 * time.Second
	countStrokes := kb.Strokes(NewDefaultTicker(tickInterval))

	const minKeyStrokesPerInterval = 5
	volume := NewVolume(minKeyStrokesPerInterval)

	for {
		n, ok := <-countStrokes
		if !ok {
			fmt.Println("Got interrupted")
			break
		}
		fmt.Printf("You've pressed %d keys in the past %v -- ", n, tickInterval)

		volume.Adjust(n)

		state := "Off"
		if volume.On() {
			state = "On"
		}
		fmt.Println("Sound", state)
	}
}
