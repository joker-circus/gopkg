package numberutil

import "testing"

func TestUint64ToFloat64(t *testing.T) {
	type args struct {
		num uint64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "A", args: args{num: 17}, want: 17},
		{name: "B", args: args{num: 0}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64ToFloat64(tt.args.num); got != tt.want {
				t.Errorf("Uint64ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound(t *testing.T) {
	type args struct {
		num  float64
		prec int32
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "", args: args{num: 0.123456789, prec: 3}, want: 0.123},
		{name: "", args: args{num: 0.123456789, prec: 4}, want: 0.1235},
		{name: "", args: args{num: 0.123456789, prec: -1}, want: 0},
		{name: "", args: args{num: 0.123456789, prec: 0}, want: 0},
		{name: "", args: args{num: 0.123456789, prec: 10}, want: 0.123456789},
		{name: "", args: args{num: 9.825, prec: 2}, want: 9.83},
		{name: "", args: args{num: 9.835, prec: 2}, want: 9.84},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round(tt.args.num, tt.args.prec); got != tt.want {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrecision(t *testing.T) {
	type args struct {
		num float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "", args: args{num: 0.123456789}, want: 9},
		{name: "", args: args{num: 0.123450000}, want: 5},
		{name: "", args: args{num: 7}, want: 0},
		{name: "", args: args{num: 7.0}, want: 0},
		{name: "", args: args{num: 0}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Precision(tt.args.num); got != tt.want {
				t.Errorf("Precision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	type args struct {
		f1   float64
		f2   float64
		prec int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "", args: args{f1: 0.3, f2: 0.30001, prec: 3}, want: true},
		{name: "", args: args{f1: 0.3, f2: 0.29999, prec: 3}, want: true},
		{name: "", args: args{f1: 0.3, f2: 0.30100, prec: 3}, want: false},
		{name: "", args: args{f1: 0.3, f2: 0.30010, prec: 3}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEqual(tt.args.f1, tt.args.f2); got != tt.want {
				t.Errorf("IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
