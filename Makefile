EXE = combusken

combusken:
	go build -gcflags -B -o $(EXE) combusken.go
