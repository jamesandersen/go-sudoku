package sudokuparser

import (
	"testing"
)

func TestParseSudokuFromFile(t *testing.T) {
	const sample800wi = "7,,,,,3,,,2,,,4,,,,1,,9,,,5,2,,9,,,,,2,,,1,5,,7,,,,,,,,,,,,9,,4,7,,,8,,,,,7,,4,8,,,3,,2,,,,5,,,9,,,3,,,,,1"
	const sample800wiFile = "../samples/800wi.png"

	if sudokuString := ParseSudokuFromFile(sample800wiFile); sudokuString != sample800wi {
		t.Error(sample800wiFile + " not parsed as " + sample800wi)
	}
}
