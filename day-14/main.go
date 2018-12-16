package main

import (
	"fmt"
	"strconv"
)

func ChocolateCharts(n int) []int {
	state := []int{3, 7}
	elf1, elf2 := 0, 1

	for len(state) < n {
		sum := state[elf1] + state[elf2]
		sumStr := strconv.Itoa(sum)
		for _, c := range sumStr {
			state = append(state, int(c-'0'))
		}
		elf1 = (elf1 + 1 + state[elf1]) % len(state)
		elf2 = (elf2 + 1 + state[elf2]) % len(state)
	}

	return state
}

func ChocolateChartsPart2(goal string) int {
	state := []int{3, 7}
	elf1, elf2 := 0, 1
	done := false

	for !done {
		sum := state[elf1] + state[elf2]
		sumStr := strconv.Itoa(sum)
		added := 0
		for _, c := range sumStr {
			state = append(state, int(c-'0'))
			added += 1
		}
		elf1 = (elf1 + 1 + state[elf1]) % len(state)
		elf2 = (elf2 + 1 + state[elf2]) % len(state)
		// Check if done
		if len(state) >= len(goal) {
			for a := 0; a < added; a++ {
				start := len(state) - len(goal) - a
				if start < 0 {
					start = 0
				}
				done = true
				// fmt.Println(state, goal)
				for i := start; i < len(state); i++ {
					j := i - start
					// fmt.Println(state[i], int(goal[j]-'0'))
					if j >= len(goal) {
						break
					}
					if state[i] != int(goal[j]-'0') {
						done = false
						break
					}
				}
				if done {
					// fmt.Println(state, goal)
					break
				}
			}
		}
	}

	return len(state) - len(goal)
}

func main() {
	test1 := ChocolateCharts(9 + 10)
	fmt.Println(test1[9 : 9+10])

	test2 := ChocolateCharts(5 + 10)
	fmt.Println(test2[5 : 5+10])

	test3 := ChocolateCharts(2018 + 10)
	fmt.Println(test3[2018 : 2018+10])

	part1 := ChocolateCharts(704321 + 10)
	fmt.Println(part1[704321 : 704321+10])

	test4 := ChocolateChartsPart2("51589")
	fmt.Println(test4)

	test5 := ChocolateChartsPart2("01245")
	fmt.Println(test5)

	test6 := ChocolateChartsPart2("92510")
	fmt.Println(test6)

	test7 := ChocolateChartsPart2("59414")
	fmt.Println(test7)

	part2 := ChocolateChartsPart2("704321")
	fmt.Println(part2)
}
