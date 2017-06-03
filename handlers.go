package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/jamesandersen/gosudoku/sudokuparser"
)

type Page struct {
	Title   string
	Image   template.URL
	Success bool
	Body    []byte
	Error   string
}

func sudokuFormHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	t, _ := template.ParseFiles("sudoku.html")
	t.Execute(w, p)
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	var p *Page
	t, _ := template.ParseFiles("sudoku.html")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovering from panic in solveHandler...", r)
			p = &Page{Title: "Error", Error: fmt.Sprintf("%v", r), Success: false}
			t.Execute(w, p)
		}
	}()

	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("sudokuFile")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)

	bytes := Buf.Bytes()
	parsed := sudokuparser.ParseSudokuFromByteArray(bytes)

	fmt.Println("Parsed sudoku: " + parsed)

	board := NewSudoku(parsed, STANDARD)
	fmt.Print("Attempting to solve Sudoku...\n")
	board.Print()
	finalBoard, success := board.Solve()

	imgDataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
	if success {
		p = &Page{Title: "Solved", Body: []byte(finalBoard.ToString()), Image: template.URL(imgDataURL), Success: true}
	} else {
		p = &Page{Title: "Not Solved", Body: []byte(finalBoard.ToString()), Image: template.URL(imgDataURL), Success: false}
	}

	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()

	//t, _ := template.ParseFiles("sudoku.html")
	t.Execute(w, p)
}
