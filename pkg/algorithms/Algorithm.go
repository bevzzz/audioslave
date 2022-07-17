package algorithms

// Algorithm - Interface for implementing algorithms
type Algorithm interface {
	Name() string                      // Algorithm name
	Adjust(condition Condition) Result // Algorithm apply function
	ToString() string                  // To String
}

// Clamp - set the boundaries
func Clamp(min, max, volume int) int {
	if volume > max {
		volume = max
	}
	if volume < min {
		volume = min
	}
	return volume
}

// CalculateMean - calculates the mean
func CalculateMean(values ...int) int {
	mean := 0
	for _, value := range values {
		mean += value
	}
	return mean / len(values)
}

// MapToRange - Maps a value between 0-1 to a value between the range given
func MapToRange(min, max int, value float64) int {
	if value < 0 {
		value = 0
	} else if value > 1 {
		value /= 100
	}

	r := max - min
	return int(float64(r)*value) + min
}

func AlgorithmByName(name string) Algorithm {
	switch name {
	case "Linear":
		return &Linear{}
	case "LinearDependentOnAverageCPM":
		return &LinearDependentOnAverageCPM{}
	case "LinearDependentOnMean":
		return &LinearDependentOnMean{}
	default:
		return nil
	}
}
