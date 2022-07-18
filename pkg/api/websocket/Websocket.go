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

func (w *Websocket) socketHandler(ctx context.Context, cancel context.CancelFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		conn, err := w.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer func() {
			cancel()
			err := conn.Close()
			log.Println(err)
		}()
		conn.SetPingHandler(nil)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			command := &Command{}
			err := conn.ReadJSON(command)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			err = conn.WriteJSON(&ACK)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				return
			}
			err = command.Decode()
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			resp, err := w.ProcessCommand(*command, cancel)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				continue
			}
			err = conn.WriteJSON(&resp)
			if err != nil {
				log.Println(err)
				conn.WriteJSON(&ERR)
				return
			}
		}
	}
}

func (w Websocket) ProcessCommand(command Command, cancel context.CancelFunc) (interface{}, error) {
	return nil, nil
}

func (w Websocket) Stop() {
	w.Application.Stop()
}
