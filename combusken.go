package main

import (
	"os"

	"github.com/mhib/combusken/engine"
	"github.com/mhib/combusken/tuning"
	"github.com/mhib/combusken/uci"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "tune":
			tuning.Tune()
		case "trace-tune":
			tuning.TraceTune()
		case "bench":
			engine.Benchmark()
		}
		return
	}
	uci := uci.NewUciProtocol(engine.NewEngine())
	uci.Run()
}
