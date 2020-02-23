# DO NOT USE THIS
# This file is provided to make it possible to compile project outside of GOPATH in OpenBench instance

EXE = combusken
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

combusken:
	mkdir -p "${GOPATH}/src/github.com/mhib"
	ln -s -f "$(ROOT_DIR)" "${GOPATH}/src/github.com/mhib/combusken"
	go build -o $(EXE) combusken.go
	rm "${GOPATH}/src/github.com/mhib/combusken"


build:
	GOOS=linux   GOARCH=amd64 go build -o combusken-linux-64       combusken.go
	GOOS=windows GOARCH=amd64 go build -o combusken-windows-64.exe combusken.go
	GOOS=darwin  GOARCH=amd64 go build -o combusken-osx-64         combusken.go
	GOOS=linux  GOARCH=arm64 go build -o combusken-arm-64          combusken.go
