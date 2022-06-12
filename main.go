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

	vc := &SpyVolumeController{100}
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


type SpyVolumeController struct {
	Volume int
}

func (s *SpyVolumeController) SetVolume(v int) {
	fmt.Println("Setting volume", v)
	s.Volume = v
}

func (s *SpyVolumeController) GetVolume() int {
	return s.Volume
}
