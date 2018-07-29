package main

import (
	"fmt"
	"time"

	"github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/engine"
)

func main() {
	pos := backend.InitialPosition
	var child backend.Position
	pos.Print()
	for i := 0; i < 600; i++ {
		start := time.Now()
		move := engine.Search(&pos)
		move.Inspect()
		pos.MakeMove(move, &child)
		pos = child
		fmt.Println(time.Now().Sub(start))
		pos.Print()
		fmt.Println("")
	}
}
