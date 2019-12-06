package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const exitError = 1

type point struct {
	X, Y int
}

type segment struct {
	Begin, End point
}

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
	return 0, data, bufio.ErrFinalToken
}

func loadWire(wireData string) []segment {
	wire := make([]segment, 0, 1<<8)
	scanner := bufio.NewScanner(strings.NewReader(wireData))
	scanner.Split(commaSplit)
	var x1, y1, x2, y2 int
	for scanner.Scan() {
		var direction string
		var amount int
		fmt.Sscanf(scanner.Text(), "%1s%04d", &direction, &amount)
		switch direction {
		case "U":
			y2 += amount
		case "D":
			y2 -= amount
		case "L":
			x2 -= amount
		case "R":
			x2 += amount
		}
		wire = append(wire, segment{Begin: point{x1, y1}, End: point{x2, y2}})
		x1 = x2
		y1 = y2
	}
	return wire
}

func loadWires(filename string) [][]segment {
	wires := make([][]segment, 0, 2)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wires = append(wires, loadWire(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file %s", filename)
	}

	return wires
}

func (s segment) horizontal() bool {
	return s.Begin.Y == s.End.Y
}

func (s segment) corrected() segment {
	if s.End.X < s.Begin.X || s.End.Y < s.Begin.Y {
		return segment{s.End, s.Begin}
	}
	return s
}

// hIntersects assumes s is horizontal and s2 isn't
func (s segment) hIntersects(s2 segment) bool {
	c1 := s.corrected()
	c2 := s2.corrected()
	return c1.Begin.X < c2.Begin.X &&
		c1.End.X > c2.End.X &&
		c2.Begin.Y < c1.Begin.Y &&
		c2.End.Y > c1.End.Y
}

// gets all the intersections of the multiple wires
func wireIntersections(wires [][]segment) []point {

	intersections := make([]point, 0, 1<<8)
	for i := 0; i < len(wires); i++ {
		for j := i + 1; j < len(wires); j++ {
			for _, iWire := range wires[i] {
				for _, jWire := range wires[j] {
					iHorizontal := iWire.horizontal()
					jHorizontal := jWire.horizontal()
					if iHorizontal != jHorizontal {
						if iHorizontal && iWire.hIntersects(jWire) {
							intersections = append(intersections, point{jWire.Begin.X, iWire.Begin.Y})
						} else if jHorizontal && jWire.hIntersects(iWire) {
							intersections = append(intersections, point{iWire.Begin.X, jWire.Begin.Y})
						}
					}
				}
			}
		}
	}
	return intersections
}

// computes the manhattan distance of a point
func (p point) distance() int {
	switch {
	case p.X < 0 && p.Y < 0:
		return -(p.X + p.Y)
	case p.X >= 0 && p.Y < 0:
		return -p.Y + p.X
	case p.X < 0 && p.Y >= 0:
		return -p.X + p.Y
	default:
		return p.X + p.Y
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func intersectDistances(wires []segment, yMap map[int][]point, xMap map[int][]point) map[point]int {
	var traversal int
	distances := make(map[point]int)
	for _, wire := range wires {
		c := wire.corrected()
		if wire.horizontal() {
			pts, ok := yMap[wire.Begin.Y]
			if ok {
				for _, pt := range pts {
					if c.Begin.X < pt.X && c.End.X > pt.X {
						distances[pt] = traversal + abs(pt.X-wire.Begin.X)
					}
				}
			}
			traversal += c.End.X - c.Begin.X
		}
		if !wire.horizontal() {
			pts, ok := xMap[wire.Begin.X]
			if ok {
				for _, pt := range pts {
					if c.Begin.Y < pt.Y && c.End.Y > pt.Y {
						distances[pt] = traversal + abs(pt.Y-wire.Begin.Y)
					}
				}
			}
			traversal += c.End.Y - c.Begin.Y
		}
	}
	return distances
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	wires := loadWires(os.Args[1])
	intersections := wireIntersections(wires)
	smallestDistance := intersections[0].distance()

	xMap := make(map[int][]point)
	yMap := make(map[int][]point)
	for _, v := range intersections {

		// compute smallest distance
		distance := v.distance()
		if smallestDistance > distance {
			smallestDistance = distance
		}

		// fill lookup maps
		xMap[v.X] = append(xMap[v.X], v)
		yMap[v.Y] = append(yMap[v.Y], v)
	}
	fmt.Printf("Part 1: %d\n", smallestDistance)

	iDistances := intersectDistances(wires[0], yMap, xMap)
	jDistances := intersectDistances(wires[1], yMap, xMap)
	shortestIntersection := iDistances[intersections[0]] + jDistances[intersections[0]]
	for _, v := range intersections {
		interDist := iDistances[v] + jDistances[v]
		if shortestIntersection > interDist {
			shortestIntersection = interDist
		}
	}
	fmt.Printf("Part 2: %d\n", shortestIntersection)
}
