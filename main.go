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

	vc := &ItchynyVolumeController{}
	volume := NewVolume(10*time.Second, 2*time.Second, vc)

	for {
		n, ok := <-countStrokes
		if !ok {
			fmt.Println("Got interrupted")
			volume.Reset()
			break
		}
		fmt.Printf("You've pressed %d keys in the past %v\n", n, tickInterval)

		volume.Adjust(n)
	}
}

