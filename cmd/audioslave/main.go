package main

import (
	"context"
	"fmt"
	"github.com/bevzzz/audioslave"
	"github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/algorithms"
	"github.com/bevzzz/audioslave/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ParseCommand()
	ctx, cancel := context.WithCancel(context.Background())

	as := audioslave.AudioSlave{
		KeystrokeCounter: keyboard.NewKeystrokeCounter(),
		VolumeController: &volume.ItchynyVolumeController{},
		Config: config.Application{
			Config: *conf,
			// TODO: determine alg by conf
			ReduceAlg: &algorithms.Linear{
				IncreaseBy: 5,
				ReduceBy:   5,
			},
			IncreaseAlg: &algorithms.Linear{
				IncreaseBy: 5,
				ReduceBy:   5,
			},
		},
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("Finishing application gracefully...")
		as.Stop()
		cancel()
		fmt.Println("Application finished")
		os.Exit(0)
	}()
	log.Println("Starting application...")
	err := as.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
