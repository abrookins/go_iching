package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Hexagram struct {
	// The canonical I Ching number of this hexagram.
	num int

	// The series of boolean values ("lines") that form the hexagram.
	lines [6]bool

	// English name of the hexagram.
	name string

	// The Unicode character of the hexagram.
	character string
}

// getWillhelmUrl returns a URL to the hexagram's entry in a copy of the
// Richard Wilhelm and Cary F. Baynes translation "I Ching: Or, Book of Changes"
// (1950) available.
func (hexagram *Hexagram) getWillhelmUrl() string {
	return fmt.Sprintf("http://www.akirarabelais.com/i/i.html#%v", hexagram.num)
}

// getLeggeUrl returns a URL to the hexagram's entry in a copy of the James
// Legge translation of the I Ching (1899).
func (hexagram *Hexagram) getLeggeUrl() string {
	var fmtString string

	if hexagram.num < 10 {
		fmtString = "%.2d"
	} else {
		fmtString = "%v"
	}

	val := fmt.Sprintf(fmtString, hexagram.num)
	return fmt.Sprintf("http://www.sacred-texts.com/ich/%v.htm", val)
}

var hexagrams [64]Hexagram

var hexagramLookupTable map[[6]bool]Hexagram

func stringsToBools(vals []string) []bool {
	bools := make([]bool, len(vals))

	for i, val := range vals {
		val, err := strconv.ParseBool(val)

		if err != nil {
			fmt.Println("Could not convert value %v to boolean!", val)
		}

		bools[i] = val
	}

	return bools
}

// Load hexagram data from the CSV data file.
func init() {
	f, err := os.Open("data.csv")

	if err != nil {
		panic(fmt.Sprintf("Could not open data.csv"))
	}

	r := csv.NewReader(f)
	rows, err := r.ReadAll()

	if err != nil {
		panic(fmt.Sprintf("Could not read data.csv"))
	}

	for i, row := range rows {
		lines := strings.Split(row[0], "|")
		var linesAsBools [6]bool
		copy(linesAsBools[:], stringsToBools(lines[:6]))
		// Lines are zero-indexed, but the hexagrams begin with `1` in external
		// references.
		hexagrams[i] = Hexagram{i + 1, linesAsBools, row[1], row[2]}
	}

	hexagramLookupTable = make(map[[6]bool]Hexagram)

	for _, hexagram := range hexagrams {
		hexagramLookupTable[hexagram.lines] = hexagram
	}
}

func GetHexagram(lines [6]bool) (Hexagram, bool) {
	hexagram, found := hexagramLookupTable[lines]

	return hexagram, found
}
