package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const exitError = 1
const goalState = 19690720
const gridSize = 100

// Define a split function that separates on commas. (stolen from https://golang.org/src/bufio/example_test.go)
func commaSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	// There is one final token to be delivered, which may be the empty string.
	// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
	// but does not trigger an error to be returned from Scan itself.
	return 0, data, bufio.ErrFinalToken
}

func loadProgram(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s", filename)
	}
	defer file.Close()

	data := make([]int, 0, 1<<10)
	scanner := bufio.NewScanner(file)
	scanner.Split(commaSplit)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		data = append(data, val)
	}

	if err := scanner.Err(); err != nil {
		os.Exit(exitError)
	}

	return data
}

func fixProgram(noun int, verb int, programData []int) []int {
	programData[1] = noun
	programData[2] = verb
	return programData
}

func executeProgram(programData []int) []int {

	for i := 0; i < len(programData); i++ {
		switch programData[i] {
		case 1:
			// log.Printf("%03d : ADD %03d %03d %03d", i, programData[i+1], programData[i+2], programData[i+3])
			programData[programData[i+3]] = programData[programData[i+1]] + programData[programData[i+2]]
			i += 3
		case 2:
			// log.Printf("%03d : MUL %03d %03d %03d", i, programData[i+1], programData[i+2], programData[i+3])
			programData[programData[i+3]] = programData[programData[i+1]] * programData[programData[i+2]]
			i += 3
		case 99:
			// log.Printf("%03d : EXIT", i)
			return programData
		default:
			log.Printf("Error token at %d: %d", i, programData[i])
		}
	}
	return programData
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	filename := os.Args[1]
	noun := 12
	verb := 2
	originalProgram := loadProgram(filename)
	programData := make([]int, len(originalProgram))
	copy(programData, originalProgram)
	programData = executeProgram(fixProgram(noun, verb, programData))

	fmt.Printf("Part 1: %d\n", programData[0])

	// Grid search for answer

	for noun = 0; noun < gridSize; noun++ {
		for verb = 0; verb < gridSize; verb++ {

			copy(programData, originalProgram)
			programData = executeProgram(fixProgram(noun, verb, programData))

			if programData[0] == goalState {
				fmt.Printf("Part 2: %d \n", 100*noun+verb)
				os.Exit(0)
			}
		}
	}
}
