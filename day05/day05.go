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

func decode(instruction int) (opcode int, modes [3]int) {
	opcode = instruction % 100
	modes[0] = instruction / 100 % 10
	modes[1] = instruction / 1000 % 10
	modes[2] = instruction / 10000 % 10
	return opcode, modes
}

func param(val int, mode int, data []int) int {
	if mode == 0 {
		return data[val]
	}
	return val
}

func executeProgram(programData []int, input []int) []int {

	output := make([]int, 0)
	for i := 0; i < len(programData); i++ {
		opcode, modes := decode(programData[i])
		switch opcode {
		case 1: // ADD
			programData[programData[i+3]] = param(programData[i+1], modes[0], programData) + param(programData[i+2], modes[1], programData)
			i += 3
		case 2: // MUL
			programData[programData[i+3]] = param(programData[i+1], modes[0], programData) * param(programData[i+2], modes[1], programData)
			i += 3
		case 3: // INP
			programData[programData[i+1]] = input[0]
			input = input[1:]
			i++
		case 4: // OUTP
			output = append(output, param(programData[i+1], modes[0], programData))
			i++
		case 5: // JNZ
			if param(programData[i+1], modes[0], programData) != 0 {
				i = param(programData[i+2], modes[1], programData) - 1
			} else {
				i += 2
			}
		case 6: // JZ
			if param(programData[i+1], modes[0], programData) == 0 {
				i = param(programData[i+2], modes[1], programData) - 1
			} else {
				i += 2
			}
		case 7: // LT
			if param(programData[i+1], modes[0], programData) < param(programData[i+2], modes[1], programData) {
				programData[programData[i+3]] = 1
			} else {
				programData[programData[i+3]] = 0
			}
			i += 3
		case 8: // EQ
			if param(programData[i+1], modes[0], programData) == param(programData[i+2], modes[1], programData) {
				programData[programData[i+3]] = 1
			} else {
				programData[programData[i+3]] = 0
			}
			i += 3
		case 99: // EXT
			return output
		default:
			log.Printf("Error token at %d: %d", i, programData[i])
		}
	}
	return output
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	filename := os.Args[1]
	originalProgram := loadProgram(filename)
	output := executeProgram(append([]int{}, originalProgram...), []int{1})

	fmt.Printf("Part 1: %+v\n", output)

	output = executeProgram(append([]int{}, originalProgram...), []int{5})
	fmt.Printf("Part 2: %+v\n", output)
}
