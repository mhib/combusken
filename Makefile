# DO NOT USE THIS
# This file is provided to make it possible to compile project outside of GOPATH in OpenBench instance

EXE = combusken
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

combusken:
	mkdir -p "${GOPATH}/src/github.com/mhib"
	ln -s -f "$(ROOT_DIR)" "${GOPATH}/src/github.com/mhib/combusken"
	go build -o $(EXE) combusken.go
