package main

import (
	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/uci"
)

func main() {
	e := engine.Engine{}
	uci := uci.NewUciProtocol(e)
	uci.Run()
}
