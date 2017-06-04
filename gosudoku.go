/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

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
	values map[string]string

	rowUnits  [][]string
	colUnits  [][]string
	sqrUnits  [][]string
	diagUnits [][]string
	allUnits  [][]string

	unitsByCell map[string][][]string
	peersByCell map[string][]string
}

// BOXES is an array of all possible cells in standard sudoku puzzle
var BOXES = cross(ROWS, COLS)

// Init - initializes a board type based on the mode
func NewSudoku(rawSudoku string, mode SudokuMode) *SudokuBoard {
	b := new(SudokuBoard)

	if len(rawSudoku) != len(BOXES) {
		panic("raw Sudoku strings should be exactly 81 characters")
	}

	b.values = make(map[string]string)
	for i, val := range rawSudoku {
		if val == '.' {
			b.values[BOXES[i]] = "123456789"
		} else {
			b.values[BOXES[i]] = string(val)
		}
	}

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
	b.peersByCell = make(map[string][]string)
	for _, cell := range BOXES {
		peersForCell := make([]string, 0)
		distinctPeers := make(map[string]bool)
		for _, unit := range b.unitsByCell[cell] {
			for _, unitCell := range unit {
				if _, present := distinctPeers[unitCell]; cell != unitCell && !present {
					peersForCell = append(peersForCell, unitCell)
					distinctPeers[unitCell] = true
				}
			}
		}
		// set the key/value with all peers of this cell
		b.peersByCell[cell] = peersForCell
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

// Solve attempts to solve the Sudoku puzzle and return the result
func (board *SudokuBoard) Solve() (*SudokuBoard, bool) {
	return board.Search()
}

// Search reduces the puzzle using constraint propagation then perform a depth-first search
// of possible box values using a box with the fewest possible options
func (board *SudokuBoard) Search() (*SudokuBoard, bool) {
	reduction := board.ReducePuzzle()

	// bail as the puzzle was not solved
	if !reduction {
		return board, false
	}

	solved := board.CountSolved()
	if solved == len(BOXES) {
		return board, true
	}

	// find cell with fewest remaining possible values
	minLength := 10
	var minKey string
	for key, val := range board.values {
		cellLen := len(val)
		if cellLen > 1 && cellLen < minLength {
			minKey = key
			minLength = len(val)
		}
	}

	// recurse to search for solution
	fmt.Print(fmt.Sprintf("Select cell %s with possible values '%s' for recursion\n", minKey, board.values[minKey]))
	for _, digit := range board.values[minKey] {
		copy := board.Copy()
		copy.values[minKey] = string(digit)
		fmt.Print(fmt.Sprintf("Recursing with %d values solved and picking %s as value for %s\n", solved, string(digit), minKey))
		resultBoard, success := copy.Search()
		if success {
			return resultBoard, true
		}
	}

	return nil, false
}

// Repeatedly propagate constraints until no further progress is made
func (board *SudokuBoard) ReducePuzzle() bool {
	progress := true
	for progress {
		numSolvedBefore := board.CountSolved()
		board.Eliminate()
		board.EliminateNakedTwins()
		board.EliminateOnlyChoice()
		numSolvedAfter := board.CountSolved()
		progress = numSolvedAfter > numSolvedBefore
		if board.Unsolvable() {
			return false
		}
	}

	return true
}

// Eliminate removes solved values from peer boxes within each unit
func (board *SudokuBoard) Eliminate() {
	for cell, val := range board.values {
		if len(val) == 1 {
			for _, peer := range board.peersByCell[cell] {
				newVal := strings.Replace(board.values[peer], val, "", 1)
				/*if newVal == "" {
					panic(fmt.Sprintf("You may have an invalid Sudoku puzzle\nelimination of %s from %s leaves no valid options", val, cell))
				}*/
				board.AssignValue(peer, newVal)
			}
		}
	}
}

// EliminateOnlyChoice finds boxes where a digit is only an option
// in one box within the unit and solve those boxes
func (board *SudokuBoard) EliminateOnlyChoice() {

	for _, unit := range board.allUnits {
		for _, num := range "123456789" {
			matches := make([]string, 0)
			for _, cell := range unit {
				for _, possibleVal := range board.values[cell] {
					if num == possibleVal {
						// the digit we're testing matches a cell in the unit
						matches = append(matches, cell)
					}
				}
			}

			if len(matches) == 1 {
				board.AssignValue(matches[0], string(num))
			}
		}
	}
}

// EliminateNakedTwins eliminates values using the naked twins strategy.
func (board *SudokuBoard) EliminateNakedTwins() bool {

	return false
}

// AssignValue assigns a solution to a cell in the SudokuBoard
func (board *SudokuBoard) AssignValue(cell string, solution string) {
	if board.values[cell] != solution {
		board.values[cell] = solution
	}
}

// CountSolved counts the number of cells which have been solved in the puzzle
func (board *SudokuBoard) CountSolved() int {
	solved := 0
	for _, val := range board.values {
		if len(val) == 1 {
			solved++
		}
	}

	return solved
}

// Unsolvable determines if the puzzle is solvable
func (board *SudokuBoard) Unsolvable() bool {
	for _, val := range board.values {
		if len(val) == 0 {
			return true
		}
	}

	return false
}

// Copy deep copies the sudoku board values
func (board *SudokuBoard) Copy() *SudokuBoard {
	boardCopy := new(SudokuBoard)

	// element copy values
	boardCopy.values = make(map[string]string)
	for key, val := range board.values {
		boardCopy.values[key] = val
	}

	// These collections do not need to be deep copied
	boardCopy.rowUnits = board.rowUnits
	boardCopy.colUnits = board.colUnits
	boardCopy.sqrUnits = board.sqrUnits
	boardCopy.diagUnits = board.diagUnits
	boardCopy.allUnits = board.allUnits
	boardCopy.unitsByCell = board.unitsByCell
	boardCopy.peersByCell = board.peersByCell

	return boardCopy
}

// ToString returns a string represenation of the puzzle
func (board *SudokuBoard) ToString() string {
	str := ""
	for i, cell := range BOXES {
		if i%9 == 0 {
			if i == 0 {
				str += "\n"
			} else {
				str += "|\n"
			}
		}
		if i%27 == 0 {
			str += "+-------+-------+-------+\n"
		}
		if i%3 == 0 {
			str += "| "
		}

		cellValue := board.values[cell]
		if len(cellValue) == 1 {
			str += cellValue + " "
		} else {
			str += ". "
		}
	}

	return str + "|\n+-------+-------+-------+\n"
}

// Print prints the puzzle
func (board *SudokuBoard) Print() {
	fmt.Print(board.ToString())
}

func main() {
	var filename string
	var mode string
	flag.StringVar(&mode, "mode", "serve", "whether to serve web app or parse additional args in CLI mode")
	flag.StringVar(&filename, "filename", "", "Sudoku puzzle image")

	flag.Parse()

	if mode == "serve" {
		fs := http.FileServer(http.Dir("web/static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/solve", solveHandler)
		http.HandleFunc("/", sudokuFormHandler)
		http.ListenAndServe(":8080", nil)
	} else if mode == "cli" {
		parsed := sudokuparser.ParseSudokuFromFile(filename)

		fmt.Println("Parsed sudoku: " + parsed)

		board := NewSudoku(parsed, STANDARD)
		fmt.Print("Attempting to solve Sudoku...\n")
		board.Print()
		finalBoard, success := board.Solve()

		if success {
			finalBoard.Print()
		} else {
			fmt.Print("Board not solved...")
		}
	}
}
