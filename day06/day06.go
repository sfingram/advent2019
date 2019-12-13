package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const exitError = 1

type body string
type orbit struct {
	parent, satellite body
}

func newOrbit(raw string) orbit {
	components := strings.Split(raw, ")")
	return orbit{
		parent:    body(components[0]),
		satellite: body(components[1]),
	}
}

func loadOrbits(filename string) []orbit {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s", filename)
	}
	defer file.Close()

	orbits := make([]orbit, 0, 1<<10)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		orbits = append(orbits, newOrbit(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		os.Exit(exitError)
	}

	return orbits
}

func makeLookups(orbits []orbit) map[body][]orbit {
	lookup := make(map[body][]orbit)
	for _, v := range orbits {
		o, _ := lookup[v.satellite]
		lookup[v.satellite] = append(o, v)
	}
	return lookup
}

func countOrbit(b body, lookups map[body][]orbit, counts map[body]int) (int, map[body]int) {
	c, ok := counts[b]
	if !ok {
		var subS int
		v := lookups[b]
		c += len(v)
		for _, p := range v {
			subS, counts = countOrbit(p.parent, lookups, counts)
			c += subS
		}
	}
	return c, counts
}

func countOrbits(orbits []orbit, lookups map[body][]orbit) int {

	counts := make(map[body]int)
	var sum, subS int
	for k := range lookups {
		subS, counts = countOrbit(k, lookups, counts)
		sum += subS
	}
	return sum
}

func path(b body, lookups map[body][]orbit) []body {
	p := make([]body, 0, 10)
	o, ok := lookups[b]
	for ok {
		p = append(p, o[0].parent)
		o, ok = lookups[o[0].parent]
	}
	return p
}

func commonAncestor(path1 []body, path2 []body) (body, int) {
	for i1, b1 := range path1 {
		for i2, b2 := range path2 {
			if b1 == b2 {
				return b1, i1 + i2
			}
		}
	}
	return body("ERROR"), -1
}

func main() {

	if len(os.Args) < 2 {
		os.Exit(exitError)
	}

	orbits := loadOrbits(os.Args[1])
	lookups := makeLookups(orbits)
	fmt.Printf("Part 1: %d\n", countOrbits(orbits, lookups))

	_, distance := commonAncestor(path(body("YOU"), lookups), path(body("SAN"), lookups))
	fmt.Printf("Part 2: %d\n", distance)
}
