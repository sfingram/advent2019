package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const exitError = 1

// fuel calculation
func fuel(mass int) int {
	return mass/3 - 2
}

// allFuel calculation with iterative exhaustion
func allFuel(mass int) int {
	var totalFuel int
	subFuel := fuel(mass)
	for subFuel > 0 {
		totalFuel += subFuel
		subFuel = fuel(subFuel)
	}
	return totalFuel
}

// fileFuel open a file and sum all fuel calculations
func fileFuel(filename string, fuelCalc func(int) int) int {
	file, err := os.Open(filename)
	if err != nil {
		os.Exit(exitError)
	}
	defer file.Close()

	// Part 1 : sum the mass of the input lines

	totalFuel := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mass, _ := strconv.Atoi(scanner.Text())
		totalFuel += fuelCalc(mass)
	}

	if err := scanner.Err(); err != nil {
		os.Exit(exitError)
	}

	return totalFuel
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	fmt.Printf("Part 1: %d\n", fileFuel(os.Args[1], fuel))
	fmt.Printf("Part 2: %d\n", fileFuel(os.Args[1], allFuel))
}
