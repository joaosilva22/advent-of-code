package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Step struct {
	ID        string
	dependsOn []*Step
	done      bool
	timer     int
}

func NewStep(ID string) *Step {
	var step Step
	step.ID = ID
	step.dependsOn = []*Step{}
	step.done = false
	step.timer = 0
	return &step
}

func (s Step) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Step %s depends on: ", s.ID))
	for _, v := range s.dependsOn {
		sb.WriteString(v.ID)
		sb.WriteString(" ")
	}
	return sb.String()
}

func (s *Step) DependsOn(o *Step) {
	s.dependsOn = append(s.dependsOn, o)
}

func (s *Step) IsReady() bool {
	if s.done || s.timer > 0 {
		return false
	}
	for _, o := range s.dependsOn {
		if !o.done {
			return false
		}
	}
	return true
}

func (s *Step) Start(baseDuration int) {
	duration := int(s.ID[0]) - 64 + baseDuration
	s.timer = duration
}

func (s *Step) Update() {
	s.timer -= 1
	if s.timer == 0 {
		s.done = true
	}
}

func SumOfItsParts(steps []*Step) string {
	var sb strings.Builder
	done := false
	for !done {
		done = true
		for _, step := range steps {
			if step.IsReady() {
				sb.WriteString(step.ID)
				step.done = true
				done = false
				break
			}
		}
	}
	return sb.String()
}

func SumOfItsPartsPart2(steps []*Step, baseDuration, workerPool int) int {
	inProgress := []*Step{}
	for _, step := range steps {
		if workerPool == 0 {
			break
		}
		if step.IsReady() {
			inProgress = append(inProgress, step)
			step.Start(baseDuration)
			workerPool -= 1
		}
	}
	second := 0
	for len(inProgress) > 0 {
		fmt.Printf("Second=%d len(inProgress)=%d\n", second, len(inProgress))
		second += 1
		toDelete := []string{}
		for _, step := range inProgress {
			step.Update()
			if step.done {
				toDelete = append(toDelete, step.ID)
				workerPool += 1
			}
		}
		for _, ID := range toDelete {
			for i, step := range inProgress {
				if step.ID == ID {
					inProgress = append(inProgress[:i], inProgress[i+1:]...)
					break
				}
			}
		}
		for _, step := range steps {
			if workerPool == 0 {
				break
			}
			if step.IsReady() {
				inProgress = append(inProgress, step)
				step.Start(baseDuration)
				workerPool -= 1
			}
		}
	}
	return second + 1
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

	steps := []*Step{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var from, to string
		format := "Step %s must be finished before step %s can begin."
		fmt.Sscanf(scanner.Text(), format, &from, &to)

		var fromStep, toStep *Step
		for _, step := range steps {
			if step.ID == from {
				fromStep = step
			}
			if step.ID == to {
				toStep = step
			}
			if fromStep != nil && toStep != nil {
				break
			}
		}
		if fromStep == nil {
			fromStep = NewStep(from)
			steps = append(steps, fromStep)
		}
		if toStep == nil {
			toStep = NewStep(to)
			steps = append(steps, toStep)
		}

		toStep.DependsOn(fromStep)
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].ID < steps[j].ID
	})

	//part1 := SumOfItsParts(steps)
	//fmt.Println(part1)

	part2 := SumOfItsPartsPart2(steps, 60, 5)
	fmt.Println(part2)
}
