package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func AlchemicalReductionPart2(polymer string) int {
	min := math.MaxInt32

	for delete := 65; delete < 91; delete++ {
		deleted := make([]bool, len(polymer))
		for charIndex := 0; charIndex < len(polymer); charIndex++ {
			char := polymer[charIndex]
			if int(char) == delete || int(char) == delete+32 {
				deleted[charIndex] = true
			}
		}

		var sb strings.Builder
		for i := 0; i < len(polymer); i++ {
			if deleted[i] {
				continue
			}
			sb.WriteByte(polymer[i])
		}

		len := AlchemicalReduction(sb.String())
		if len < min {
			min = len
		}

	}

	return min
}

func AlchemicalReduction(polymer string) int {
	reacted := true
	deleted := make([]bool, len(polymer))

	for reacted {
		reacted = false
		prevIndex := 0
		for deleted[prevIndex] {
			if prevIndex < len(polymer)-1 {
				prevIndex++
			} else {
				break
			}
		}
		for currIndex := prevIndex + 1; currIndex < len(polymer); currIndex++ {
			if deleted[currIndex] {
				continue
			}
			curr := polymer[currIndex]
			prev := polymer[prevIndex]

			if int(curr)-int(prev) == 32 || int(curr)-int(prev) == -32 {
				deleted[currIndex] = true
				deleted[prevIndex] = true
				reacted = true
				break
			}

			prevIndex = currIndex
		}
	}

	var sb strings.Builder
	for i := 0; i < len(polymer); i++ {
		if deleted[i] {
			continue
		}
		sb.WriteByte(polymer[i])
	}

	return len(sb.String())
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

	var polymer string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		polymer = scanner.Text()
	}

	part1 := AlchemicalReduction(polymer)
	fmt.Println(part1)

	part2 := AlchemicalReductionPart2(polymer)
	fmt.Println(part2)
}
