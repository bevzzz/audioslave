package main

import (
	"github.com/eiannone/keyboard"
	"time"
)

type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type DefaultTicker struct{
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
	keystrokes <-chan keyboard.KeyEvent
}

func NewKeyboard() *Keyboard {
	keystrokes, _ := keyboard.GetKeys(10) // TODO: say, do I need to test it?
	return &Keyboard{keystrokes}
}

func (k *Keyboard) Strokes(tick Ticker) <-chan int {
	ch := make(chan int)
	go func() {
		defer tick.Stop()

		var counter int
		for {
			select {
			case _, ok := <-k.keystrokes:
				if !ok {
					ch <- counter
				}
				counter++
			case <-tick.C():
				if counter == 0 {
					continue
				}
				ch <- counter
				counter = 0
			}
		}
	}()
	return ch
}

func (k *Keyboard) Close() error {
	err := keyboard.Close()
	return err // TODO: wrap error
}