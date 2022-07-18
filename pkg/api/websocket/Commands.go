package websocket

import "fmt"

var (
	ACK = Command{
		Type:    "ACK",
		Payload: nil,
	}
	ERR = Command{
		Type:    "ERR",
		Payload: nil,
	}
)

type Command struct {
	Type    string
	Payload interface{}
}

type commandAlg struct {
	Name string
}

func (c *Command) Decode() error {
	switch c.Type {
	case "ALG":
		_, ok := c.Payload.(commandAlg)
		if !ok {
			return fmt.Errorf("could not convert")
		}
	case "STOP":
		// no payload needed
	case "STATUS":
		// no payload needed
	}
	return nil
}
