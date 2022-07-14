package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := parseCommand()
	kc := NewKeystrokeCounter()

	countStrokes := kc.Count(NewDefaultTicker(conf.Interval))
	defer kc.Stop()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("KILLING")
		os.Exit(1)
	}()

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
