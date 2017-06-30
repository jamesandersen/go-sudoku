package main

import (
	"reflect"
	"testing"
)

const (
	simple    = "..3.2.6..9..3.5..1..18.64....81.29..7.......8..67.82....26.95..8..2.3..9..5.1.3.."
	harder    = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	diagonal1 = "2.............62....1....7...6..8...3...9...7...6..4...4....8....52.............3"
	diagonal2 = "4.......3..9.........1...7.....1.8.....5.9.....1.2.....3...5.........7..7.......8"
)

func TestStandardSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)
		if unitsLength := len(board.allUnits); unitsLength != 27 {
			t.Error("Standard Sudoku board should test 27 units")
		}
	}
}

func TestDiagonalSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(diagonal1, DIAGONAL)
		if unitsLength := len(board.allUnits); unitsLength != 29 {
			t.Error("Diagonal Sudoku board should test 29 units")
		}
	}
}

func TestUnitsByCell(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)
		if unitsLength := len(board.unitsByCell); unitsLength != len(BOXES) {
			t.Error("Incorrect number of units in unitsByCell")
		}

		unitsA1 := board.unitsByCell["A1"]
		foundRow := false
		foundCol := false
		for i := range unitsA1 {
			if reflect.DeepEqual(unitsA1[i], board.rowUnits[0]) {
				foundRow = true
			}

			if reflect.DeepEqual(unitsA1[i], board.colUnits[0]) {
				foundCol = true
			}
		}

		if !foundRow || !foundCol {
			t.Error("Units for A1 does not include one of either row A or col 1 units")
		}
	}
}

func TestPeersByCell(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)
		if unitsLength := len(board.peersByCell); unitsLength != len(BOXES) {
			t.Error("Incorrect number of units in peersByCell")
		}

		peersA1 := board.peersByCell["A1"]
		foundA9 := false
		foundI1 := false
		for i := range peersA1 {
			if peersA1[i] == "A9" {
				foundA9 = true
			}

			if peersA1[i] == "I1" {
				foundI1 = true
			}
		}

		if !foundA9 || !foundI1 {
			t.Error("Peers for A1 does not include one of either A9 or I1")
		}
	}
}

func TestValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)
		if unitsLength := len(board.values); unitsLength != len(BOXES) {
			t.Error("Incorrect number of values in peersByCell")
		}
	}
}

func TestDeepCopy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)

		board2 := board.Copy()
		board2.values["A1"] = CellValue{Value: "4", Source: RESOLVED}
		if board.values["A1"] == board2.values["A1"] {
			t.Error("Board is not deep copied")
		}
	}
}

func TestSolve(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := NewSudoku(simple, STANDARD)

		resultBoard, success := board.Solve()
		if !success {
			t.Error("Board was not solved")
		} else {
			resultBoard.Print()
		}
	}
}
