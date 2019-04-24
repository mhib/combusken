combusken:
	mkdir -p "${GOPATH}/src/github.com/mhib"
	ln -s -f `pwd` "${GOPATH}/src/github.com/mhib/"
	go build combusken.go
