package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Cart struct {
	x    int
	y    int
	dir  Direction
	step int
}

func NewCart(x, y int, dir Direction) *Cart {
	cart := Cart{x, y, dir, 0}
	return &cart
}

func (c *Cart) update(track [][]rune) {
	switch c.dir {
	case Left:
		c.x -= 1
	case Right:
		c.x += 1
	case Up:
		c.y -= 1
	case Down:
		c.y += 1
	}

	switch track[c.y][c.x] {
	case '\\':
		if c.dir == Right {
			c.dir = Down
		} else if c.dir == Up {
			c.dir = Left
		} else if c.dir == Left {
			c.dir = Up
		} else {
			c.dir = Right
		}
	case '/':
		if c.dir == Up {
			c.dir = Right
		} else if c.dir == Down {
			c.dir = Left
		} else if c.dir == Right {
			c.dir = Up
		} else {
			c.dir = Down
		}
	case '+':
		if c.step == 0 {
			if c.dir == Left {
				c.dir = Down
			} else if c.dir == Right {
				c.dir = Up
			} else if c.dir == Down {
				c.dir = Right
			} else {
				c.dir = Left
			}
		} else if c.step == 2 {
			if c.dir == Left {
				c.dir = Up
			} else if c.dir == Right {
				c.dir = Down
			} else if c.dir == Down {
				c.dir = Left
			} else {
				c.dir = Right
			}
		}
		c.step = (c.step + 1) % 3
	}
}

func (c *Cart) toRune() rune {
	switch c.dir {
	case Left:
		return '<'
	case Right:
		return '>'
	case Up:
		return '^'
	case Down:
		return 'v'
	}
	return '?'
}

type Track struct {
	layout [][]rune
	carts  []*Cart
}

func NewTrack(layout [][]rune, carts []*Cart) *Track {
	track := Track{layout, carts}
	return &track
}

func (t *Track) findFirstCollision(debug bool) (int, int) {
	sort.Slice(t.carts, func(i, j int) bool {
		return t.carts[i].y < t.carts[j].y
	})

	x, y := -1, -1
	done := false

	for !done {
		for i, curr := range t.carts {
			if debug {
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
			}
			curr.update(t.layout)
			for j, other := range t.carts {
				if i == j {
					continue
				}
				if curr.x == other.x && curr.y == other.y {
					x = curr.x
					y = curr.y
					done = true
					break
				}
			}
			if done {
				break
			}
			if debug {
				t.display()
			}
		}
	}
	return x, y
}

func (t *Track) findLastStanding() (int, int) {
	x, y := -1, -1
	done := false
	iter := 0
	for !done {
		iter += 1
		sort.Slice(t.carts, func(i, j int) bool {
			return t.carts[i].y < t.carts[j].y
		})

		toRemove := []int{}
		for i, curr := range t.carts {
			curr.update(t.layout)
			for j, other := range t.carts {
				if i == j {
					continue
				}
				if curr.x == other.x && curr.y == other.y {
					fmt.Println("crash", curr.x, curr.y)
					toRemove = append(toRemove, i)
					toRemove = append(toRemove, j)
				}
			}
		}
		sort.Slice(toRemove, func(i, j int) bool {
			return toRemove[i] > toRemove[j]
		})
		for _, r := range toRemove {
			t.carts = append(t.carts[:r], t.carts[r+1:]...)
		}
		if len(t.carts) == 1 {
			fmt.Println(iter)
			done = true
			x = t.carts[0].x
			y = t.carts[0].y
		}
	}
	return x, y
}

func (t *Track) display() {
	for row, line := range t.layout {
		for col, val := range line {
			isCart := false
			for _, cart := range t.carts {
				if cart.x == col && cart.y == row {
					fmt.Printf("%c", cart.toRune())
					isCart = true
					break
				}
			}
			if !isCart {
				fmt.Printf("%c", val)
			}
		}
		fmt.Println()
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

	layout := [][]rune{}
	carts := []*Cart{}

	row := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		layout = append(layout, []rune{})
		data := scanner.Text()
		for col, c := range data {
			switch c {
			case 'v':
				cart := NewCart(col, row, Down)
				carts = append(carts, cart)
				layout[row] = append(layout[row], '|')
			case '^':
				cart := NewCart(col, row, Up)
				carts = append(carts, cart)
				layout[row] = append(layout[row], '|')
			case '>':
				cart := NewCart(col, row, Right)
				carts = append(carts, cart)
				layout[row] = append(layout[row], '-')
			case '<':
				cart := NewCart(col, row, Left)
				carts = append(carts, cart)
				layout[row] = append(layout[row], '-')
			default:
				layout[row] = append(layout[row], c)
			}
		}
		row++
	}

	track := NewTrack(layout, carts)
	// part1X, part1Y := track.findFirstCollision(false)
	// fmt.Println(part1X, part1Y)

	part2X, part2Y := track.findLastStanding()
	fmt.Println(part2X, part2Y)
}
