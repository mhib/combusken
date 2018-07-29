package main

import (
	"fmt"
	"github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/engine"
)

func main() {
	pos := backend.InitialPosition
	var child backend.Position
	pos.Print()
	for i := 0; i < 600; i++ {
		move := engine.Search(&pos)
		move.Inspect()
		pos.MakeMove(move, &child)
		pos = child
		print("\033[H\033[2J")
		pos.Print()
		fmt.Println("")
	}
}
