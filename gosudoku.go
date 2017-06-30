/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/jamesandersen/gosudoku/sudokuparser"
	"github.com/nytimes/gziphandler"
)

func main() {
	var filename string
	var mode string
	flag.StringVar(&mode, "mode", "serve", "whether to serve web app or parse additional args in CLI mode")
	flag.StringVar(&filename, "filename", "", "Sudoku puzzle image")

	flag.Parse()

	if mode == "serve" {
		fs := http.FileServer(http.Dir("web/static"))
		http.Handle("/static/", gziphandler.GzipHandler(http.StripPrefix("/static/", fs)))
		http.HandleFunc("/solve", solveHandler)
		http.HandleFunc("/", sudokuFormHandler)
		http.ListenAndServe(":8080", nil)
	} else if mode == "cli" {
		parsed, points := sudokuparser.ParseSudokuFromFile(filename)

		fmt.Println("Parsed sudoku: " + parsed)
		if len(points) == 4 {
			fmt.Print(fmt.Sprintf("TopLeft (%c, %c)\n", points[0].X, points[0].Y))
			fmt.Print(fmt.Sprintf("TopRight (%c, %c)\n", points[1].X, points[1].Y))
			fmt.Print(fmt.Sprintf("BottomRight (%c, %c)\n", points[2].X, points[2].Y))
			fmt.Print(fmt.Sprintf("BottomLeft (%c, %c)\n", points[3].X, points[3].Y))
		}

		board := NewSudoku(parsed, STANDARD)
		fmt.Print("Attempting to solve Sudoku...\n")
		board.Print()
		finalBoard, success := board.Solve()

		if success {
			finalBoard.Print()
		} else {
			fmt.Print("Board not solved...")
		}
	}
}
