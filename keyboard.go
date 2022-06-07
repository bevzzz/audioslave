package main

import (
	"time"
)

type Keyboard struct {
	keystrokes chan int
}

func (k *Keyboard) Strokes(interval time.Duration) chan int {
	ch := make(chan int, 2)
	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()

		var counter int
		for {
			select {
			case _, ok := <-k.keystrokes:
				if !ok {
					ch <- counter
				}
				counter++
			case <-tick.C:
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
