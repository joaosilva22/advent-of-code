package main

import "testing"

func TestChronalCalibration(t *testing.T) {
	cases := []struct {
		in   []int64
		want int64
	}{
		{[]int64{1, 2, 3}, 6},
		{[]int64{}, 0},
		{[]int64{-1, -2, -3}, -6},
		{[]int64{1, -2, 3}, 2},
	}
	for _, c := range cases {
		got := ChronalCalibration(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestChronalCalibrationPart2(t *testing.T) {
	cases := []struct {
		in   []int64
		want int64
	}{
		{[]int64{+1, -1}, 0},
		{[]int64{+3, +3, +4, -2, -4}, 10},
		{[]int64{-6, +3, +8, +5, -6}, 5},
		{[]int64{+7, +7, -2, -7, -4}, 14},
	}
	for _, c := range cases {
		got := ChronalCalibrationPart2(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
