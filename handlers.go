package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"image/gif"
	"image/png"

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

func gifToPng(file multipart.File) ([]byte, error) {
	img, err := gif.Decode(file)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = png.Encode(writer, img)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Encoded %v byte PNG image\n", b.Len())
	return b.Bytes(), nil
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	var p *page
	isGif := false
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
	defer file.Close()
	fmt.Printf("File name %s\n", header.Filename)
	for key, value := range header.Header {
		for _, val := range value {
			if strings.ToLower(key) == "content-type" && strings.ToLower(val) == "image/gif" {
				isGif = true
			}
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var bytes []byte
	if isGif {
		fmt.Print("Gif image detected, converting to PNG\n")
		bytes, err = gifToPng(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Copy the file data to my buffer
		io.Copy(&Buf, file)
		bytes = Buf.Bytes()
	}

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
