#!/bin/bash

mkdir -p dist

SRC="./cmd/migrate/"

echo "Building binaries..."

GOOS=linux GOARCH=amd64 go build -o dist/drift-linux-amd64 $SRC

GOOS=windows GOARCH=amd64 go build -o dist/drift-windows-amd64.exe $SRC

GOOS=darwin GOARCH=amd64 go build -o dist/drift-darwin-amd64 $SRC

GOOS=darwin GOARCH=arm64 go build -o dist/drift-darwin-arm64 $SRC

echo "Done! Check the /dist folder."
