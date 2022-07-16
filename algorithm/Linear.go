package algorithm

import "fmt"

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
	return Result{Volume: VolumeSetBoundaries(condition.MinValue, condition.MaxValue, volume)}
}
