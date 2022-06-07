package main

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestStrokes(t *testing.T) {
	// Expect that 2 groups of keystrokes are identified
	var got []int
	var intervals []time.Duration

	interval := time.Duration(50)
	want := []int{2, 5, 6}

	// Create a fake channel through which keystrokes will be sent
	keyChan := make(chan int, 6)
	kb := Keyboard{keystrokes: keyChan}

	nStrokes := kb.Strokes(interval * time.Millisecond)

	// Simulate keystrokes being sent through the channel
	for _, n := range want {
		start := time.Now()

		for i := 0; i < n; i++ {
			keyChan <- n
		}
		// Strokes() should post count every 1s
		got = append(got, <-nStrokes)
		intervals = append(intervals, time.Since(start))
	}

	t.Run("number of keystrokes captured", func(t *testing.T) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("sends intervals at a specified rate", func(t *testing.T) {
		for _, i := range intervals {
			ms := time.Duration(i.Milliseconds())
			if math.Abs(float64(ms-interval)) > 20 {
				t.Fatalf("want %dms, got %dms", interval, ms)
			}
		}
	})

}
