/*
 * Calculate an I Ching oracle.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// The value of these constants follows the traditional three-coin divination
// technique. See: http://en.wikipedia.org/wiki/I_Ching_divination
const (
	// A broken line changing into a solid line.
	LINE_OLD_YIN = 6
	// A solid line.
	LINE_YOUNG_YANG = 7
	// A broken line.
	LINE_YOUNG_YIN = 8
	// A solid line changing into a broken line.
	LINE_OLD_YANG = 9
)

const (
	// A line that is either new or old yin.
	YIN = false
	// A line that is either new or old yang.
	YANG = true
)

const (
	// The number of coins to use when randomly generating lines for a hexagram.
	NUM_COINS = 3
	// The number of lines to generate to make a hexagram.
	NUM_LINES = 6
)

type Line int

// getLine simulates tossing three coins to generate a random I Ching line,
// following the traditional three-coin toss method of divination.
func getLine() Line {
	var line Line = 0

	for i := 0; i < NUM_COINS; i++ {
		rand.Seed(time.Now().UnixNano())
		var flip = rand.Intn(100)%2 == 0

		// For the origin of these values, see:
		// http://en.wikipedia.org/wiki/I_Ching_divination
		if flip {
			line += 3
		} else {
			line += 2
		}
	}

	return line
}

// getLines returns a six lines that form a Hexagram.
func getLines() [6]Line {
	var lines [6]Line

	for i := 0; i < NUM_LINES; i++ {
		lines[i] = getLine()
	}

	return lines
}

// linesToBools converts an array of Lines that may include "young"
// and "old" yin and yang values into an array of booleans, in which each
// boolean represents whether the line is a yin or yang value.
func linesToBools(lines [6]Line) [6]bool {
	var bools [6]bool

	for i, line := range lines {
		if line == LINE_YOUNG_YANG || line == LINE_OLD_YANG {
			bools[i] = YANG
		} else {
			bools[i] = YIN
		}
	}

	return bools
}

// getNextHexagram checks if any lines in `lines` are changing, and if so,
// find and return the Hexagram the new lines form.
func getNextHexagram(lines [6]Line) (Hexagram, bool) { 
	var nextLines [6]Line

	for i, line := range lines {
		switch {
		case line == LINE_OLD_YANG:
			nextLines[i] = LINE_YOUNG_YIN
		case line == LINE_OLD_YIN:
			nextLines[i] = LINE_YOUNG_YANG
		default:
			nextLines[i] = line	
		}
	}

	return GetHexagram(linesToBools(nextLines))
}

// reportHexagram prints details about a hexagram.
func reportHexagram(hexagram Hexagram) {
	fmt.Printf("#%v: %v (%v)\n", hexagram.num, hexagram.name, hexagram.character)
	fmt.Println(hexagram.getWillhelmUrl())
	fmt.Println(hexagram.getLeggeUrl())
}

func RunOracle() {
	lines := getLines()
	hexagram, _ := GetHexagram(linesToBools(lines))
	nextHexagram, _ := getNextHexagram(lines)

	fmt.Printf("Lines: %v\n", lines)
	reportHexagram(hexagram)

	if nextHexagram != hexagram {
		fmt.Println("\nHexagram is changing. Next hexagram:")
		reportHexagram(nextHexagram)
	}
}
