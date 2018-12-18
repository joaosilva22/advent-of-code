package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	// "math/rand"
	"os"
)

type Tile rune

const (
	Water      Tile = '~'
	WaterTrail Tile = '|'
	Sand       Tile = '.'
	Clay       Tile = '#'
)

type Vein struct {
	x0 int
	x1 int
	y0 int
	y1 int
}

type Ground struct {
	layout [][]Tile
	spawn  int
	width  int
	height int
	water  [][2]int
}

func NewGround(scan []Vein) *Ground {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, vein := range scan {
		if vein.x0 < minX {
			minX = vein.x0
		}
		if vein.x1 > maxX {
			maxX = vein.x1
		}
		if vein.y0 < minY {
			minY = vein.y0
		}
		if vein.y1 > maxY {
			maxY = vein.y1
		}
	}

	minX = minX - 1
	maxX = maxX + 1

	width := maxX - minX + 1
	height := maxY - minY + 1
	layout := [][]Tile{}

	for row := 0; row < height; row++ {
		layout = append(layout, []Tile{})
		for col := 0; col < width; col++ {
			layout[row] = append(layout[row], Sand)
		}
	}

	for _, vein := range scan {
		for row := vein.y0; row <= vein.y1; row++ {
			for col := vein.x0; col <= vein.x1; col++ {
				layout[row-minY][col-minX] = Clay
			}
		}
	}

	spawn := 500 - minX

	g := Ground{layout, spawn, width, height, [][2]int{}}
	return &g
}

func (g *Ground) Print() {
	for _, line := range g.layout {
		for _, tile := range line {
			fmt.Printf("%c", tile)
		}
		fmt.Println()
	}
}

func (g *Ground) PrintToFile(name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range g.layout {
		for _, tile := range line {
			_, err := w.WriteString(fmt.Sprintf("%c", tile))
			if err != nil {
				log.Fatal(err)
			}
		}
		w.WriteString("\n")
	}

	w.Flush()
}

func (g *Ground) Update() {
	waterTiles := [][2]int{}
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			if g.layout[row][col] == WaterTrail {
				pos := [2]int{row, col}
				waterTiles = append(waterTiles, pos)
			}
		}
	}

	for _, tile := range waterTiles {
		row, col := tile[0], tile[1]
		// Always move down if possible
		if row+1 >= g.height {
			g.layout[row][col] = WaterTrail
		} else if g.layout[row+1][col] == Sand || g.layout[row+1][col] == WaterTrail {
			g.layout[row+1][col] = WaterTrail
		} else {
			// If you can't move down, move to the side
			// sides := [2]int{-1, 1}
			side := 1

			if g.layout[row][col+side] == Sand || g.layout[row][col+side] == WaterTrail {
				g.layout[row][col+side] = WaterTrail
			}
			if g.layout[row][col-side] == Sand || g.layout[row][col-side] == WaterTrail {
				g.layout[row][col-side] = WaterTrail
			}
		}

		// g.Print()
		// var val string
		// fmt.Scanf("%s", val)

	}

	for row := 0; row < g.height; row++ {
		streak := false
		begin := -1
		for col := 0; col < g.width; col++ {
			tile := g.layout[row][col]

			if !streak {
				if tile == Clay {
					streak = true
					begin = col + 1
				}
			} else if streak {
				if tile == Clay {
					if col-begin == 1 {
						if row+1 < g.height && g.layout[row+1][begin] != Clay && g.layout[row+1][begin] != Water {
							continue
						}
					}
					for i := begin; i < col; i++ {
						g.layout[row][i] = Water
					}
					begin = col + 1
				} else if tile != WaterTrail {
					streak = false
				}
			}
		}
	}

	if g.layout[0][g.spawn] != Water {
		g.layout[0][g.spawn] = WaterTrail
		g.water = append(g.water, [2]int{0, g.spawn})
	}
}

func (g *Ground) CountWaterTiles() int {
	count := 0
	for _, line := range g.layout {
		for _, tile := range line {
			if tile == Water || tile == WaterTrail {
				count++
			}
		}
	}
	return count
}

func (g *Ground) CountWaterAtRest() int {
	count := 0
	for _, line := range g.layout {
		for _, tile := range line {
			if tile == Water {
				count++
			}
		}
	}
	return count
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

	scan := []Vein{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var a0, a1, a2 int

		_, err := fmt.Sscanf(scanner.Text(), "x=%d, y=%d..%d", &a0, &a1, &a2)
		if err == nil {
			vein := Vein{a0, a0, a1, a2}
			scan = append(scan, vein)
			continue
		}

		_, err = fmt.Sscanf(scanner.Text(), "y=%d, x=%d..%d", &a0, &a1, &a2)
		if err == nil {
			vein := Vein{a1, a2, a0, a0}
			scan = append(scan, vein)
			continue
		}
	}

	g := NewGround(scan)

	prevWaterTiles := -1
	count := 0
	for count < 20 {
		// g.Print()
		g.Update()
		waterTiles := g.CountWaterTiles()

		if waterTiles == prevWaterTiles {
			count++
		} else {
			prevWaterTiles = waterTiles
			count = 0
		}

		// var val string
		// fmt.Scanf("%s", val)
	}

	g.PrintToFile("please")
	fmt.Println(prevWaterTiles)
	fmt.Println(g.CountWaterAtRest())
}
