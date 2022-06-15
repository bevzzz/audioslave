package main

import (
	"github.com/eiannone/keyboard"
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

// KeystrokeCounter calculates average typing speed.
// Wraps around `keyboard` package.
type KeystrokeCounter struct {
	keystrokes <-chan keyboard.KeyEvent
}

// NewKeystrokeCounter creates a KeystrokeCounter and passes a channel where key-events are posted.
func NewKeystrokeCounter() *KeystrokeCounter {
	// TODO: check error
	keystrokes, _ := keyboard.GetKeys(10)
	return &KeystrokeCounter{keystrokes}
}

// Count starts a goroutine that counts keystrokes in per specified interval.
// It returns a channel where the count is posted on every "tick".
//
// Ctrl+C (keyboard.KeyCtrlC) stops the counter and closes the channel.
func (k *KeystrokeCounter) Count(tick Ticker) <-chan int {
	ch := make(chan int)
	go func() {
		defer tick.Stop()

		var counter int
		for {
			select {
			case ks, _ := <-k.keystrokes:
				if ks.Key == keyboard.KeyCtrlC {
					close(ch)
					return
				}
				counter++
			case <-tick.C():
				ch <- counter
				counter = 0
			}
		}
	}()
	return ch
}

// Stop closes the key-event channel.
func (k *KeystrokeCounter) Stop() {
	// TODO: handle (wrap) error
	keyboard.Close()
}
