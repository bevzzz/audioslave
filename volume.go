package main

type Volume struct {
	on        bool
	threshold int
}

func NewVolume(threshold int) *Volume {
	return &Volume{on: true, threshold: threshold}
}

func (v *Volume) On() bool {
	return v.on
}

func (v *Volume) Adjust(nKeystrokes int) {
	if nKeystrokes >= v.threshold {
		v.on = !v.on
	}
}
