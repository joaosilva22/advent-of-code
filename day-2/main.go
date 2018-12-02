package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Checksum(input []string) int {
	exactlyTwo, exactlyThree := 0, 0
	for _, id := range input {
		freqs := make(map[byte]int)
		for i := 0; i < len(id); i++ {
			freqs[id[i]] += 1
		}
		foundTwo, foundThree := false, false
		for _, freq := range freqs {
			if foundTwo && foundThree {
				break
			}
			if freq == 2 && !foundTwo {
				exactlyTwo += 1
				foundTwo = true
			}
			if freq == 3 && !foundThree {
				exactlyThree += 1
				foundThree = true
			}
		}
	}
	return exactlyTwo * exactlyThree
}

func CommonLetters(input []string) string {
	for i, id1 := range input {
		for _, id2 := range input[i:] {
			index := 0
			differences := 0
			for k := 0; k < len(id1) && k < len(id2); k++ {
				if id1[k] != id2[k] {
					differences += 1
					index = k
				}
			}
			if differences == 1 {
				return id1[:index] + id1[index+1:]
			}
		}
	}
	return "Not found"
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

	var input []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := Checksum(input)
	fmt.Println(result)

	commonLetters := CommonLetters(input)
	fmt.Println(commonLetters)
}
