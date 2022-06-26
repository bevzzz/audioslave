package main

import (
	"fmt"
	"os"
)

func main() {

	kc := NewKeystrokeCounter()
	defer func() {
		kc.Stop()
	}()

	conf, parseMessage, err := parseFlags(os.Args[0], os.Args[1:])
	if wantsHelp(err) {
		fmt.Println(parseMessage)
		os.Exit(2)
	}

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
