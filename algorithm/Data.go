package algorithm

import "time"

// Condition - Values for the algorithm to calculate the new value
type Condition struct {
	AverageCPM int // average CPM

	Interval time.Duration // Interval per previous value

	PreviousStrokes []int // history of previous strokes
	CurrentStrokes  int   // current stroke count in the interval

	PreviousVolumes []int // previous calculated volumes

	MinValue int // minimum value
	MaxValue int // maximum value
}

// Result - Algorithm result
type Result struct {
	Volume int
}
