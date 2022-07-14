package main

import (
	"context"
	"github.com/eiannone/keyboard"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func fakeKeyboardStroker() (<-chan keyboard.KeyEvent, context.CancelFunc) {
	ch := make(chan keyboard.KeyEvent)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for true {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- randomKey()
			}
		}
	}()
	return ch, cancel
}

func randomKey() keyboard.KeyEvent {
	keys := []keyboard.Key{keyboard.KeyCtrlTilde,
		keyboard.KeyCtrl2,
		keyboard.KeyCtrlSpace,
		keyboard.KeyCtrlA,
		keyboard.KeyCtrlB,
		keyboard.KeyCtrlC,
		keyboard.KeyCtrlD,
		keyboard.KeyCtrlE,
		keyboard.KeyCtrlF,
		keyboard.KeyCtrlG,
		keyboard.KeyBackspace,
		keyboard.KeyCtrlH,
		keyboard.KeyTab,
		keyboard.KeyCtrlI,
		keyboard.KeyCtrlJ,
		keyboard.KeyCtrlK,
		keyboard.KeyCtrlL,
		keyboard.KeyEnter,
		keyboard.KeyCtrlM,
		keyboard.KeyCtrlN,
		keyboard.KeyCtrlO,
		keyboard.KeyCtrlP,
		keyboard.KeyCtrlQ,
		keyboard.KeyCtrlR,
		keyboard.KeyCtrlS,
		keyboard.KeyCtrlT,
		keyboard.KeyCtrlU,
		keyboard.KeyCtrlV,
		keyboard.KeyCtrlW,
		keyboard.KeyCtrlX,
		keyboard.KeyCtrlY,
		keyboard.KeyCtrlZ,
		keyboard.KeyEsc,
		keyboard.KeyCtrlLsqBracket,
		keyboard.KeyCtrl3,
		keyboard.KeyCtrl4,
		keyboard.KeyCtrlBackslash,
		keyboard.KeyCtrl5,
		keyboard.KeyCtrlRsqBracket,
		keyboard.KeyCtrl6,
		keyboard.KeyCtrl7,
		keyboard.KeyCtrlSlash,
		keyboard.KeyCtrlUnderscore,
		keyboard.KeySpace,
		keyboard.KeyBackspace2,
		keyboard.KeyCtrl8}
	length := len(keys)
	i := rand.Intn(length)
	return keyboard.KeyEvent{
		Key:  keys[i],
		Rune: ' ',
		Err:  nil,
	}
}

func TestCount(t *testing.T) {

	t.Run("number of keystrokes is captured", func(t *testing.T) {
		var got []int
		want := []int{2, 5, 6}

		_, keyChan, strokeCount := createCounterWithFakeChannel(t)

		// Simulate keystrokes being sent through the channel
		for _, n := range want {
			for i := 0; i < n; i++ {
				keyChan <- keyboard.KeyEvent{}
			}
			// Count() should post count at an interval
			got = append(got, <-strokeCount)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("uses ticket to send updates at a certain interval", func(t *testing.T) {
		ch, cancelF := fakeKeyboardStroker()
		kc := DefaultKeystrokeCounter{keystrokes: ch}
		defer cancelF()

		spyTicker := newSpyTicker(5 * time.Millisecond)
		kc.Count(spyTicker)
		time.Sleep(time.Millisecond)
		spyTicker.Stop()

		if spyTicker.Calls == 0 {
			t.Error("expected calls to ticker, but did not get any")
		}
	})

	t.Run("count channel is closed after Ctrl+C", func(t *testing.T) {
		_, keyChan, strokeCount := createCounterWithFakeChannel(t)

		// Imitate sending 'interrupt event'
		keyChan <- keyboard.KeyEvent{Key: keyboard.KeyCtrlC}

		// Interrupt signal (-1) is sent
		if _, open := <-strokeCount; open {
			t.Error("the channel is open, expected closed")
		}
	})

	t.Run("goroutine does not send on closed channel", func(t *testing.T) {
		kc, keyChan, _ := createCounterWithFakeChannel(t)

		// Imitate sending 'interrupt event'
		keyChan <- keyboard.KeyEvent{Key: keyboard.KeyCtrlC}

		// This should not cause panic
		kc.Stop()
	})

	t.Run("keystrokes channel is flushed on every tick", func(t *testing.T) {
		// TODO: empty k.keystrokes on every tick; measurements from the "last interval" are now affecting the "next" count
	})
}

// createCounterWithFakeChannel returns a DefaultKeystrokeCounter,
// its Count() channel, and another fake keystrokes channel.
func createCounterWithFakeChannel(t testing.TB) (DefaultKeystrokeCounter, chan keyboard.KeyEvent, <-chan int) {
	t.Helper()

	// Create a fake channel through which custom keystrokes can be sent
	keyChan := make(chan keyboard.KeyEvent)
	kc := DefaultKeystrokeCounter{keystrokes: keyChan}

	// Start stroke count
	strokeCount := kc.Count(NewDefaultTicker(0 * time.Millisecond))

	return kc, keyChan, strokeCount
}

// spyTicker implements main.Ticker interface for testing purposes.
type spyTicker struct {
	interval time.Duration
	bufSize  int
	Calls    int
	c        <-chan time.Time
}

// newSpyTicker creates a spyTicker with a pre-filled ping channel.
func newSpyTicker(interval time.Duration) *spyTicker {
	n := 5
	c := make(chan time.Time, n)

	// Populate channel in advance, so values can be read immediately
	for i := 0; i < n; i++ {
		c <- time.Time{}
	}

	return &spyTicker{
		interval: interval,
		c:        c,
		bufSize:  n,
	}
}

// C returns the ping channel.
func (s *spyTicker) C() <-chan time.Time {
	return s.c
}

// Stop calculates the number of values that were read from the ping channel.
func (s *spyTicker) Stop() {
	s.Calls = s.bufSize - len(s.c)
}
