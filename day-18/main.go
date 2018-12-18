package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Acre rune

const (
	Open       Acre = '.'
	Trees      Acre = '|'
	Lumberyard Acre = '#'
)

func printLayout(layout [][]Acre) {
	for _, line := range layout {
		for _, acre := range line {
			fmt.Printf("%c", acre)
		}
		fmt.Println()
	}
}

func countAdjacent(layout [][]Acre, row, col int, acre Acre) int {
	count := 0
	if row-1 >= 0 {
		if layout[row-1][col] == acre {
			count++
		}
		if col-1 >= 0 {
			if layout[row-1][col-1] == acre {
				count++
			}
		}
		if col+1 < len(layout[row]) {
			if layout[row-1][col+1] == acre {
				count++
			}
		}
	}
	if row+1 < len(layout) {
		if layout[row+1][col] == acre {
			count++
		}
		if col-1 >= 0 {
			if layout[row+1][col-1] == acre {
				count++
			}
		}
		if col+1 < len(layout[row]) {
			if layout[row+1][col+1] == acre {
				count++
			}
		}
	}
	if col-1 >= 0 {
		if layout[row][col-1] == acre {
			count++
		}
	}
	if col+1 < len(layout[row]) {
		if layout[row][col+1] == acre {
			count++
		}
	}

	return count
}

func updateLayout(layout [][]Acre) [][]Acre {
	newLayout := [][]Acre{}

	for row, line := range layout {
		newLine := []Acre{}
		for col, acre := range line {
			if acre == Open {
				trees := countAdjacent(layout, row, col, Trees)
				if trees >= 3 {
					newLine = append(newLine, Trees)
				} else {
					newLine = append(newLine, Open)
				}
			}
			if acre == Trees {
				lumberyards := countAdjacent(layout, row, col, Lumberyard)
				if lumberyards >= 3 {
					newLine = append(newLine, Lumberyard)
				} else {
					newLine = append(newLine, Trees)
				}
			}
			if acre == Lumberyard {
				lumberyards := countAdjacent(layout, row, col, Lumberyard)
				trees := countAdjacent(layout, row, col, Trees)

				if lumberyards >= 1 && trees >= 1 {
					newLine = append(newLine, Lumberyard)
				} else {
					newLine = append(newLine, Open)
				}
			}
		}
		newLayout = append(newLayout, newLine)
	}
	return newLayout
}

func resourceValue(layout [][]Acre) (int, int, int) {
	lumberyards := 0
	trees := 0
	for _, line := range layout {
		for _, tile := range line {
			if tile == Lumberyard {
				lumberyards++
			}
			if tile == Trees {
				trees++
			}
		}
	}
	return lumberyards * trees, lumberyards, trees
}

func copyLayout(layout [][]Acre) [][]Acre {
	copy := [][]Acre{}
	for _, line := range layout {
		lineCopy := []Acre{}
		for _, acre := range line {
			lineCopy = append(lineCopy, acre)
		}
		copy = append(copy, lineCopy)
	}
	return copy
}

func compareLayouts(l1, l2 [][]Acre) bool {
	for row, _ := range l1 {
		for col, _ := range l2 {
			if l1[row][col] != l2[row][col] {
				return false
			}
		}
	}
	return true
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

	var layout [][]Acre

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []Acre{}
		for _, a := range scanner.Text() {
			row = append(row, Acre(a))
		}
		layout = append(layout, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	otherLayout := layout
	limit := 1000000000 - (28 * 35713285)

	// limit := 18025 - (28 * 100)

	for i := 0; i < limit; i++ {
		newLayout := updateLayout(layout)
		// printLayout(newLayout)
		val, lumb, trees := resourceValue(newLayout)
		fmt.Println(val, lumb, trees, i)
		// fmt.Println()
		if compareLayouts(layout, newLayout) || compareLayouts(newLayout, otherLayout) {
			fmt.Println("Steady state!")
			otherLayout = copyLayout(layout)
			layout = newLayout
			break
		}
		otherLayout = copyLayout(layout)
		layout = newLayout
	}

	// printLayout(layout)
	fmt.Println(resourceValue(layout))
}
