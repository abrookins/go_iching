/*
 * Calculate an I Ching oracle.
 */

package iching

import (
	"math/rand"
	"time"
)

type Line int

type Reading struct {
	// The question asked of the oracle.
	Question string

	// The six generated lines of the reading.
	Lines [6]Line

	// The hexagram the lines form.
	Hexagram *Hexagram

	// The hexagram the lines are changing into (if applicable).
	NextHexagram *Hexagram
}

const (
	// The number of coins to use when randomly generating lines for a hexagram.
	NUM_COINS = 3

	// The number of lines to generate to make a hexagram.
	NUM_LINES = 6
)

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

// getLines returns an array of six lines that form a Hexagram.
func getLines() [6]Line {
	var lines [6]Line

	for i := 0; i < NUM_LINES; i++ {
		lines[i] = getLine()
	}

	return lines
}


func GetReading(question string) *Reading {
	lines := getLines()
	hexagram, _ := GetHexagram(linesToBools(lines))
	nextHexagram, _ := getNextHexagram(lines)

	return &Reading{question, lines, hexagram, nextHexagram}
}
