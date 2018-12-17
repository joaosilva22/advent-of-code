package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type Point struct {
	row int
	col int
}

func (p *Point) FirstThan(o Point) bool {
	if p.row < o.row {
		return true
	}
	return p.row == o.row && p.col < o.col
}

type Entity rune

const (
	Goblin    Entity = 'G'
	Elf       Entity = 'E'
	Wall      Entity = '#'
	FreeSpace Entity = '.'
)

type Map struct {
	layout      [][]Entity
	hp          [][]int
	attackPower int
}

func NewMap(layout [][]Entity, attackPower int) *Map {
	m := Map{}
	m.layout = layout
	m.attackPower = attackPower
	m.hp = [][]int{}
	for row, line := range layout {
		m.hp = append(m.hp, []int{})
		for _, entity := range line {
			if entity == Elf || entity == Goblin {
				m.hp[row] = append(m.hp[row], 200)
			} else {
				m.hp[row] = append(m.hp[row], 0)
			}
		}
	}
	return &m
}

func (m *Map) Get(p Point) Entity {
	return m.layout[p.row][p.col]
}

func (m *Map) GetAdjacentSquares(p Point) []Point {
	squares := []Point{}
	if p.row-1 > 0 {
		squares = append(squares, Point{p.row - 1, p.col})
	}
	if p.col-1 > 0 {
		squares = append(squares, Point{p.row, p.col - 1})
	}
	if p.col+1 < len(m.layout[p.row]) {
		squares = append(squares, Point{p.row, p.col + 1})
	}
	if p.row+1 < len(m.layout) {
		squares = append(squares, Point{p.row + 1, p.col})
	}
	return squares
}

func (m *Map) GetFreeAdjacentSquares(p Point) []Point {
	squares := m.GetAdjacentSquares(p)
	freeSquares := []Point{}
	for _, square := range squares {
		if m.Get(square) == FreeSpace {
			freeSquares = append(freeSquares, square)
		}
	}
	return freeSquares
}

func (m *Map) GetDestinations(from Point) []Point {
	targetEntity := Goblin
	if m.Get(from) == Goblin {
		targetEntity = Elf
	}
	res := []Point{}
	for row, line := range m.layout {
		for col, _ := range line {
			p := Point{row, col}
			if m.Get(p) == targetEntity {
				squares := m.GetFreeAdjacentSquares(p)
				res = append(res, squares...)
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].FirstThan(res[j])
	})
	return res
}

type Queue struct {
	contents []Point
}

func NewQueue() *Queue {
	q := Queue{[]Point{}}
	return &q
}

func (q *Queue) Enqueue(p Point) {
	q.contents = append(q.contents, p)
}

func (q *Queue) Pop() Point {
	res := q.contents[0]
	q.contents = q.contents[1:]
	return res
}

func (q *Queue) Contains(p Point) bool {
	contains := false
	for _, point := range q.contents {
		if p == point {
			contains = true
		}
	}
	return contains
}

func (q *Queue) IsEmpty() bool {
	return len(q.contents) == 0
}

type Set struct {
	contents []Point
}

func NewSet() *Set {
	s := Set{[]Point{}}
	return &s
}

func (s *Set) Contains(p Point) bool {
	contains := false
	for _, point := range s.contents {
		if p == point {
			contains = true
			break
		}
	}
	return contains
}

func (s *Set) Add(p Point) {
	if !s.Contains(p) {
		s.contents = append(s.contents, p)
	}
}

func (m *Map) GetShortestPath(from, to Point) ([]Point, bool) {
	border := NewQueue()
	visited := NewSet()

	meta := map[Point]Point{}

	border.Enqueue(from)

	for !border.IsEmpty() {
		root := border.Pop()

		if root == to {
			path := []Point{root}
			for val, ok := meta[root]; ok; val, ok = meta[root] {
				path = append(path, val)
				root = val
			}
			orderedPath := []Point{}
			for i := len(path) - 1; i >= 0; i-- {
				orderedPath = append(orderedPath, path[i])
			}
			return orderedPath, true
		}

		children := m.GetFreeAdjacentSquares(root)
		for _, child := range children {
			if visited.Contains(child) {
				continue
			}

			if !border.Contains(child) {
				border.Enqueue(child)
				meta[child] = root
			}
		}

		visited.Add(root)
	}

	return []Point{}, false
}

func (m *Map) IsNextToTarget(p Point) bool {
	targetEntity := Goblin
	if m.Get(p) == Goblin {
		targetEntity = Elf
	}

	adjacentSquares := m.GetAdjacentSquares(p)
	for _, square := range adjacentSquares {
		if m.Get(square) == targetEntity {
			return true
		}
	}

	return false
}

func (m *Map) Set(p Point, e Entity) {
	m.layout[p.row][p.col] = e
}

func (m *Map) Move(from, to Point) {
	entity := m.Get(from)
	m.Set(to, entity)
	m.Set(from, FreeSpace)

	hp := m.GetHP(from)
	m.SetHP(to, +hp)
	m.SetHP(from, -hp)
}

func (m *Map) GetTargetsInRange(p Point) []Point {
	targetEntity := Goblin
	if m.Get(p) == Goblin {
		targetEntity = Elf
	}

	targets := []Point{}

	adjacentSquares := m.GetAdjacentSquares(p)
	for _, square := range adjacentSquares {
		if m.Get(square) == targetEntity {
			targets = append(targets, square)
		}
	}

	return targets
}

func (m *Map) GetHP(target Point) int {
	return m.hp[target.row][target.col]
}

func (m *Map) SetHP(target Point, diff int) {
	m.hp[target.row][target.col] += diff
}

func (m *Map) Eliminate(target Point) {
	m.hp[target.row][target.col] = 0
	m.layout[target.row][target.col] = FreeSpace
}

func (m *Map) Attack(target Point) {
	if m.Get(target) == Goblin {
		m.SetHP(target, -m.attackPower)
	} else {
		m.SetHP(target, -3)
	}
	if m.GetHP(target) < 0 {
		m.Eliminate(target)
	}
}

func (m *Map) GetEntitiesInReadingOrder() []Point {
	positions := []Point{}
	for row, line := range m.layout {
		for col, entity := range line {
			if entity == Elf || entity == Goblin {
				positions = append(positions, Point{row, col})
			}
		}
	}
	return positions
}

func (m *Map) Print() {
	for row, line := range m.layout {
		for _, entity := range line {
			fmt.Printf("%c", entity)
		}
		fmt.Print(" ")
		for _, hp := range m.hp[row] {
			fmt.Printf("%3d", hp)
		}
		fmt.Println()
	}
}

func (m *Map) Over() bool {
	positions := m.GetEntitiesInReadingOrder()

	elves := false
	for _, pos := range positions {
		if m.Get(pos) == Elf {
			elves = true
			break
		}
	}

	goblins := false
	for _, pos := range positions {
		if m.Get(pos) == Goblin {
			goblins = true
			break
		}
	}

	return !elves || !goblins
}

func (m *Map) Update() bool {
	positions := m.GetEntitiesInReadingOrder()
	dead := NewSet()
	for _, point := range positions {
		if m.Over() {
			fmt.Println("Ended in middle of round")
			return false
		}

		if dead.Contains(point) {
			continue
		}

		// Move...
		if !m.IsNextToTarget(point) {
			shortestPath := make([]Point, 10000)

			for _, dest := range m.GetDestinations(point) {
				path, found := m.GetShortestPath(point, dest)
				if found && len(path) < len(shortestPath) {
					shortestPath = path
				}
				if found && len(path) == len(shortestPath) && path[1].FirstThan(shortestPath[1]) {
					shortestPath = path
				}
			}

			if len(shortestPath) == 10000 {
				continue
			}

			moveTo := shortestPath[1]

			// fmt.Printf("Moving from %v to %v\n", point, moveTo)
			m.Move(point, moveTo)
			point = moveTo
		}

		// Atack!
		if m.IsNextToTarget(point) {
			possibleTargets := m.GetTargetsInRange(point)
			target := possibleTargets[0]
			for _, t := range possibleTargets {
				if m.GetHP(t) < m.GetHP(target) {
					target = t
				}
			}
			// fmt.Printf("%v is attacking %v\n", point, target)
			m.Attack(target)
			if m.GetHP(target) == 0 {
				dead.Add(target)
			}
		}
	}

	return true
}

func (m *Map) CountElves() int {
	count := 0
	for _, row := range m.layout {
		for _, entity := range row {
			if entity == Elf {
				count++
			}
		}
	}
	return count
}

func BeverageBandits(m *Map, debug bool) int {
	if debug {
		fmt.Println("Initial State:")
		m.Print()
		fmt.Println()
	}

	round := 0
	for !m.Over() {
		completed := m.Update()
		if completed {
			round++
		}
		if debug {
			fmt.Printf("Round %d\n", round)
			m.Print()

			var val int
			fmt.Print("Press enter to continue...\n")
			fmt.Scanf("%d", &val)
		}
	}

	totalHP := 0
	positions := m.GetEntitiesInReadingOrder()
	for _, pos := range positions {
		totalHP += m.GetHP(pos)
	}

	fmt.Println(round, totalHP)
	return round * totalHP
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

	layout := [][]Entity{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []Entity{}
		for _, r := range scanner.Text() {
			line = append(line, Entity(r))
		}
		layout = append(layout, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for AP := 0; AP < 100; AP++ {
		duplicate := make([][]Entity, len(layout))
		for i := range layout {
			duplicate[i] = make([]Entity, len(layout[i]))
			copy(duplicate[i], layout[i])
		}
		m := NewMap(duplicate, AP)
		initialElves := m.CountElves()
		part2 := BeverageBandits(m, false)
		if initialElves == m.CountElves() {
			fmt.Println("Result:")
			fmt.Println(part2, AP)
			break
		}
	}
}
