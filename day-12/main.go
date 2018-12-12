package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func SubterraneanSustainability(initialState string, lookup map[string]string, gens int) int64 {
	var leftmost int64 = 0
	state := initialState
	for gen := 0; gen < gens; gen++ {
		start := 0
		for i := 0; i < 5; i++ {
			if state[i] == '.' {
				start++
			} else {
				break
			}
		}
		for i := 0; i < 5-start; i++ {
			state = "." + state
		}
		end := 0
		for i := len(state) - 1; i >= len(state)-5; i-- {
			if state[i] == '.' {
				end++
			} else {
				break
			}
		}
		for i := 0; i < 5-end; i++ {
			state = state + "."
		}

		leftmost -= 5 - int64(start) - 2

		var sb strings.Builder
		for i := 2; i < len(state)-2; i++ {
			next := lookup[state[i-2:i+3]]
			if next == "" {
				next = "."
			}
			sb.WriteString(next)
		}

		nextState := sb.String()
		if strings.Trim(nextState, ".") == strings.Trim(state, ".") {
			remainingGens := int64(gens - (gen + 1))
			leftmost += remainingGens
			state = nextState
			break
		}
		state = nextState
	}
	var res int64 = 0
	for i := 0; i < len(state); i++ {
		if state[i] == '#' {
			res += int64(i) + leftmost
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		message := fmt.Sprintf("Usage: %s <input-file>", os.Args[0])
		log.Fatal(message)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var initialState string
	lookup := map[string]string{}

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "initial state: %s", &initialState)

	scanner.Scan()
	for scanner.Scan() {
		var key, val string
		fmt.Sscanf(scanner.Text(), "%s => %s", &key, &val)
		lookup[key] = val
	}

	part1 := SubterraneanSustainability(initialState, lookup, 20)
	fmt.Println(part1)

	part2 := SubterraneanSustainability(initialState, lookup, 50000000000)
	fmt.Println(part2)
}
