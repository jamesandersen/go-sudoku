package sudokuparser

import (
	"testing"
)

func TestParseSudokuFromFile(t *testing.T) {
	const sample800wi = "7....3..2..4...1.9..52.9....2..15.7...........9.47..8....7.48..3.2...5..9..3....1"
	const sample800wiFile = "../samples/800wi.png"

	if sudokuString := ParseSudokuFromFile(sample800wiFile); sudokuString != sample800wi {
		t.Error(sample800wiFile + " not parsed as " + sample800wi)
	}
}

func TestTrainSudoku(t *testing.T) {
	if sudokuString := TrainSudoku("train_config.csv"); sudokuString != "Testing Train Setup" {
		t.Error("Unexpected response from training Sudoku")
	}
}
