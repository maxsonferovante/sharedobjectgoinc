#!/bin/sh

apk add go 
go mod download
go mod tidy
go build -buildmode=c-shared -o libcsv.so main.go
