package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Instruction struct {
	name string
	A    int
	B    int
	C    int
}

type Program struct {
	IP           int
	registers    [6]int
	operations   map[string]func([6]int, int, int, int) [6]int
	instructions []Instruction
}

func NewProgram(IP int, instructions []Instruction) *Program {
	p := Program{}
	p.IP = IP
	p.instructions = instructions
	p.registers = [6]int{}
	p.registers[0] = 1
	p.operations = map[string]func([6]int, int, int, int) [6]int{}
	p.operations["addr"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] + reg[B]
		return reg
	}
	p.operations["addi"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] + B
		return reg
	}
	p.operations["mulr"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] * reg[B]
		return reg
	}
	p.operations["muli"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] * B
		return reg
	}
	p.operations["banr"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] & reg[B]
		return reg
	}
	p.operations["bani"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] & B
		return reg
	}
	p.operations["borr"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] | reg[B]
		return reg
	}
	p.operations["bori"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A] | B
		return reg
	}
	p.operations["setr"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = reg[A]
		return reg
	}
	p.operations["seti"] = func(reg [6]int, A, B, C int) [6]int {
		reg[C] = A
		return reg
	}
	p.operations["gtir"] = func(reg [6]int, A, B, C int) [6]int {
		if A > reg[B] {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	p.operations["gtri"] = func(reg [6]int, A, B, C int) [6]int {
		if reg[A] > B {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	p.operations["gtrr"] = func(reg [6]int, A, B, C int) [6]int {
		if reg[A] > reg[B] {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	p.operations["eqir"] = func(reg [6]int, A, B, C int) [6]int {
		if A == reg[B] {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	p.operations["eqri"] = func(reg [6]int, A, B, C int) [6]int {
		if reg[A] == B {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	p.operations["eqrr"] = func(reg [6]int, A, B, C int) [6]int {
		if reg[A] == reg[B] {
			reg[C] = 1
		} else {
			reg[C] = 0
		}
		return reg
	}
	return &p
}

func (p *Program) Execute(debug bool) [6]int {
	for p.registers[p.IP] < len(p.instructions) {
		i := p.instructions[p.registers[p.IP]]
		op := p.operations[i.name]

		if debug {
			var r, v int
			fmt.Print("Modify Register: ")
			_, err := fmt.Scanf("%d %d", &r, &v)
			if err == nil {
				p.registers[r] = v
			}
		}

		p.registers = op(p.registers, i.A, i.B, i.C)
		p.registers[p.IP] += 1

		if debug {
			fmt.Printf("%v %v\n", i, p.registers)
		}
	}
	return p.registers
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

	IP := 0
	instructions := []Instruction{}

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "#ip %d", &IP)

	for scanner.Scan() {
		var name string
		var A, B, C int
		fmt.Sscanf(scanner.Text(), "%s %d %d %d", &name, &A, &B, &C)

		i := Instruction{name, A, B, C}
		instructions = append(instructions, i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	p := NewProgram(IP, instructions)
	reg := p.Execute(true)
	fmt.Println(reg, reg[0])
}
