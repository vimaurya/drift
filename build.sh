#!/bin/bash

mkdir -p dist

SRC="./cmd/migrate/"

echo "Building binaries..."

GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/drift $SRC

GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/drift.exe $SRC

GOOS=darwin GOARCH=amd64 go build -o dist/darwin-amd64/drift $SRC

GOOS=darwin GOARCH=arm64 go build -o dist/darwin-arm64/drift $SRC

echo "Done! Check the /dist folder."
