package main

import (
	"fmt"
	"github.com/mhib/combusken/backend"
	"time"
)

func main() {
	backend.InitBB()
	position := backend.InitialPosition
	start := time.Now()
	fmt.Println(backend.Perft(&position, 3))
	fmt.Println(time.Now().Sub(start))
}
