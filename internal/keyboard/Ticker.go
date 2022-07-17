package keyboard

import "time"

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
