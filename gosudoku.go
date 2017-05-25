/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/jamesandersen/gosudoku/sudokuparser"
)

// SudokuMode represents a variant of sudoku
type SudokuMode int

const (
	// STANDARD Each block, row and column should contain digits 1-9
	STANDARD SudokuMode = iota

	// DIAGONAL Diagonals should also contains digits 1-9
	DIAGONAL SudokuMode = iota
)

const (
	// Row labels
	ROWS = "ABCDEFGHI"

	// Column labels
	COLS = "123456789"
)

// A SudokuBoard represents a sudoku board
type SudokuBoard struct {
	values [9][9]string

	rowUnits  [][]string
	colUnits  [][]string
	sqrUnits  [][]string
	diagUnits [][]string
	allUnits  [][]string

	unitsByCell map[string][][]string
}

// BOXES is an array of all possible cells in standard sudoku puzzle
var BOXES = cross(ROWS, COLS)

// Init - initializes a board type based on the mode
func (b SudokuBoard) Init(mode SudokuMode) SudokuBoard {
	for _, row := range ROWS {
		b.rowUnits = append(b.rowUnits, cross(string(row), COLS))
	}

	for _, col := range COLS {
		b.colUnits = append(b.colUnits, cross(ROWS, string(col)))
	}

	for _, rBlk := range [3]string{"ABC", "DEF", "GHI"} {
		for _, cBlk := range [3]string{"123", "456", "789"} {
			b.sqrUnits = append(b.sqrUnits, cross(rBlk, cBlk))
		}
	}

	if mode == DIAGONAL {
		unit := []string{}
		for i := 0; i < len(COLS); i++ {
			unit = append(unit, string(ROWS[i])+string(COLS[i]))
		}
		b.diagUnits = append(b.diagUnits, unit)

		for i := 0; i < len(COLS); i++ {
			unit = append(unit, string(ROWS[i])+string(COLS[len(COLS)-i-1]))
		}
		b.diagUnits = append(b.diagUnits, unit)
	}

	// add all other units into single units list
	for _, u := range [][][]string{b.rowUnits, b.colUnits, b.sqrUnits, b.diagUnits} {
		b.allUnits = append(b.allUnits, u...)
	}

	// initialize map of cells to units which include the cell
	b.unitsByCell = make(map[string][][]string)
	for _, cell := range BOXES {
		unitsForCell := make([][]string, 0)
		for _, unit := range b.allUnits {
			for _, unitCell := range unit {
				if cell == unitCell {
					unitsForCell = append(unitsForCell, unit)
				}
			}
		}
		// set the key/value with all units this cell belongs to
		b.unitsByCell[cell] = unitsForCell
	}

	// initialize map of cells to all peer cells
	b.unitsByCell = make(map[string][][]string)
	for _, cell := range BOXES {
		unitsForCell := make([][]string, 0)
		for _, unit := range b.allUnits {
			for _, unitCell := range unit {
				if cell == unitCell {
					unitsForCell = append(unitsForCell, unit)
				}
			}
		}
		// set the key/value with all units this cell belongs to
		b.unitsByCell[cell] = unitsForCell
	}

	return b
}

// Cross product of elements in a and elements in b.
func cross(a string, b string) []string {
	values := []string{}
	for _, outer := range a {
		for _, inner := range b {
			values = append(values, string(outer)+string(inner))
		}
	}

	return values
}

func main() {
	var filename string

	flag.StringVar(&filename, "filename", "", "Sudoku puzzle image")

	flag.Parse()

	parsed := sudokuparser.ParseSudokuFromFile(filename)

	fmt.Println("Parsed sudoku: " + parsed)
}
