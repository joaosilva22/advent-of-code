package main

import "fmt"

func PowerLevelXY(x, y, serial int) int {
	rackId := x + 10
	powerLvl := ((((rackId*y + serial) * rackId) / 100) % 10) - 5
	return powerLvl
}

func PowerLevel3X3(x, y, serial int, memo [][]int) int {
	totalPowerLvl := 0
	for row := y; row < y+3; row++ {
		for col := x; col < x+3; col++ {
			if memo[row][col] == 0 {
				memo[row][col] = PowerLevelXY(col, row, serial)
			}
			totalPowerLvl += memo[row][col]
		}
	}
	return totalPowerLvl
}

func PowerLevelNXN(x, y, n, serial int, memo [][]int) int {
	totalPowerLvl := 0
	for row := y; row < y+n; row++ {
		for col := x; col < x+n; col++ {
			if memo[row][col] == 0 {
				memo[row][col] = PowerLevelXY(col, row, serial)
			}
			totalPowerLvl += memo[row][col]
		}
	}
	return totalPowerLvl
}

func ChronalCharge(serial int) (int, int) {
	maxX, maxY := 0, 0
	maxPowerLvl := 0

	memo := make([][]int, 301)
	for y := 0; y <= 300; y++ {
		memo[y] = make([]int, 301)
	}

	for y := 1; y <= 298; y++ {
		for x := 1; x <= 298; x++ {
			powerLvl := PowerLevel3X3(x, y, serial, memo)
			if powerLvl > maxPowerLvl {
				maxPowerLvl = powerLvl
				maxX = x
				maxY = y
			}
		}
	}

	return maxX, maxY
}

func ChronalChargePart2(serial int) (int, int, int) {
	maxX, maxY := 0, 0
	maxPowerLvl, maxSize := 0, 0

	memo := make([][]int, 301)
	for y := 0; y <= 300; y++ {
		memo[y] = make([]int, 301)
	}

	for size := 1; size <= 300; size++ {
		for y := 1; y <= 301-size; y++ {
			for x := 1; x <= 301-size; x++ {
				powerLvl := PowerLevelNXN(x, y, size, serial, memo)
				if powerLvl > maxPowerLvl {
					maxPowerLvl = powerLvl
					maxX = x
					maxY = y
					maxSize = size
				}
			}
		}
	}

	return maxX, maxY, maxSize
}

func main() {
	fmt.Println(ChronalCharge(6878))

	fmt.Println(ChronalChargePart2(6878))
}
