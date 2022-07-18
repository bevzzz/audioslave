package websocket

import "fmt"

type CommandType string

var (
	ACKCommand  CommandType = "ACK"
	ERRCommand  CommandType = "ERR"
	ALGCommand  CommandType = "ALG"
	STOPCommand CommandType = "STOP"
)

var (
	ACK = Command{
		Type:    ACKCommand,
		Payload: nil,
	}
	ERR = Command{
		Type:    ERRCommand,
		Payload: nil,
	}
)

type Command struct {
	Type    CommandType
	Payload any
}

type commandAlg struct {
	Type     string
	Reduce   bool
	Increase bool
	Data     any
}

// Decode - decodes the command
func (c *Command) Decode() error {
	switch c.Type {
	case ALGCommand:
		_, ok := c.Payload.(commandAlg)
		if !ok {
			return fmt.Errorf("could not convert")
		}
	case STOPCommand:
		// no payload needed
	}
	return nil
}
