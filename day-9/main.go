package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Marble struct {
	val  int
	next *Marble
	prev *Marble
}

func NewMarble(val int) *Marble {
	marble := Marble{}
	marble.val = val
	return &marble
}

func (m *Marble) SetNext(o *Marble) {
	m.next = o
	o.prev = m
}

func (m *Marble) Next(n int) *Marble {
	next := m
	for i := n; i > 0; i-- {
		next = next.next
	}
	return next
}

func (m *Marble) Prev(n int) *Marble {
	prev := m
	for i := n; i > 0; i-- {
		prev = prev.prev
	}
	return prev
}

func MarbleManiaXL(players, marbles int) int {
	scores := make([]int, players+1)

	currentMarble := NewMarble(0)
	currentMarble.SetNext(currentMarble)

	for i := 1; i <= marbles; i++ {
		currentPlayer := ((i - 1) % players) + 1
		if i%23 == 0 {
			scores[currentPlayer] += i + currentMarble.Prev(7).val
			currentMarble.Prev(8).SetNext(currentMarble.Prev(6))
			currentMarble = currentMarble.Prev(6)
		} else {
			marble := NewMarble(i)
			marble.SetNext(currentMarble.Next(2))
			currentMarble.Next(1).SetNext(marble)
			currentMarble = marble
		}
	}

	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	return maxScore
}

func MarbleMania(players, marbles int) int {
	circle := []int{}
	scores := make([]int, players+1)

	currentMarble := 0
	circle = append(circle, 0)

	for marble := 1; marble <= marbles; marble++ {
		currentPlayer := ((marble - 1) % players) + 1
		if marble%23 == 0 {
			remove := ((currentMarble-7)%len(circle) + len(circle)) % len(circle)
			scores[currentPlayer] += marble + circle[remove]
			circle = append(circle[:remove], circle[remove+1:]...)
			currentMarble = remove
		} else {
			currentMarble = (currentMarble+1)%len(circle) + 1
			circle = append(circle, 0)
			copy(circle[currentMarble+1:], circle[currentMarble:])
			circle[currentMarble] = marble
		}
	}

	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	return maxScore
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

	var input string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = scanner.Text()
	}

	var players, marbles int
	format := "%d players; last marble is worth %d points"
	fmt.Sscanf(input, format, &players, &marbles)

	maxScore := MarbleMania(players, marbles)
	fmt.Println(maxScore)

	maxScoreX100 := MarbleManiaXL(players, marbles*100)
	fmt.Println(maxScoreX100)
}
