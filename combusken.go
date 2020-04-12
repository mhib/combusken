package main

import (
	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/uci"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "tune":
			evaluation.Tune()
		case "trace-tune":
			evaluation.TraceTune()
		case "bench":
			engine.Benchmark()
		}
	}
	uci := uci.NewUciProtocol(engine.NewEngine())
	uci.Run()
}
