package main

import "C"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"libcsv/pkg/builders"
	"libcsv/pkg/valitations"
	"os"
	"sync"
)

var mtx sync.Mutex

func loadFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		message := fmt.Sprintf("Error in read: %s", err)
		err := errors.New(message)
		panic(err)
	}
	defer file.Close()

	records, err := ioutil.ReadAll(file)
	if err != nil {
		message := fmt.Sprintf("Error in read: %s", err)
		err := errors.New(message)
		panic(err)

	}
	return string(records)

}

//export processCsv
func processCsv(csv, selectedColumns, rowFilterDefinitions *C.char) {
	mtx.Lock()
	defer mtx.Unlock()
	goCsv := C.GoString(csv)
	goSelectedColumns := C.GoString(selectedColumns)
	goRowFilterDefinitions := C.GoString(rowFilterDefinitions)

	valitations.ValidateCSV(goCsv)

	matrix := builders.BuilderMatrix(goCsv)
	err, _ := valitations.ValidationSelectedColumns(matrix, goSelectedColumns)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	err, _ = valitations.ValidationFilterDefinitions(matrix, goRowFilterDefinitions)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	filteredMatrix, err := builders.BuilderFilteredMatrix(matrix, goRowFilterDefinitions)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	selectedMatrix := builders.BuilderSelectedMatrix(filteredMatrix, goSelectedColumns)
	builders.BuilderApreciation(selectedMatrix)

}

//export processCsvFile
func processCsvFile(csvFilePath, selectedColumns, rowFilterDefinitions *C.char) {

	mtx.Lock()
	defer mtx.Unlock()

	goCsv := C.GoString(csvFilePath)
	goSelectedColumns := C.GoString(selectedColumns)
	goRowFilterDefinitions := C.GoString(rowFilterDefinitions)

	records := loadFile(goCsv)

	valitations.ValidateCSV(records)

	matrix := builders.BuilderMatrix(records)
	err, _ := valitations.ValidationSelectedColumns(matrix, goSelectedColumns)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	err, _ = valitations.ValidationFilterDefinitions(matrix, goRowFilterDefinitions)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	filteredMatrix, err := builders.BuilderFilteredMatrix(matrix, goRowFilterDefinitions)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}

	selectedMatrix := builders.BuilderSelectedMatrix(filteredMatrix, goSelectedColumns)
	builders.BuilderApreciation(selectedMatrix)
}

func main() {
	// This is necessary to build the shared library
	// go build -o libcsv.so -buildmode=c-shared main.go
}
