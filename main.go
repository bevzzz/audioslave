package main

import (
	"fmt"
)

func main() {
	conf := parseCommand()

	kc := NewKeystrokeCounter()
	defer func() {
		kc.Stop()
	}()

	countStrokes := kc.Count(NewDefaultTicker(conf.Interval))

	vc := &ItchynyVolumeController{}

	output := NewOutput(conf.Window, conf.Interval, conf.AverageCpm, vc, conf.MinVolume)

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

// TODO: add command line arguments
// TODO: write README.md and synopsis
// TODO: add test coverage to github actions
