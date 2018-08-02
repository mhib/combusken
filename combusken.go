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
	for i := 0; i < 90; i++ {
		start := time.Now()
		end := time.After(30 * time.Second)
		move := engine.Search(&pos, end)
		move.Inspect()
		if move == 0 {
			return
		}
		pos.MakeMove(move, &child)
		pos = child
		fmt.Println(time.Now().Sub(start))
		pos.Print()
		fmt.Println("")
	}
}
