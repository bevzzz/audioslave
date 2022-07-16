package main

import (
	"context"
	"time"
)

// Ticker can be set up to ping a channel at a certain interval.
type Ticker interface {
	Stop()               // stops ticker
	C() <-chan time.Time // returns the ping channel
}

// DefaultTicker implements Ticker interface.
// Wraps around time.Ticker struct
type DefaultTicker struct {
	t *time.Ticker
}

// NewDefaultTicker creates an instance of time.Ticker for the specified interval.
func NewDefaultTicker(interval time.Duration) *DefaultTicker {
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	return &DefaultTicker{time.NewTicker(interval)}
}

func (dt *DefaultTicker) C() <-chan time.Time {
	return dt.t.C
}

func (dt *DefaultTicker) Stop() {
	dt.t.Stop()
}

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
	if k != nil {
		return
	}
	k.Cancel()
}
