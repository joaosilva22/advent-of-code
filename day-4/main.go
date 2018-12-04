package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type EventKind int

const (
	WakeUp     EventKind = 0
	FallAsleep EventKind = 1
	StartShift EventKind = 2
)

type Event struct {
	id   int
	time time.Time
	kind EventKind
}

func FindSleepyGuard(events []Event) (int, int) {
	times := map[int][]int{}

	var guard int
	var min int

	for _, event := range events {
		switch event.kind {
		case StartShift:
			guard = event.id
			fmt.Printf("Guard %d started it's shift\n", guard)
		case FallAsleep:
			fmt.Printf("Guard %d fell asleep at %v\n", guard, event.time)
			min = event.time.Minute()
		case WakeUp:
			fmt.Printf("Guard %d woke up at %v\n", guard, event.time)
			for i := min; i < event.time.Minute(); i++ {
				if times[guard] == nil {
					times[guard] = make([]int, 60)
				}
				times[guard][i] += 1
			}
		}
	}

	sleepy := -1
	max := 0
	for id, time := range times {
		total := 0
		for _, amount := range time {
			total += amount
		}
		if total > max {
			max = total
			sleepy = id
		}
	}

	fmt.Printf("Guard %d sleeps the most (%d minutes)\n", sleepy, max)

	chosen := 0
	max = 0
	for i, count := range times[sleepy] {
		if count > max {
			chosen = i
			max = count
		}
	}

	// Part 2
	mostFrequentGuard := -1
	mostFrequentMin := -1
	maxMinutesSlept := 0

	for i := 0; i < 60; i++ {
		for guardId, time := range times {
			if time[i] > maxMinutesSlept {
				maxMinutesSlept = time[i]
				mostFrequentGuard = guardId
				mostFrequentMin = i
			}
		}
	}

	fmt.Printf("Chosen minute is %d\n", chosen)

	return chosen * sleepy, mostFrequentGuard * mostFrequentMin
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

	events := []Event{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var year, month, day, hour, min int

		parts := strings.SplitAfterN(scanner.Text(), "]", 2)
		message := strings.TrimSpace(parts[1])

		format := "[%d-%d-%d %d:%d]"
		fmt.Sscanf(parts[0], format, &year, &month, &day, &hour, &min)

		time := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)

		var event Event
		if message == "wakes up" {
			event = Event{id: -1, time: time, kind: WakeUp}
		} else if message == "falls asleep" {
			event = Event{id: -1, time: time, kind: FallAsleep}
		} else {
			var id int
			format := "Guard #%d begins shift"
			fmt.Sscanf(message, format, &id)
			event = Event{id: id, time: time, kind: StartShift}
		}

		events = append(events, event)
	}

	sort.Slice(events, func(a, b int) bool {
		return events[a].time.Before(events[b].time)
	})

	result, result2 := FindSleepyGuard(events)
	fmt.Println(result, result2)
}
