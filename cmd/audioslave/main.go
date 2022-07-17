package main

import (
	"fmt"
	keyboard2 "github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ParseCommand()
	kc := keyboard2.NewKeystrokeCounter()

	countStrokes := kc.Count(keyboard2.NewDefaultTicker(conf.Interval))
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
