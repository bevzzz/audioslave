package algorithm

// Algorithm - Interface for implementing algorithms
type Algorithm interface {
	Name() string                      // Algorithm name
	Adjust(condition Condition) Result // Algorithm apply function
}

func VolumeSetBoundaries(min, max, volume int) int {
	if volume > max {
		volume = max
	}
	if volume < min {
		volume = min
	}
	return volume
}
