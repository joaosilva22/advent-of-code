package main

import "testing"

func TestPowerLevelXY(t *testing.T) {
	cases := []struct {
		x      int
		y      int
		serial int
		want   int
	}{
		{3, 5, 8, 4},
		{122, 79, 57, -5},
		{217, 196, 39, 0},
		{101, 153, 71, 4},
	}
	for _, c := range cases {
		got := PowerLevelXY(c.x, c.y, c.serial)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestChronalCharge(t *testing.T) {
	cases := []struct {
		serial int
		wantX  int
		wantY  int
	}{
		{18, 33, 45},
		{42, 21, 61},
	}
	for _, c := range cases {
		gotX, gotY := ChronalCharge(c.serial)
		if gotX != c.wantX || gotY != c.wantY {
			t.Errorf("Expected %d,%d got %d,%d", c.wantX, c.wantY, gotX, gotY)
		}
	}
}
