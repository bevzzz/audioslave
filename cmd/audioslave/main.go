package main

import (
	"fmt"
	"github.com/bevzzz/audioslave/config"
	"github.com/bevzzz/audioslave/keyboard"
	"github.com/bevzzz/audioslave/volume"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ParseCommand()
	kc := keyboard.NewKeystrokeCounter()

	countStrokes := kc.Count(keyboard.NewDefaultTicker(conf.Interval))
	defer kc.Stop()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("KILLING")
		os.Exit(1)
	}()

	vc := &volume.ItchynyVolumeController{}

	output := volume.NewOutput(conf.Window, conf.Interval, conf.AverageCpm, vc, conf.MinVolume)

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
