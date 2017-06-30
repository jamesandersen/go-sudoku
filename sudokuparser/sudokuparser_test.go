package sudokuparser

import (
	"testing"
)

func TestParseSudokuFromFile(t *testing.T) {
	const sample800wi = "7....3..2..4...1.9..52.9....2..15.7...........9.47..8....7.48..3.2...5..9..3....1"
	const sample800wiFile = "../samples/800wi.png"

	if sudokuString, _ := ParseSudokuFromFile(sample800wiFile); sudokuString != sample800wi {
		t.Error(sample800wiFile + " not parsed as " + sample800wi + ": \n" + sudokuString)
	}
}

func TestParseNewsprintSudokuFromFile(t *testing.T) {
	const sampleNewsprint = "2....61..1...92.8...7.....4.298......7..5..2......735.4.....9...8.41...7..36....5"
	const sampleNewsprintFile = "../samples/NewsprintSudoku.jpg"

	if sudokuString, _ := ParseSudokuFromFile(sampleNewsprintFile); sudokuString != sampleNewsprint {
		t.Error(sampleNewsprintFile + " not parsed as " + sampleNewsprint + ": \n" + sudokuString)
	}
}

func TestTrainSudoku(t *testing.T) {
	if sudokuString := TrainSudoku("train_config.csv"); sudokuString != "97.73" {
		t.Error("Unexpected response from training Sudoku: " + sudokuString)
	}
}
