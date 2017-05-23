package gosudoku

import (
	"fmt"
	"runtime"
	"testing"
)

func TestStandardSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := SudokuBoard{}
		board = board.Init(STANDARD)
		if unitsLength := len(board.all_units); unitsLength != 27 {
			t.Error("Standard Sudoku board should test 27 units")
		}
	}
}

func TestDiagonalSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := SudokuBoard{}
		board = board.Init(DIAGONAL)
		if unitsLength := len(board.all_units); unitsLength != 29 {
			t.Error("Diagonal Sudoku board should test 29 units")
		}
	}
}

func TestParseSudoku(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("Current test filename: " + filename)
	filename = "samples/800wi.png"
	sudoku_string := parseSudoku(filename)
	fmt.Println("Parsed sudoku: " + sudoku_string)
}
