package main

import "testing"

func TestAlchemicalReduction(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"dabAcCaCBAcCcaDA", 10},
		{"dDaBcCbA", 0},
	}
	for _, c := range cases {
		got := AlchemicalReduction(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestAlchemicalReductionPart2(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"dabAcCaCBAcCcaDA", 4},
	}
	for _, c := range cases {
		got := AlchemicalReductionPart2(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
