#!/bin/bash

mkdir -p build

echo "Building cako.io-cli-darwin-amd64"
GOOS=darwin GOARCH=amd64 go build -o build/cako.io-cli-darwin-amd64

echo "Building cako.io-cli-darwin-arm64"
GOOS=darwin GOARCH=arm64 go build -o build/cako.io-cli-darwin-arm64

echo "Building cako.io-cli-linux-amd64"
GOOS=linux GOARCH=amd64 go build -o build/cako.io-cli-linux-amd64

echo "Building cako.io-cli-linux-arm64"
GOOS=linux GOARCH=arm64 go build -o build/cako.io-cli-linux-arm64

echo "Building cako.io-cli-windows-amd64.exe"
GOOS=windows GOARCH=amd64 go build -o build/cako.io-cli-windows-amd64.exe

echo "Building cako.io-cli-windows-arm64.exe"
GOOS=windows GOARCH=arm64 go build -o build/cako.io-cli-windows-arm64.exe