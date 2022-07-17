package algorithms

import "testing"

func TestCalculateMean(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Mean",
			args: struct{ values []int }{values: []int{3, 5, 7}},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateMean(tt.args.values...); got != tt.want {
				t.Errorf("CalculateMean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	type args struct {
		min    int
		max    int
		volume int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Clamp not applied",
			args: struct {
				min    int
				max    int
				volume int
			}{min: 3, max: 5, volume: 4},
			want: 4,
		},
		{
			name: "Clamp applied",
			args: struct {
				min    int
				max    int
				volume int
			}{min: 3, max: 5, volume: 7},
			want: 5,
		},
		{
			name: "Clamp applied",
			args: struct {
				min    int
				max    int
				volume int
			}{min: 3, max: 5, volume: 1},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.args.min, tt.args.max, tt.args.volume); got != tt.want {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToRange(t *testing.T) {
	type args struct {
		min   int
		max   int
		value float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "MapToRange 1",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 0, max: 100, value: 0.5},
			want: 50,
		},
		{
			name: "MapToRange 2",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 0, max: 100, value: 50},
			want: 50,
		},
		{
			name: "MapToRange 3",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 20, max: 70, value: 0.5},
			want: 45,
		},
		{
			name: "MapToRange 4",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 20, max: 70, value: 0.25},
			want: 32,
		},
		{
			name: "MapToRange 5",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 20, max: 70, value: 25},
			want: 32,
		},
		{
			name: "MapToRange 6",
			args: struct {
				min   int
				max   int
				value float64
			}{min: 40, max: 70, value: 0.8},
			want: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToRange(tt.args.min, tt.args.max, tt.args.value); got != tt.want {
				t.Errorf("MapToRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
