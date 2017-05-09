/*
Package go-sudoku implements a simple library for solving sudoku puzzles.

The syntax of the regular expressions accepted is:

    regexp:
        concatenation { '|' concatenation }
    concatenation:
        { closure }
    closure:
        term [ '*' | '+' | '?' ]
    term:
        '^'
        '$'
        '.'
        character
        '[' [ '^' ] character-ranges ']'
        '(' regexp ')'
*/
package gosudoku

// SudokuMode represents a variant of sudoku
type SudokuMode int

const (
	// STANDARD Each block, row and column should contain digits 1-9
	STANDARD SudokuMode = iota

	// DIAGONAL Diagonals should also contains digits 1-9
	DIAGONAL SudokuMode = iota
)

// A SudokuBoard represents a sudoku board
type SudokuBoard struct {
	values [9][9]string

	units []string
}

// Init - initializes a board type based on the mode
func (b SudokuBoard) Init(mode SudokuMode) SudokuBoard {
	b.units = append(b.units, "foo")

	return b
}
