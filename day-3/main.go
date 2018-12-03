package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type claim struct {
	id     int
	x      int
	y      int
	width  int
	height int
}

func GetIntactClaim(claims []claim) int {
	fabric := map[int]map[int]int{}
	intact := -1
	for _, claim := range claims {
		for row := claim.y; row < claim.y+claim.height; row++ {
			for col := claim.x; col < claim.x+claim.width; col++ {
				if fabric[row] == nil {
					fabric[row] = map[int]int{}
				}
				fabric[row][col] += 1
			}
		}
	}
	for _, claim := range claims {
		valid := true
		for row := claim.y; row < claim.y+claim.height; row++ {
			for col := claim.x; col < claim.x+claim.width; col++ {
				if fabric[row][col] != 1 {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			intact = claim.id
			break
		}
	}
	return intact
}

func GetOverlappingInches(claims []claim) int {
	fabric := map[int]map[int]int{}
	overlapping := 0

	for _, claim := range claims {
		for row := claim.y; row < claim.y+claim.height; row++ {
			for col := claim.x; col < claim.x+claim.width; col++ {
				if fabric[row] == nil {
					fabric[row] = map[int]int{}
				}
				fabric[row][col] += 1
				if fabric[row][col] == 2 {
					overlapping += 1
				}
			}
		}
	}

	return overlapping
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

	var claims []claim

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var id, x, y, width, height int
		_, err := fmt.Sscanf(scanner.Text(), "#%d @ %d,%d: %dx%d", &id, &x, &y, &width, &height)
		if err != nil {
			log.Fatal(err)
		}
		claim := claim{id, x, y, width, height}
		claims = append(claims, claim)
	}

	overlapping := GetOverlappingInches(claims)
	fmt.Println(overlapping)

	intact := GetIntactClaim(claims)
	fmt.Println(intact)
}
