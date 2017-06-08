package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/jamesandersen/gosudoku/sudokuparser"
)

type page struct {
	Title   string
	Image   template.URL
	Success bool
	Body    string
	Values  map[string]CellValue
	Points  []sudokuparser.Point2d
	Error   string
}

func sudokuFormHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p := &page{Title: "TestPage", Body: "--Solution goes here--."}
		t, _ := template.ParseFiles("web/sudoku.html")
		t.Execute(w, p)
	case "POST":
		solveHandler(w, r)
	default:
		// Give an error message.
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	var p *page
	t, _ := template.ParseFiles("web/sudoku.html")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovering from panic in solveHandler...", err)
			p = &page{Title: "Error", Error: fmt.Sprintf("%v", err), Success: false}
			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				js, jerr := json.Marshal(p)
				if jerr != nil {
					http.Error(w, jerr.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				http.Error(w, string(js), http.StatusInternalServerError)
			} else {
				t.Execute(w, p)
			}
		}
	}()

	var Buf bytes.Buffer
	file, header, err := r.FormFile("sudokuFile")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Printf("File name %s\n", header.Filename)
	// Copy the file data to my buffer
	io.Copy(&Buf, file)

	bytes := Buf.Bytes()
	parsed, points := sudokuparser.ParseSudokuFromByteArray(bytes)

	fmt.Println("Parsed sudoku: " + parsed)

	board := NewSudoku(parsed, STANDARD)
	fmt.Print("Attempting to solve Sudoku...\n")
	board.Print()
	finalBoard, success := board.Solve()

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		if success {
			p = &page{Title: "Solved", Body: finalBoard.ToString(), Success: true, Values: finalBoard.values, Points: points, Error: ""}
		} else {
			p = &page{Title: "Not Solved", Body: finalBoard.ToString(), Success: false, Error: ""}
		}
		js, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		imgDataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
		if success {
			p = &page{Title: "Solved", Body: finalBoard.ToString(), Image: template.URL(imgDataURL), Success: true}
		} else {
			p = &page{Title: "Not Solved", Body: finalBoard.ToString(), Image: template.URL(imgDataURL), Success: false}
		}
		t.Execute(w, p)
	}
}
