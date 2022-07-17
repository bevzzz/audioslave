package keyboard

import (
	"reflect"
	"testing"
	"time"
)

type KeyLoggerMock struct {
	Amount int
}

func (k *KeyLoggerMock) GetKey() Key {
	k.Amount--
	if k.Amount < 0 {
		return Key{
			Empty:   true,
			Rune:    ' ',
			Keycode: 0,
		}
	}
	return Key{
		Empty:   false,
		Rune:    ' ',
		Keycode: 00,
	}
}

func TestCount(t *testing.T) {

	t.Run("number of keystrokes is captured", func(t *testing.T) {
		var got []int
		want := []int{2, 5, 6}

		// Simulate keystrokes being sent through the channel
		for _, n := range want {
			kl := &KeyLoggerMock{Amount: n}
			counter := DefaultKeystrokeCounter{
				Keylogger: kl,
				Cancel:    nil,
			}
			ch := counter.Count(NewDefaultTicker(0))
			got = append(got, <-ch)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("uses ticker to send updates at a certain interval", func(t *testing.T) {
		kc := DefaultKeystrokeCounter{
			Keylogger: &KeyLoggerMock{},
			Cancel:    nil,
		}

		spyTicker := newSpyTicker(5 * time.Millisecond)
		kc.Count(spyTicker)
		time.Sleep(time.Millisecond)
		spyTicker.Stop()

		if spyTicker.Calls == 0 {
			t.Error("expected calls to ticker, but did not get any")
		}
	})

	t.Run("goroutine does not send on closed channel", func(t *testing.T) {
		want := 1
		kl := &KeyLoggerMock{Amount: want}
		counter := DefaultKeystrokeCounter{
			Keylogger: kl,
			Cancel:    nil,
		}
		ch := counter.Count(NewDefaultTicker(0))
		n := <-ch
		if n != want {
			t.Errorf("got %d, want %d", n, want)
		}
		counter.Stop()
		n = <-ch
		if n != 0 {
			t.Errorf("got %d, want %d", n, 0)
		}
	})

	t.Run("keystrokes channel is flushed on every tick", func(t *testing.T) {
		// TODO: empty k.keystrokes on every tick; measurements from the "last interval" are now affecting the "next" count
	})
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
