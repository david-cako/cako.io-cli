#!/bin/bash

mkdir -p build

echo "Building cako-darwin-amd64"
GOOS=darwin GOARCH=amd64 go build -o build/cako-darwin-amd64

echo "Building cako-darwin-arm64"
GOOS=darwin GOARCH=arm64 go build -o build/cako-darwin-arm64

echo "Building cako-linux-amd64"
GOOS=linux GOARCH=amd64 go build -o build/cako-linux-amd64

echo "Building cako-linux-arm64"
GOOS=linux GOARCH=arm64 go build -o build/cako-linux-arm64

echo "Building cako-windows-amd64.exe"
GOOS=windows GOARCH=amd64 go build -o build/cako-windows-amd64.exe

echo "Building cako-windows-arm64.exe"
GOOS=windows GOARCH=arm64 go build -o build/cako-windows-arm64.exe