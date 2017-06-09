/*
Package go-sudoku implements a simple library for solving sudoku puzzles.
*/
package sudokuparser

// FOO CPPFLAGS: -I/Users/jandersen/anaconda3/envs/opencvenv/include -I/Users/jandersen/anaconda3/envs/opencvenv/include/opencv2
// Foo CPPFLAGS: -I/usr/local/Cellar/opencv3/3.2.0/include -I/usr/local/Cellar/opencv3/3.2.0/include/opencv2

/*
#cgo darwin CPPFLAGS: -I/usr/local/Cellar/opencv3/3.2.0/include -I/usr/local/Cellar/opencv3/3.2.0/include/opencv2
#cgo darwin CXXFLAGS: --std=c++1z -stdlib=libc++
#cgo darwin LDFLAGS: -L/usr/local/Cellar/opencv3/3.2.0/lib -lopencv_core -lopencv_highgui -lopencv_imgcodecs -lopencv_imgproc -lopencv_ml -lopencv_objdetect -lopencv_photo
#cgo linux CPPFLAGS: -I/usr/include -I/usr/include/opencv2 -I/usr/local/include -I/usr/local/include/opencv2
#cgo linux CXXFLAGS: --std=c++1z
#cgo linux LDFLAGS: -L/usr/lib -lopencv_core -lopencv_highgui -lopencv_imgcodecs -lopencv_imgproc -lopencv_ml -lopencv_objdetect -lopencv_photo
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
	"strings"
	"unsafe"
)

// Point2d represents a point in 2d space
type Point2d struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var svmModelPath string

// ParseSudokuFromFile parses a Sudoku puzzle using a file path to a Sudoku image
func ParseSudokuFromFile(filename string) (string, []Point2d) {
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

// ParseSudokuFromByteArray parses a Sudoku puzzle from an image byte array
func ParseSudokuFromByteArray(data []byte) (string, []Point2d) {
	if svmModelPath == "" {
		svmModelPath = setupSVMModel()
		svmEnvVar := C.GoString(C.SVM_MODEL_VAR)
		err := os.Setenv(svmEnvVar, svmModelPath)
		if err != nil {
			panic(err)
		}

		fmt.Print(fmt.Sprintf("Set environment variable %s=%s\n", svmEnvVar, svmModelPath))
	}

	parsed := C.CString(strings.Repeat("0", 81))
	defer C.free(unsafe.Pointer(parsed))

	p := C.CBytes(data)
	defer C.free(unsafe.Pointer(p))

	// float32 is standard type compatible with C
	gridCoords := []float32{-1, -1, -1, -1, -1, -1, -1, -1}

	C.ParseSudoku((*C.char)(p), C.int(len(data)), (*C.float)(unsafe.Pointer(&gridCoords[0])), true, parsed)

	coords := []Point2d{}
	for i := 0; i < 4; i++ {
		x := i * 2
		y := (i * 2) + 1
		if gridCoords[x] > -1 && gridCoords[y] > -1 {
			fmt.Print(fmt.Sprintf("Grid coord (%f, %f)\n", gridCoords[i*2], gridCoords[(i*2)+1]))
			coords = append(coords, Point2d{X: int(gridCoords[x]), Y: int(gridCoords[y])})
		}
	}
	goString := C.GoString(parsed)

	return goString, coords
}

func setupSVMModel() string {
	tmpDir := path.Join(os.TempDir(), "sudokusolver")
	os.MkdirAll(tmpDir, os.ModePerm)

	tmpModelFile := path.Join(tmpDir, "sudokuSVMModel.yml")
	fmt.Println("Checking for SVM model file at " + tmpModelFile)
	finfo, err := os.Stat(tmpModelFile)
	if err != nil {
		// create the file
		data, err := Asset("data/model4.yml")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(tmpModelFile, data, os.ModePerm)
		if err != nil {
			panic(err)
		}

		fmt.Println("SVM model file written to " + tmpModelFile)
	} else if finfo.IsDir() {
		panic(tmpModelFile + " is a directory")
	}

	return tmpModelFile
}

// Parse a Sudoku puzzle from an image byte array
func TrainSudoku(trainConfigFile string) string {

	parsed := C.TrainSudoku(C.CString(trainConfigFile))

	return C.GoString(parsed)
}
