package main

import (
	"fmt"
	"os"
	"strconv"
)

func makeDigits(x int) [7]int {
	return [7]int{
		x % 10,
		x / 10 % 10,
		x / 100 % 10,
		x / 1000 % 10,
		x / 10000 % 10,
		x / 100000 % 10,
		x / 1000000 % 10,
	}
}

func valid(x int) (int, int) {
	digits := makeDigits(x)
	counter := make(map[int]int)
	for i := 1; i < 6; i++ {
		if digits[i] > digits[i-1] {
			return 0, 0
		}
		if digits[i] == digits[i-1] {
			counter[digits[i]]++
		}
	}
	var has1, has2 int
	for _, v := range counter {
		if v == 1 {
			has2 = 1
		}
		if v > 0 {
			has1 = 1
		}
	}
	return has1, has2
}

func main() {
	start, _ := strconv.Atoi(os.Args[1])
	stop, _ := strconv.Atoi(os.Args[2])
	var part1, part2 int
	for x := start; x <= stop; x++ {
		has1, has2 := valid(x)
		part1 += has1
		part2 += has2
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}
