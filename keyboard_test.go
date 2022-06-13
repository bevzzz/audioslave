package main

import (
	"github.com/eiannone/keyboard"
	"reflect"
	"testing"
	"time"
)

func TestStrokes(t *testing.T) {

	t.Run("number of keystrokes is captured", func(t *testing.T) {
		var got []int
		want := []int{2, 5, 6}

		// Create a fake channel through which keystrokes will be sent
		_, keyChan, strokeCount := createKeyboardWithFakeChannel(t)

		// Simulate keystrokes being sent through the channel
		for _, n := range want {
			for i := 0; i < n; i++ {
				keyChan <- keyboard.KeyEvent{}
			}
			// Strokes() should post count every 1s
			got = append(got, <-strokeCount)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("sends intervals at a specified rate", func(t *testing.T) {
		kb := NewKeyboard()
		defer func() {
			kb.Close()
		}()

		spyTicker := NewSpyTicker(5 * time.Millisecond)
		kb.Strokes(spyTicker)
		time.Sleep(time.Millisecond)
		spyTicker.Stop()

		if spyTicker.Calls == 0 {
			t.Error("expected calls to ticker, but did not get any")
		}
	})

	t.Run("count channel is closed after Ctrl+C", func(t *testing.T) {
		_, keyChan, strokeCount := createKeyboardWithFakeChannel(t)

		// Imitate sending 'interrupt event'
		keyChan <- keyboard.KeyEvent{Key: keyboard.KeyCtrlC}

		// Interrupt signal (-1) is sent
		_, open := <-strokeCount
		if open {
			t.Error("the channel is open, expected closed")
		}
	})

	t.Run("goroutine does not send on closed channel", func(t *testing.T) {
		kb, keyChan, _ := createKeyboardWithFakeChannel(t)

		// Imitate sending 'interrupt event'
		keyChan <- keyboard.KeyEvent{Key: keyboard.KeyCtrlC}

		// This should not cause panic
		kb.Close()
	})

	t.Run("0 strokes sent through the channel", func(t *testing.T) {
		kb := NewKeyboard()
		defer func() {
			kb.Close()
		}()

		// Start stroke count
		strokeCount := kb.Strokes(NewSpyTicker(0 * time.Millisecond))

		want := 0
		for i := 0; i < 3; i++ {
			select {
			case got, _ := <-strokeCount:
				if got != want {
					t.Fatalf("got %q, want %q", got, want)
				}
			case <-time.After(1*time.Millisecond):
				t.Fatalf("expected a value in the strokeCount channel")
			}
		}
	})
}

func createKeyboardWithFakeChannel(t testing.TB) (Keyboard, chan keyboard.KeyEvent, <-chan int) {
	t.Helper()

	// Create a fake channel through which keystrokes can be sent
	keyChan := make(chan keyboard.KeyEvent)
	kb := Keyboard{keystrokes: keyChan}

	// Start stroke count
	strokeCount := kb.Strokes(NewDefaultTicker(0 * time.Millisecond))

	return kb, keyChan, strokeCount
}

type SpyTicker struct {
	interval time.Duration
	Calls    int
	c        <-chan time.Time
}

func NewSpyTicker(interval time.Duration) *SpyTicker {
	c := make(chan time.Time, 50)
	for i := 0; i < 50; i++ {
		c <- time.Time{}
	}

	return &SpyTicker{
		interval: interval,
		c:        c,
	}
}

func (s *SpyTicker) C() <-chan time.Time {
	return s.c
}

func (s *SpyTicker) Stop() {
	s.Calls = 50 - len(s.c)
}
//
//func (s *SpyTicker) Tick() {
//	interface{}(s.c).(chan time.Time) <- time.Now()
//}