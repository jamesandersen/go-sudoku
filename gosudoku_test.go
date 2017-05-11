package gosudoku

import "testing"

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
