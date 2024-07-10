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

	filteredMatrix := builders.BuilderFilteredMatrix(matrix, goRowFilterDefinitions)

	selectedMatrix := builders.BuilderSelectedMatrix(filteredMatrix, goSelectedColumns)
	builders.BuilderApreciation(selectedMatrix)

}

//export processCsvFile
func processCsvFile(csvFilePath, selectedColumns, rowFilterDefinitions *C.char) {

	goCsv := C.GoString(csvFilePath)

	records := loadFile(goCsv)

	csv := C.CString(records)

	processCsv(csv, selectedColumns, rowFilterDefinitions)
}

func main() {
	// This is necessary to build the shared library
	// go build -o libcsv.so -buildmode=c-shared main.go
}

/* csv := "header1,header2,header3\n1,2,3\n4,5,6\n7,8,9"
processCsv(csv, "header1,header3", "header1>1\nheader3<8")
// output
// header1,header3
// 4,6

// processCsv(csv, "header1,header3", "header1=4\nheader2>3")
// */
// filePath := "data.csv"
// processCsvFile(filePath, "col1,col3,col4", "col1>l1c1\ncol3>l2c3")
