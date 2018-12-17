package main

import "testing"

func TestPointFirstThan(t *testing.T) {
	cases := []struct {
		in   Point
		want bool
	}{
		{Point{0, 0}, false},
		{Point{0, 1}, true},
		{Point{1, 0}, true},
	}
	point := Point{0, 0}
	for _, c := range cases {
		got := point.FirstThan(c.in)
		if got != c.want {
			t.Errorf("Expected %v got %v", c.want, got)
		}
	}
}

func TestMapGetFreeAdjacentSquares(t *testing.T) {
	layout := [][]Entity{
		{'#', '#', '#', '#', '#', '#', '#'},
		{'#', 'E', '.', '.', 'G', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '#'},
		{'#', '.', 'G', '.', '#', 'G', '#'},
		{'#', '#', '#', '#', '#', '#', '#'},
	}
	mp := NewMap(layout)
	cases := []struct {
		in   Point
		want int
	}{
		{Point{1, 1}, 2},
		{Point{3, 2}, 3},
		{Point{1, 4}, 2},
		{Point{3, 5}, 1},
	}
	for _, c := range cases {
		got := mp.GetFreeAdjacentSquares(c.in)
		if len(got) != c.want {
			t.Errorf("Expected %d got %d", c.want, len(got))
		}
	}
}

func TestMapGetDestinations(t *testing.T) {
	layout := [][]Entity{
		{'#', '#', '#', '#', '#', '#', '#'},
		{'#', 'E', '.', '.', 'G', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '#'},
		{'#', '.', 'G', '.', '#', 'G', '#'},
		{'#', '#', '#', '#', '#', '#', '#'},
	}
	mp := NewMap(layout)
	cases := []struct {
		in   Point
		want int
	}{
		{Point{1, 1}, 6},
		{Point{3, 2}, 2},
		{Point{1, 4}, 2},
		{Point{3, 5}, 2},
	}
	for _, c := range cases {
		got := mp.GetDestinations(c.in)
		if len(got) != c.want {
			t.Errorf("Expected %d got %d", c.want, len(got))
		}
	}
}

func TestMapGetShortestPath(t *testing.T) {
	layout := [][]Entity{
		{'#', '#', '#', '#', '#', '#', '#'},
		{'#', 'E', '.', '.', 'G', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '#'},
		{'#', '.', 'G', '.', '#', 'G', '#'},
		{'#', '#', '#', '#', '#', '#', '#'},
	}
	mp := NewMap(layout)
	cases := []struct {
		from Point
		to   Point
		want int
	}{
		{Point{1, 1}, Point{3, 1}, 3},
		{Point{1, 1}, Point{1, 5}, 0},
	}
	for _, c := range cases {
		got, _ := mp.GetShortestPath(c.from, c.to)
		if len(got) != c.want {
			t.Errorf("Expected %d got %d", c.want, len(got))
		}
	}
}

func TestMapIsNextToTarget(t *testing.T) {
	layout := [][]Entity{
		{'#', '#', '#', '#', '#', '#', '#'},
		{'#', 'E', '.', '.', 'G', 'E', '#'},
		{'#', '.', '.', '.', '#', '.', '#'},
		{'#', '.', 'G', '.', '#', 'G', '#'},
		{'#', '#', '#', '#', '#', '#', '#'},
	}
	mp := NewMap(layout)
	cases := []struct {
		in   Point
		want bool
	}{
		{Point{1, 1}, false},
		{Point{1, 4}, true},
		{Point{1, 5}, true},
	}
	for _, c := range cases {
		got := mp.IsNextToTarget(c.in)
		if got != c.want {
			t.Errorf("Expected %v got %v", c.want, got)
		}
	}
}
