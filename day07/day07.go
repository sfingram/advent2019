package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
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

func executeProgramChannel(
	name string,
	programData []int,
	input chan int,
	output chan int,
	wg *sync.WaitGroup) {

	defer close(output)

	outp := make([]int, 0)
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
			// log.Printf("%s INP recv", name)
			programData[programData[i+1]] = <-input
			// log.Printf("%s INP %d", name, programData[programData[i+1]])
			i++
		case 4: // OUTP
			v := param(programData[i+1], modes[0], programData)
			outp = append(outp, v)
			// log.Printf("%s OUTP send %d", name, v)
			output <- v
			// log.Printf("%s OUTP wrote %d", name, v)
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
			// log.Printf("%s EXT", name)
			if wg != nil {
				wg.Done()
			}
			return
		default:
			log.Printf("Error token at %d: %d", i, programData[i])
		}
	}
}

// Perm calls f with each permutation of a.
func Perm(a []int, ch chan []int) {
	perm(a, ch, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, ch chan []int, i int) {
	if i > len(a) {
		ch <- append([]int{}, a...)
		return
	}
	perm(a, ch, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, ch, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	filename := os.Args[1]
	originalProgram := loadProgram(filename)

	ch := make(chan []int)
	go func() {
		Perm([]int{0, 1, 2, 3, 4}, ch)
		close(ch)
	}()

	var max int
	for p := range ch {
		output := []int{0}
		for _, v := range p {
			output = executeProgram(append([]int{}, originalProgram...), []int{v, output[0]})
		}
		if max < output[0] {
			max = output[0]
		}
	}
	fmt.Printf("Part 1: %d\n", max)

	// part Two: channels galore

	ch = make(chan []int)
	phases := []int{9, 8, 7, 6, 5}
	go func() {
		Perm(phases, ch)
		close(ch)
	}()

	max = 0
	names := [5]string{"a", "b", "c", "d", "e"}
	for p := range ch {
		var wg sync.WaitGroup
		wg.Add(len(phases) - 1)
		pipe := make([]chan int, 0, len(p))
		for range p {
			pipe = append(pipe, make(chan int))
		}
		for i := range p {
			if i < len(p)-1 {
				go executeProgramChannel(names[i], append([]int{}, originalProgram...), pipe[i], pipe[(i+1)%len(p)], &wg)
			} else {
				go executeProgramChannel(names[i], append([]int{}, originalProgram...), pipe[i], pipe[(i+1)%len(p)], nil)
			}
		}
		for i, v := range p {
			pipe[i] <- v // Provide each amplifier its phase setting at its first input instruction
		}
		pipe[0] <- 0        // To start the process, a 0 signal is sent to amplifier A's input exactly once
		wg.Wait()           // wait for everyone but the last amp
		answer := <-pipe[0] // read the final output from that amp
		if max < answer {
			max = answer
		}
	}
	fmt.Printf("Part 2: %d\n", max)
}
