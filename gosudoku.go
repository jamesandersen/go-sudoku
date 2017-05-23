/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/jamesandersen/gosudoku/sudokusolver"
)

func main() {
	var filename string

	flag.StringVar(&filename, "filename", "", "Sudoku puzzle image")

	flag.Parse()

	parsed := sudokusolver.ParseSudokuFromFile(filename)

	fmt.Println("Parsed sudoku: " + parsed)
}
