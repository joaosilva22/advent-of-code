package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Coordinate struct {
	ID int
	x  int
	y  int
}

func (c Coordinate) ManhattanDistance(o Coordinate) int {
	return abs(c.x-o.x) + abs(c.y-o.y)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func printGrid(grid map[int][]int, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Printf("%2d ", grid[y][x])
		}
		fmt.Println()
	}
}

func ChronalCoordinates(coords []Coordinate) int {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32

	for _, coord := range coords {
		minX = min(minX, coord.x)
		minY = min(minY, coord.y)
		maxX = max(maxX, coord.x)
		maxY = max(maxY, coord.y)
	}

	fmt.Printf("minX=%d maxX=%d\n", minX, maxX)
	fmt.Printf("minY=%d maxY=%d\n", minY, maxY)

	width := maxX - minX + 1
	height := maxY - minY + 1

	fmt.Printf("width=%d height=%d\n", width, height)

	grid := make(map[int][]int, height)

	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			other := Coordinate{0, x + minX, y + minY}
			minDistance := math.MaxInt32
			for _, coord := range coords {
				distance := coord.ManhattanDistance(other)
				if distance < minDistance {
					minDistance = distance
					grid[y][x] = coord.ID
				} else if distance == minDistance {
					grid[y][x] = -1
				}
			}
		}
	}

	printGrid(grid, width, height)

	maxNonInfiniteArea := 0
	for _, coord := range coords {
		area := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if grid[y][x] == coord.ID {
					area += 1
					if x == 0 || x == width-1 || y == 0 || y == height-1 {
						area = math.MinInt32
					}
				}
			}
		}
		fmt.Printf("Area of coord ID=%d (%d, %d) is %d\n", coord.ID, coord.x, coord.y, area)
		if area > maxNonInfiniteArea {
			maxNonInfiniteArea = area
		}
	}

	return maxNonInfiniteArea
}

func ChronalCoordinatesPart2(coords []Coordinate, maxTotalDistance int) int {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32

	for _, coord := range coords {
		minX = min(minX, coord.x)
		minY = min(minY, coord.y)
		maxX = max(maxX, coord.x)
		maxY = max(maxY, coord.y)
	}

	fmt.Printf("minX=%d maxX=%d\n", minX, maxX)
	fmt.Printf("minY=%d maxY=%d\n", minY, maxY)

	width := maxX - minX + 1
	height := maxY - minY + 1

	fmt.Printf("width=%d height=%d\n", width, height)

	grid := make(map[int][]int, height)

	safeArea := 0

	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			other := Coordinate{0, x + minX, y + minY}
			totalDistance := 0
			for _, coord := range coords {
				distance := coord.ManhattanDistance(other)
				totalDistance += distance
			}
			if totalDistance < maxTotalDistance {
				safeArea += 1
			}
		}
	}

	return safeArea
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

	var coords []Coordinate

	scanner := bufio.NewScanner(file)
	coordID := 1
	for scanner.Scan() {
		var x, y int
		fmt.Sscanf(scanner.Text(), "%d, %d", &x, &y)
		coord := Coordinate{coordID, x, y}
		coords = append(coords, coord)
		coordID += 1
	}

	for _, coord := range coords {
		fmt.Println(coord)
	}

	part1 := ChronalCoordinates(coords)
	fmt.Println(part1)

	part2 := ChronalCoordinatesPart2(coords, 10000)
	fmt.Println(part2)
}
