package websocket

import (
	"context"
	"fmt"
	"github.com/bevzzz/audioslave"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Websocket struct {
	Application *audioslave.AudioSlave
	Port        string
	Upgrader    websocket.Upgrader
}

// Start - starts the websocket API and the application
func (w *Websocket) Start(ctx context.Context) error {
	w.Upgrader = websocket.Upgrader{} // use default options
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		server := &http.Server{Addr: fmt.Sprintf("localhost:%s", w.Port),
			Handler: w.socketHandler(ctx, cancel)}
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
			w.Application.Stop()
		}
	}()
	err := w.Application.Start(ctx)
	if err != nil {
		return err
	}
	return nil
}

// socketHandler - handler func for incoming request
func (w *Websocket) socketHandler(ctx context.Context, cancel context.CancelFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		conn, err := w.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer func() {
			err := conn.Close()
			log.Println(err)
			cancel()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			command := &Command{}
			// read json
			err := conn.ReadJSON(command)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			// write ack
			err = conn.WriteJSON(&ACK)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				return
			}
			// decode command
			err = command.Decode()
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			// process command
			resp, err := w.ProcessCommand(*command, cancel)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			// write resp. ack or data structure
			err = conn.WriteJSON(&resp)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				return
			}
		}
	}
}

// ProcessCommand - Process a command and returns a response
func (w *Websocket) ProcessCommand(command Command, cancel context.CancelFunc) (any, error) {
	switch command.Type {
	case STOPCommand:
		w.Stop()
		cancel()
	case ALGCommand:
		payload := command.Payload.(commandAlg)
		err := w.Application.ChangeAlg(payload.Type, payload.Data, payload.Increase, payload.Reduce)
		if err != nil {
			return nil, err
		}
	case PAUSECommand:
		w.Application.Pause()
	case RESUMECommand:
		w.Application.Resume()
	case RELOADCONFIGCommand:
		w.Application.ReloadConfig()
	default:
	}
	return &ACK, nil
}

// Stop - stops the application underneath
func (w *Websocket) Stop() {
	w.Application.Stop()
}
