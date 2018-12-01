package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ChronalCalibrationPart2(input []int64) int64 {
	var freq int64 = 0
	freqs := make(map[int64]int)
	freqs[freq] = 1

	for {
		for _, v := range input {
			freq += v
			freqs[freq] += 1
			if freqs[freq] == 2 {
				return freq
			}
		}
	}
}

func ChronalCalibration(input []int64) int64 {
	var result int64 = 0
	for _, v := range input {
		result += v
	}
	return result
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

	var input []int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := strconv.ParseInt(scanner.Text(), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		input = append(input, i)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := ChronalCalibration(input)
	fmt.Println(result)

	resultPart2 := ChronalCalibrationPart2(input)
	fmt.Println(resultPart2)
}
