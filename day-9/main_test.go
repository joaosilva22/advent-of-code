package main

import "testing"

func TestMarbleMania(t *testing.T) {
	cases := []struct {
		players int
		marbles int
		want    int
	}{
		{9, 25, 32},
		{10, 1618, 8317},
		{13, 7999, 146373},
		{17, 1104, 2764},
		{21, 6111, 54718},
		{30, 5807, 37305},
	}
	for _, c := range cases {
		got := MarbleManiaXL(c.players, c.marbles)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
