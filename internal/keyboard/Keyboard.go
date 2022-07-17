package keyboard

import (
	"context"
)

type KeystrokeCounter interface {
	Count(ticker Ticker) <-chan int
	Stop()
}

// DefaultKeystrokeCounter calculates average typing speed.
// Wraps around `keyboard` package.
type DefaultKeystrokeCounter struct {
	Keylogger KeyLogger
	Cancel    context.CancelFunc
}

// NewKeystrokeCounter creates a DefaultKeystrokeCounter and passes a channel where key-events are posted.
func NewKeystrokeCounter() *DefaultKeystrokeCounter {
	return &DefaultKeystrokeCounter{Keylogger: NewKeyLogger(), Cancel: nil}
}

// Count starts a goroutine that counts keystrokes in per specified interval.
// It returns a channel where the count is posted on every "tick".
//
// Ctrl+C (keyboard.KeyCtrlC) stops the counter and closes the channel.
func (k *DefaultKeystrokeCounter) Count(tick Ticker) <-chan int {
	ch := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	k.Cancel = cancel
	go func() {
		defer tick.Stop()
		counter := 0
		for {
			ks := k.Keylogger.GetKey()
			if !ks.Empty {
				counter++
			}
			select {
			case <-tick.C():
				ch <- counter
				counter = 0
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return ch
}

func (k *DefaultKeystrokeCounter) Stop() {
	if k.Cancel != nil {
		k.Cancel()
	}
}
