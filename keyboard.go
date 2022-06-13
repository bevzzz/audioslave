package main

import (
	"github.com/eiannone/keyboard"
	"time"
)

type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type DefaultTicker struct {
	t *time.Ticker
}

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

type Keyboard struct {
	// TODO: rename to `KeystrokeCounter`
	keystrokes <-chan keyboard.KeyEvent
}

func NewKeyboard() *Keyboard {
	// TODO: check error
	keystrokes, _ := keyboard.GetKeys(10)
	return &Keyboard{keystrokes}
}

func (k *Keyboard) Strokes(tick Ticker) <-chan int {
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
				// TODO: empty k.keystrokes on every tick; measurements from the "last interval" are now affecting the "next" count
				ch <- counter
				counter = 0
			}
		}
	}()
	return ch
}

func (k *Keyboard) Close() {
	// TODO: handle (wrap) error
	keyboard.Close()
}
