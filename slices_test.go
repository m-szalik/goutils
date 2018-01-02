package goutils

import "testing"

var sliceA = []interface{}{"A", "b", 11, 3.14}

func TestSliceIndexOf(t *testing.T) {
	type args struct {
		slice []interface{}
		e     interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name:"NotFound", args:args{slice:sliceA, e:"bar"}, want:-1},
		{name:"Zero", args:args{slice:sliceA, e:"A"}, want:0},
		{name:"3.14", args:args{slice:sliceA, e:3.14}, want:0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceIndexOf(tt.args.slice, tt.args.e); got != tt.want {
				t.Errorf("SliceIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	type args struct {
		slice []interface{}
		e     interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name:"NotFound", args:args{slice:sliceA, e:"bar"}, want:false},
		{name:"Zero", args:args{slice:sliceA, e:"A"}, want:true},
		{name:"3.14", args:args{slice:sliceA, e:3.14}, want:true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceContains(tt.args.slice, tt.args.e); got != tt.want {
				t.Errorf("SliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
