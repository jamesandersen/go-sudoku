package gosudoku

import "testing"

func TestSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := SudokuBoard{}
		board = board.Init(STANDARD)
		if unitsLength := len(board.units); unitsLength != 1 {
			t.Error("Should have single value")
		}
	}

}
