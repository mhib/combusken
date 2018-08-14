package main

import (
	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/uci"
)

func main() {
	e := engine.NewEngine()
	uci := uci.NewUciProtocol(e)
	uci.Run()
}
