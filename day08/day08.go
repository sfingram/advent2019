package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type dimension struct {
	width, height int
}

func (d dimension) size() int {
	return d.width * d.height
}

type image struct {
	bounds dimension
	layers [][]int
}

func loadImage(filename string, bounds dimension) image {
	b, _ := ioutil.ReadFile(filename)
	s := string(b)
	output := image{
		bounds: bounds,
		layers: make([][]int, len(s)/bounds.size()),
	}
	for i, v := range s {
		if i%bounds.size() == 0 {
			output.layers[i/bounds.size()] = make([]int, bounds.size())
		}
		output.layers[i/bounds.size()][i%bounds.size()] = int(v - '0')
	}
	return output
}

func valCount(input []int, val int) int {
	var k int
	for _, v := range input {
		if v == val {
			k++
		}
	}
	return k
}

func imageChecksum(input image) int {

	minZeros := 1 << 16
	var zeroLayer int
	for i := range input.layers {
		nz := valCount(input.layers[i], 0)
		if minZeros > nz {
			minZeros = nz
			zeroLayer = i
		}
	}
	return valCount(input.layers[zeroLayer], 1) * valCount(input.layers[zeroLayer], 2)
}

func renderImage(input image) {
	for row := 0; row < input.bounds.height; row++ {
		for col := 0; col < input.bounds.width; col++ {
			pixel := 2
			for layer := 0; pixel == 2 && layer < len(input.layers); layer++ {
				pixel = input.layers[layer][col+row*input.bounds.width]
			}
			if pixel == 1 {
				fmt.Printf("*")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {

	if len(os.Args) < 4 {
		log.Fatalf("need filename width height")
	}

	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])
	image := loadImage(os.Args[1], dimension{width, height})
	fmt.Printf("Part 1: %d\n", imageChecksum(image))
	fmt.Printf("Part 2: \n")
	renderImage(image)
}
