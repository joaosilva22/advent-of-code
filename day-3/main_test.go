package main

import "testing"

func TestGetOverlappingInches(t *testing.T) {
	cases := []struct {
		in   []claim
		want int
	}{
		{[]claim{
			claim{1, 1, 3, 4, 4},
			claim{2, 3, 1, 4, 4},
			claim{3, 5, 5, 2, 2},
		}, 4},
	}
	for _, c := range cases {
		got := GetOverlappingInches(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestGetIntactClaim(t *testing.T) {
	cases := []struct {
		in   []claim
		want int
	}{
		{[]claim{
			claim{1, 1, 3, 4, 4},
			claim{2, 3, 1, 4, 4},
			claim{3, 5, 5, 2, 2},
		}, 3},
	}
	for _, c := range cases {
		got := GetIntactClaim(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
