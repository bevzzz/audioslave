package main

import (
	"context"
	"fmt"
	"github.com/bevzzz/audioslave"
	"github.com/bevzzz/audioslave/internal/keyboard"
	"github.com/bevzzz/audioslave/internal/util"
	"github.com/bevzzz/audioslave/internal/volume"
	"github.com/bevzzz/audioslave/pkg/algorithms"
	"github.com/bevzzz/audioslave/pkg/api/websocket"
	"github.com/bevzzz/audioslave/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ParseCommand()
	if conf.Verbose {
		log.Printf("Commands parsed:\n%s\n", util.PrettyPrint(conf))
	}
	ctx, cancel := context.WithCancel(context.Background())

	as := &audioslave.AudioSlave{
		KeystrokeCounter: keyboard.NewKeystrokeCounter(),
		VolumeController: &volume.ItchynyVolumeController{},
		Config: &config.Application{
			Config: *conf,
			// Default algs
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
	w := websocket.Websocket{
		Application: as,
		Port:        "10001",
	}
	var err error
	switch conf.Mode {
	case "cli":
		err = as.Start(ctx)
	case "websocket":
		err = w.Start(ctx)
	default:
		log.Fatalf("mode %s not found", conf.Mode)
	}
	if err != nil {
		log.Fatal(err)
	}
}
