package main

import (
	"fmt"
	"time"

	"github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/engine"
)

func main() {
	pos := backend.ParseFen("8/8/8/8/4k3/8/PPPPPPPP/4K3 w - -")
	var child backend.Position
	pos.Print()
	for i := 0; i < 90; i++ {
		start := time.Now()
		move := engine.Search(&pos)
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
