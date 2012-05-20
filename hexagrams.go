package iching

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Hexagram struct {
	// The canonical I Ching number of this hexagram.
	Num int

	// The series of boolean values ("lines") that form the hexagram.
	Lines [6]bool

	// English name of the hexagram.
	Name string

	// The Unicode character of the hexagram.
	Character string

	// Description of the hexagram from Wikipedia.
	Description string

	// URLs to Internet translations for this hexagram.
	TranslationUrls []string
}

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


// getWillhelmUrl returns a URL to the entry for hexagram with `hexagramNum` in
// an internet copy of the Richard Wilhelm and Cary F. Baynes translation "I
// Ching: Or, Book of Changes" (1950).
func getWillhelmUrl(hexagramNum int) string {
	return fmt.Sprintf("http://www.akirarabelais.com/i/i.html#%v", hexagramNum)
}

// getLeggeUrl returns a URL to the entry for hexagram with `hexagramNum` in an
// internet copy of the James Legge translation of the I Ching (1899).
func getLeggeUrl(hexagramNum int) string {
	var fmtString string

	if hexagramNum < 10 {
		fmtString = "%.2d"
	} else {
		fmtString = "%v"
	}

	val := fmt.Sprintf(fmtString, hexagramNum)
	return fmt.Sprintf("http://www.sacred-texts.com/ich/ic%v.htm", val)
}

// getTranslationUrls returns a slice of string URLs to internet translations
// for the hexagram identified by `hexagramNum`.
//
// TODO: There has to be a way to do this using a slice or array of function
// types (or pointers), and then iterate over the slice and run the functions,
// but I can't figure it out yet in Go.
func getTranslationUrls(hexagramNum int) []string {
	return []string{getWillhelmUrl(hexagramNum), getLeggeUrl(hexagramNum)}
}

var hexagrams [64]Hexagram

var hexagramLookupTable map[[6]bool]Hexagram

// Load hexagram data from the CSV data file.
func init() {
	_, filename, _, _ := runtime.Caller(1)

	f, err := os.Open(path.Join(path.Dir(filename), "data.csv"))

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

		// Convert lines slice to an array.
		copy(linesAsBools[:], stringsToBools(lines[:6]))

		// Lines are zero-indexed and ordered by hexagram, but the order of hexagrams
		// begins with `1` in external references.
		hexagrams[i] = Hexagram{i + 1, linesAsBools, row[1], row[2], row[3], getTranslationUrls(i)}
	}

	hexagramLookupTable = make(map[[6]bool]Hexagram)

	for _, hexagram := range hexagrams {
		hexagramLookupTable[hexagram.Lines] = hexagram
	}
}

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
func getNextHexagram(lines [6]Line) (*Hexagram, bool) { 
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

// GetAllHexagrams returns a pointer to `hexagrams`.
func GetAllHexagrams() *[64]Hexagram {
	return &hexagrams
}

// GetHexagramByNum returns the hexagram whose number is `num`. Hexagrams are
// traditionally ordered by specific numbers and they appear in the `hexagrams`
// array in this order.
func GetHexagramByNum(num int) (*Hexagram, bool) {
	// Hexagram numbers are one-based, while our array of hexagrams is zero-based.
	num = num - 1
	numHexagrams := len(hexagrams) - 1
	var hexagram Hexagram

	if num <= numHexagrams && num >= 0 {
		hexagram := hexagrams[num]
		return &hexagram, true
	}

	return &hexagram, false
}

// GetHexagram returns the hexagram whose series of "yin" and "yang" lines
// matches the given array of booleans in `lines`.
func GetHexagram(lines [6]bool) (*Hexagram, bool) {
	hexagram, found := hexagramLookupTable[lines]
	return &hexagram, found
}
