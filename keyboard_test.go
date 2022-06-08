package main

import (
	"reflect"
	"testing"
	"time"
)

func TestStrokes(t *testing.T) {

	t.Run("number of keystrokes is captured", func(t *testing.T) {
		var got []int
		want := []int{2, 5, 6}

		// Create a fake channel through which keystrokes will be sent
		keyChan := make(chan int, 6)
		kb := Keyboard{keystrokes: keyChan}

		strokeCount := kb.Strokes(NewDefaultTicker(0*time.Millisecond))

		// Simulate keystrokes being sent through the channel
		for _, n := range want {
			for i := 0; i < n; i++ {
				keyChan <- n
			}
			// Strokes() should post count every 1s
			got = append(got, <-strokeCount)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("sends intervals at a specified rate", func(t *testing.T) {
		// Create a fake channel through which keystrokes will be sent
		keyChan := make(chan int, 6)
		kb := Keyboard{keystrokes: keyChan}

		spyTicker := NewSpyTicker(10 * time.Millisecond)
		kb.Strokes(spyTicker)
		time.Sleep(time.Millisecond)
		spyTicker.Stop()

		if spyTicker.Calls == 0 {
			t.Error("expected calls to ticker, but did not get any")
		}
	})
}

type SpyTicker struct {
	interval time.Duration
	Calls int
	c <-chan time.Time
}

func NewSpyTicker(interval time.Duration) *SpyTicker {
	c := make(chan time.Time, 50)
	for i := 0; i < 50; i++ {
		c <- time.Time{}
	}

	return &SpyTicker{
		interval: interval,
		c: c,
	}
}

func (s *SpyTicker) C() <-chan time.Time {
	return s.c
}

func (s *SpyTicker) Stop() {
	s.Calls = 50 - len(s.c)
}