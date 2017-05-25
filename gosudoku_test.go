package main

import (
	"reflect"
	"testing"
)

func TestStandardSudokuBoardInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := SudokuBoard{}
		board = board.Init(STANDARD)
		if unitsLength := len(board.allUnits); unitsLength != 27 {
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
		if unitsLength := len(board.allUnits); unitsLength != 29 {
			t.Error("Diagonal Sudoku board should test 29 units")
		}
	}
}

func TestUnitsByCell(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	} else {
		board := SudokuBoard{}
		board = board.Init(STANDARD)
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
		board := SudokuBoard{}
		board = board.Init(STANDARD)
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
