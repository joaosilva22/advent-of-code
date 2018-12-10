package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Point struct {
	x  int
	y  int
	vx int
	vy int
}

func NewPoint(x, y, vx, vy int) *Point {
	point := Point{x, y, vx, vy}
	return &point
}

func (p *Point) Update(delta int) {
	p.x += delta * p.vx
	p.y += delta * p.vy
}

func (p *Point) Rewind(delta int) {
	p.x -= delta * p.vx
	p.y -= delta * p.vy
}

func calculateBounds(points []*Point) (minX, minY, maxX, maxY int) {
	minX, minY = math.MaxInt32, math.MaxInt32
	maxX, maxY = math.MinInt32, math.MinInt32
	for _, point := range points {
		if point.x < minX {
			minX = point.x
		}
		if point.y < minY {
			minY = point.y
		}
		if point.x > maxX {
			maxX = point.x
		}
		if point.y > maxY {
			maxY = point.y
		}
	}
	return
}

func printPoints(points []*Point, minX, maxX, minY, maxY int) {
	width := maxX - minX + 1
	height := maxY - minY + 1

	sky := make([][]string, height)
	for _, point := range points {
		if sky[point.y-minY] == nil {
			sky[point.y-minY] = make([]string, width)
		}
		sky[point.y-minY][point.x-minX] = "#"
	}

	for row := 0; row < len(sky); row++ {
		for col := 0; col < len(sky[row]); col++ {
			fmt.Printf("%1s", sky[row][col])
		}
		fmt.Println()
	}
}

func TheStarsAlign(points []*Point) {
	done := false
	for !done {
		minX, minY, maxX, maxY := calculateBounds(points)
		fmt.Printf("X=[%d, %d], Y=[%d, %d]\n", minX, maxX, minY, maxY)

		var cmd string
		var arg int
		fmt.Print("> ")
		fmt.Scanf("%s %d\n", &cmd, &arg)

		switch cmd {
		case "f":
			delta := arg
			if delta == 0 {
				delta = 1
			}
			for _, point := range points {
				point.Update(delta)
			}
		case "b":
			fmt.Println("Backward", arg)
			delta := arg
			if delta == 0 {
				delta = 1
			}
			for _, point := range points {
				point.Rewind(delta)
			}
		case "p":
			printPoints(points, minX, maxX, minY, maxY)
		case "q":
			done = true
		default:
			fmt.Println("Unknown command")
		}
	}
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

	points := []*Point{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var x, y, vx, vy int
		format := "position=<%d, %d> velocity=<%d, %d>"
		fmt.Sscanf(scanner.Text(), format, &x, &y, &vx, &vy)
		point := NewPoint(x, y, vx, vy)
		points = append(points, point)
	}

	TheStarsAlign(points)
}
