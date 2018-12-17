package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Operation struct {
	name string
	op   func(int, int, int)
}

func NewOperation(name string, op func(int, int, int)) *Operation {
	operation := Operation{name, op}
	return &operation
}

type Sample struct {
	before [4]int
	cmd    [4]int
	after  [4]int
}

func NewSample(before, cmd, after [4]int) Sample {
	s := Sample{before, cmd, after}
	return s
}

func ChronalClassification(samples []Sample) int {
	count := map[Sample]int{}

	for _, sample := range samples {
		fmt.Println(sample)
		count[sample] = 0
		m := NewMachine(sample.before)
		for _, op := range m.operations {
			op.op(sample.cmd[1], sample.cmd[2], sample.cmd[3])
			if m.registers == sample.after {
				count[sample] += 1
				// fmt.Println(count[sample])
			}
			m.SetState(sample.before)
		}
	}

	res := 0
	for sample, cnt := range count {
		if cnt >= 3 {
			// fmt.Println(cnt)
			res++
		} else if cnt == 1 {
			fmt.Println(sample)
		}
	}
	return res
}

type Set struct {
	contents []string
}

func NewSet() *Set {
	s := Set{[]string{}}
	return &s
}

func (s *Set) Contains(o string) bool {
	for _, i := range s.contents {
		if i == o {
			return true
		}
	}
	return false
}

func (s *Set) Add(o string) {
	if !s.Contains(o) {
		s.contents = append(s.contents, o)
	}
}

func (s *Set) Len() int {
	return len(s.contents)
}

func (s *Set) Remove(o string) {
	for i, j := range s.contents {
		if j == o {
			s.contents = append(s.contents[:i], s.contents[i+1:]...)
		}
	}
}

func (s *Set) Copy() *Set {
	contents := make([]string, len(s.contents))
	copy(contents, s.contents)
	copy := Set{contents}
	return &copy
}

func ChronalClassificationPart2(samples []Sample, program [][4]int) int {
	possibleOps := map[int]*Set{}
	for i := 0; i < 16; i++ {
		possibleOps[i] = NewSet()
		possibleOps[i].Add("addr")
		possibleOps[i].Add("addi")
		possibleOps[i].Add("mulr")
		possibleOps[i].Add("muli")
		possibleOps[i].Add("banr")
		possibleOps[i].Add("bani")
		possibleOps[i].Add("borr")
		possibleOps[i].Add("bori")
		possibleOps[i].Add("setr")
		possibleOps[i].Add("seti")
		possibleOps[i].Add("gtir")
		possibleOps[i].Add("gtri")
		possibleOps[i].Add("gtrr")
		possibleOps[i].Add("eqir")
		possibleOps[i].Add("eqri")
		possibleOps[i].Add("eqrr")
	}

	fmt.Println("Initial restrictions: ")
	for code, set := range possibleOps {
		fmt.Println(code, set)
	}

	for _, sample := range samples {
		m := NewMachine(sample.before)
		noMatch := []string{}
		for _, op := range m.operations {
			op.op(sample.cmd[1], sample.cmd[2], sample.cmd[3])
			if m.registers != sample.after {
				// fmt.Println("NOMATCH", op.name)
				noMatch = append(noMatch, op.name)
			}
			m.SetState(sample.before)
		}
		if len(noMatch) < len(m.operations) {
			for _, opName := range noMatch {
				possibleOps[sample.cmd[0]].Remove(opName)
			}
		}
	}

	// Clean the restrictions
	// for code, mnemonics := range possibleOps {
	// 	if mnemonics.Len() == 1 {
	// 		mnemonic := mnemonics.contents[0]
	// 		for otherCode, otherMnemonics := range possibleOps {
	// 			if otherCode == code {
	// 				continue
	// 			}
	// 			if otherMnemonics.Contains(mnemonic) {
	// 				otherMnemonics.Remove(mnemonic)
	// 			}
	// 		}
	// 	}
	// }

	fmt.Println("Initial restrictions: ")
	for code, set := range possibleOps {
		fmt.Println(code, set)
	}

	_, possibleOps = Backtrack(possibleOps)

	fmt.Println("Final restrictions: ")
	for code, set := range possibleOps {
		fmt.Println(code, set)
	}

	m := NewMachine([4]int{0, 0, 0, 0})

	opsMap := map[int]string{}
	for k, v := range possibleOps {
		opsMap[k] = v.contents[0]
	}

	for _, cmd := range program {
		op := opsMap[cmd[0]]
		m.Execute(op, cmd[1], cmd[2], cmd[3])
	}

	return m.registers[0]
}

func Backtrack(restrictions map[int]*Set) (bool, map[int]*Set) {
	fmt.Println("Restrictions")
	for k, v := range restrictions {
		fmt.Println(k, v)
	}

	// Clean the restrictions
	for code, mnemonics := range restrictions {
		if mnemonics.Len() == 1 {
			mnemonic := mnemonics.contents[0]
			for otherCode, otherMnemonics := range restrictions {
				if otherCode == code {
					continue
				}
				if otherMnemonics.Contains(mnemonic) {
					otherMnemonics.Remove(mnemonic)
				}
			}
		}
	}

	smallestCode := -1
	smallest := NewSet()
	for code, s := range restrictions {
		// fmt.Println(code, s)
		if s.Len() == 0 {
			fmt.Printf("Could not find a mnemonic for code %d\n", code)
			return false, nil
		}
		if s.Len() > 1 {
			if smallest.Len() == 0 {
				smallestCode = code
				smallest = s
			} else if s.Len() < smallest.Len() {
				smallestCode = code
				smallest = s
			}
		}
	}

	if smallest.Len() == 0 {
		fmt.Println("Found the solution:")
		for k, v := range restrictions {
			fmt.Println(k, v)
		}
		return true, restrictions
	}

	fmt.Printf("Found the smallest set to be %d: %v\n", smallestCode, smallest)

	for _, mnemonic := range smallest.contents {
		// Create a copy of the restrictions and try that one
		newRestrictions := map[int]*Set{}
		for k, v := range restrictions {
			// Remove the mnemonic from all the other restrictions, except the smallest
			newSet := v.Copy()
			if k == smallestCode {
				newSet = NewSet()
				newSet.Add(mnemonic)
			} else {
				newSet.Remove(mnemonic)
			}
			newRestrictions[k] = newSet
		}
		// Try with the new restrictions
		fmt.Printf("Trying to set %d to %s\n", smallestCode, mnemonic)
		solved, r := Backtrack(newRestrictions)
		if solved {
			return true, r
		}
	}

	return false, nil
}

type Machine struct {
	operations []*Operation
	registers  [4]int
}

func NewMachine(initialState [4]int) *Machine {
	m := Machine{}
	operations := []*Operation{}

	addr := NewOperation("addr", func(A, B, C int) {
		m.registers[C] = m.registers[A] + m.registers[B]
	})
	operations = append(operations, addr)

	addi := NewOperation("addi", func(A, B, C int) {
		m.registers[C] = m.registers[A] + B
	})
	operations = append(operations, addi)

	mulr := NewOperation("mulr", func(A, B, C int) {
		m.registers[C] = m.registers[A] * m.registers[B]
	})
	operations = append(operations, mulr)

	muli := NewOperation("muli", func(A, B, C int) {
		m.registers[C] = m.registers[A] * B
	})
	operations = append(operations, muli)

	banr := NewOperation("banr", func(A, B, C int) {
		m.registers[C] = (m.registers[A] & m.registers[B])
	})
	operations = append(operations, banr)

	bani := NewOperation("bani", func(A, B, C int) {
		m.registers[C] = (m.registers[A] & B)
	})
	operations = append(operations, bani)

	borr := NewOperation("borr", func(A, B, C int) {
		m.registers[C] = (m.registers[A] | m.registers[B])
	})
	operations = append(operations, borr)

	bori := NewOperation("bori", func(A, B, C int) {
		m.registers[C] = (m.registers[A] | B)
	})
	operations = append(operations, bori)

	setr := NewOperation("setr", func(A, B, C int) {
		m.registers[C] = m.registers[A]
	})
	operations = append(operations, setr)

	seti := NewOperation("seti", func(A, B, C int) {
		m.registers[C] = A
	})
	operations = append(operations, seti)

	gtir := NewOperation("gtir", func(A, B, C int) {
		if A > m.registers[B] {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, gtir)

	gtri := NewOperation("gtri", func(A, B, C int) {
		if m.registers[A] > B {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, gtri)

	gtrr := NewOperation("gtrr", func(A, B, C int) {
		if m.registers[A] > m.registers[B] {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, gtrr)

	eqir := NewOperation("eqir", func(A, B, C int) {
		if A == m.registers[B] {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, eqir)

	eqri := NewOperation("eqri", func(A, B, C int) {
		if m.registers[A] == B {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, eqri)

	eqrr := NewOperation("eqrr", func(A, B, C int) {
		if m.registers[A] == m.registers[B] {
			m.registers[C] = 1
			return
		}
		m.registers[C] = 0
	})
	operations = append(operations, eqrr)

	var registers [4]int
	copy(registers[:], initialState[:])
	m.registers = registers
	m.operations = operations
	return &m
}

func (m *Machine) SetState(state [4]int) {
	var registers [4]int
	copy(registers[:], state[:])
	m.registers = registers
}

func (m *Machine) Execute(mnemonic string, A, B, C int) {
	for _, op := range m.operations {
		if op.name == mnemonic {
			op.op(A, B, C)
			return
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

	samples := []Sample{}
	program := [][4]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		before := [4]int{}
		n, _ := fmt.Sscanf(scanner.Text(), "Before: [%d, %d, %d, %d]", &before[0], &before[1], &before[2], &before[3])
		if n > 0 {
			scanner.Scan()
			sampleCmd := [4]int{}
			fmt.Sscanf(scanner.Text(), "%d %d %d %d", &sampleCmd[0], &sampleCmd[1], &sampleCmd[2], &sampleCmd[3])

			scanner.Scan()
			after := [4]int{}
			fmt.Sscanf(scanner.Text(), "After:  [%d, %d, %d, %d]", &after[0], &after[1], &after[2], &after[3])

			sample := NewSample(before, sampleCmd, after)
			samples = append(samples, sample)
		} else {
			cmd := [4]int{}
			fmt.Sscanf(scanner.Text(), "%d %d %d %d", &cmd[0], &cmd[1], &cmd[2], &cmd[3])
			program = append(program, cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1 := ChronalClassification(samples)
	fmt.Println(part1)

	part2 := ChronalClassificationPart2(samples, program)
	fmt.Println(part2)
}
