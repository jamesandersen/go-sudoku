/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package sudokusolver

/*
#cgo CPPFLAGS: -I/Users/jandersen/anaconda3/envs/opencvenv/include -I/Users/jandersen/anaconda3/envs/opencvenv/include/opencv2
#cgo CXXFLAGS: --std=c++1z -stdlib=libc++
#cgo darwin LDFLAGS: -L/usr/local/Cellar/opencv3/3.2.0/lib -lopencv_core -lopencv_highgui -lopencv_imgcodecs -lopencv_imgproc -lopencv_ml -lopencv_objdetect -lopencv_photo
#include <stdlib.h>
#include "sudoku_parser.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// SudokuMode represents a variant of sudoku
type SudokuMode int

const (
	// STANDARD Each block, row and column should contain digits 1-9
	STANDARD SudokuMode = iota

	// DIAGONAL Diagonals should also contains digits 1-9
	DIAGONAL SudokuMode = iota
)

const (
	// Row labels
	ROWS = "ABCDEFGHI"

	// Column labels
	COLS = "123456789"
)

// A SudokuBoard represents a sudoku board
type SudokuBoard struct {
	values [9][9]string

	row_units  [][]string
	col_units  [][]string
	sqr_units  [][]string
	diag_units [][]string
	all_units  [][]string
}

// Init - initializes a board type based on the mode
func (b SudokuBoard) Init(mode SudokuMode) SudokuBoard {
	for _, row := range ROWS {
		b.row_units = append(b.row_units, cross(string(row), COLS))
	}

	for _, col := range COLS {
		b.col_units = append(b.col_units, cross(ROWS, string(col)))
	}

	for _, rBlk := range [3]string{"ABC", "DEF", "GHI"} {
		for _, cBlk := range [3]string{"123", "456", "789"} {
			b.sqr_units = append(b.sqr_units, cross(rBlk, cBlk))
		}
	}

	if mode == DIAGONAL {
		unit := []string{}
		for i := 0; i < len(COLS); i++ {
			unit = append(unit, string(ROWS[i])+string(COLS[i]))
		}
		b.diag_units = append(b.diag_units, unit)

		for i := 0; i < len(COLS); i++ {
			unit = append(unit, string(ROWS[i])+string(COLS[len(COLS)-i-1]))
		}
		b.diag_units = append(b.diag_units, unit)
	}

	// add all other units into single units list
	for _, u := range [][][]string{b.row_units, b.col_units, b.sqr_units, b.diag_units} {
		b.all_units = append(b.all_units, u...)
	}

	return b
}

// Cross product of elements in a and elements in b.
func cross(a string, b string) []string {
	values := []string{}
	for _, outer := range a {
		for _, inner := range b {
			values = append(values, string(outer)+string(inner))
		}
	}

	return values
}

// Parse a Sudoku puzzle using a file path to a Sudoku image
func ParseSudokuFromFile(filename string) string {
	if !path.IsAbs(filename) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Getwd failed: %s", err))
		}
		filename = path.Clean(path.Join(cwd, filename))
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return ParseSudokuFromByteArray(data)
}

// Parse a Sudoku puzzle from an image byte array
func ParseSudokuFromByteArray(data []byte) string {
	p := C.CBytes(data)
	defer C.free(p)

	parsed := C.ParseSudoku((*C.char)(p), C.int(len(data)), true)

	return C.GoString(parsed)
}
