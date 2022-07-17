package algorithms

import "fmt"

// Linear - Linear alg which dependent on the current strokes
type Linear struct {
	IncreaseBy int
	ReduceBy   int
}

func (l Linear) Name() string {
	return fmt.Sprintf("Linear: Increase by: %d, Reduce by: %d", l.IncreaseBy, l.ReduceBy)
}

func (l Linear) Adjust(condition Condition) Result {
	previousVolume := condition.PreviousVolumes[len(condition.PreviousVolumes)-1]
	volume := previousVolume
	if condition.CurrentStrokes > condition.PreviousStrokes[len(condition.PreviousStrokes)-1] {
		volume += l.IncreaseBy
	} else {
		volume -= l.ReduceBy
	}
	return Result{Volume: Clamp(condition.MinValue, condition.MaxValue, volume)}
}

// LinearDependentOnAverageCPM - Linear alg which dependent on the averageCPM
type LinearDependentOnAverageCPM struct {
	Linear
}

func (l LinearDependentOnAverageCPM) Name() string {
	return fmt.Sprintf("LinearDependentOnAverageCPM: Increase by: %d, Reduce by: %d", l.IncreaseBy, l.ReduceBy)
}

func (l LinearDependentOnAverageCPM) Adjust(condition Condition) Result {
	volume := condition.PreviousVolumes[len(condition.PreviousVolumes)-1]
	if condition.CurrentStrokes > condition.AverageCPM {
		volume += l.IncreaseBy
	} else {
		volume -= l.ReduceBy
	}
	return Result{Volume: Clamp(condition.MinValue, condition.MaxValue, volume)}
}

// LinearDependentOnMean - Linear alg which dependent on the mean strokes of a certain timeFrame
type LinearDependentOnMean struct {
	Linear
	TimeFrame int
}

func (l LinearDependentOnMean) Name() string {
	return fmt.Sprintf("LinearDependentOnMean: Increase by: %d, Reduce by: %d", l.IncreaseBy, l.ReduceBy)
}

func (l LinearDependentOnMean) Adjust(condition Condition) Result {
	previousStrokesLength := len(condition.PreviousStrokes)
	if previousStrokesLength > l.TimeFrame {
		previousStrokesLength = l.TimeFrame
	}
	mean := CalculateMean(condition.PreviousStrokes[:previousStrokesLength]...)
	volume := condition.PreviousVolumes[len(condition.PreviousVolumes)-1]
	if condition.CurrentStrokes > mean {
		volume += l.IncreaseBy
	} else {
		volume -= l.ReduceBy
	}
	return Result{Volume: Clamp(condition.MinValue, condition.MaxValue, volume)}
}
