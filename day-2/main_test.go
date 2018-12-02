package main

import "testing"

func TestChecksum(t *testing.T) {
	cases := []struct {
		in   []string
		want int
	}{
		{[]string{
			"abcdef",
			"bababc",
			"abbcde",
			"abcccd",
			"aabcdd",
			"abcdee",
			"ababab",
		}, 12},
	}
	for _, c := range cases {
		got := Checksum(c.in)
		if got != c.want {
			t.Errorf("Expected %d got %d", c.want, got)
		}
	}
}

func TestCommonLetters(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{
			"abcde",
			"fghij",
			"klmno",
			"pqrst",
			"fguij",
			"axcye",
			"wvxyz",
		}, "fgij"},
	}
	for _, c := range cases {
		got := CommonLetters(c.in)
		if got != c.want {
			t.Errorf("Expected %s got %s", c.want, got)
		}
	}
}
