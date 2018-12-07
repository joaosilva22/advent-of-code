package main

import "testing"

func TestChronalCoordinates(t *testing.T) {
	cases := []struct {
		in   []Coordinate
		want int
	}{
		{[]Coordinate{
			Coordinate{1, 1, 1},
			Coordinate{2, 1, 6},
			Coordinate{3, 8, 3},
			Coordinate{4, 3, 4},
			Coordinate{5, 5, 5},
			Coordinate{6, 8, 9},
		}, 17},
	}
	for _, c := range cases {
		got := ChronalCoordinates(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestChronalCoordinatesPart2(t *testing.T) {
	cases := []struct {
		in   []Coordinate
		want int
	}{
		{[]Coordinate{
			Coordinate{1, 1, 1},
			Coordinate{2, 1, 6},
			Coordinate{3, 8, 3},
			Coordinate{4, 3, 4},
			Coordinate{5, 5, 5},
			Coordinate{6, 8, 9},
		}, 16},
	}
	for _, c := range cases {
		got := ChronalCoordinatesPart2(c.in, 32)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
