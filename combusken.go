package main

import (
	"fmt"
	"github.com/mhib/combusken/backend"
	"time"
)

func main() {
	position := backend.InitialPosition
	start := time.Now()
	fmt.Println(backend.Perft(&position, 6))
	fmt.Println(time.Now().Sub(start))
}
