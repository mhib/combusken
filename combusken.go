package main

import (
	"fmt"
	"github.com/mhib/combusken/backend"
)

func main() {
	backend.InitBB()
	position := backend.InitialPosition
	fmt.Println(backend.Perft(&position, 6))
}
