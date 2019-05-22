package main

import (
	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/uci"
	"os"
)

func main() {
	if len(os.Args) > 0 && os.Args[1] == "tune" {
		evaluation.Tune()
		return
	}
	e := engine.NewEngine()
	uci := uci.NewUciProtocol(e)
	uci.Run()
}
