package main

import "testing"

func TestRegularMap(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$", 23},
		{"^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$", 31},
		{"^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$", 18},
		{"^ENWWW(NEEE|SSE(EE|N))$", 10},
		{"^WNE$", 3},
	}
	for _, c := range cases {
		got := RegularMapV2(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}
