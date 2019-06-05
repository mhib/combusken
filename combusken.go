package main

import (
	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/uci"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "tune" {
			evaluation.Tune()
		} else if os.Args[1] == "bench" {
			engine.Benchmark()
		}
		return
	}
	e := engine.NewEngine()
	uci := uci.NewUciProtocol(e)
	uci.Run()
}
