package sudokusolver

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

func TestParseSudokuFromFile(t *testing.T) {
	const sample800wi = "7,,,,,3,,,2,,,4,,,,1,,9,,,5,2,,9,,,,,2,,,1,5,,7,,,,,,,,,,,,9,,4,7,,,8,,,,,7,,4,8,,,3,,2,,,,5,,,9,,,3,,,,,1"
	const sample800wiFile = "../samples/800wi.png"
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("Current test filename: " + filename)

	if sudokuString := ParseSudokuFromFile(sample800wiFile); sudokuString != sample800wi {
		t.Error(sample800wiFile + " not parsed as " + sample800wi)
	}
}
