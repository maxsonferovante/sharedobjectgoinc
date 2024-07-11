#!/bin/sh

echo "Installing go"
apk add go 

echo "Installing dependencies"
go mod download
go mod tidy

echo "Building go binary into shared library"
go build -buildmode=c-shared -o libcsv.so main.go

echo "Copying shared library to /usr/lib"
cp libcsv.so /usr/lib/libcsv.so

echo "Copying shared library to /usr/local/lib"
cp libcsv.so /usr/local/lib/libcsv.so

echo "Running ldconfig"
ldconfig

echo "Completed building shared library"



