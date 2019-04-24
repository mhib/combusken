# DO NOT USE THIS
# This file is provided to make it possible to compile project outside of GOPATH in OpenBench instance

EXE = combusken

combusken:
	mkdir -p "${GOPATH}/src/github.com/mhib"
	ln -s -f `pwd` "${GOPATH}/src/github.com/mhib/"
	go build -o $(EXE) combusken.go
