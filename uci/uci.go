// Based on https://github.com/ChizhovVadim/CounterGo/blob/master/shell/uciprotocol.go
package uci

import . "github.com/mhib/combusken/engine"
import "github.com/mhib/combusken/backend"
import "fmt"
import "context"

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type UciProtocol struct {
	commands  map[string]func()
	messages  chan interface{}
	engine    Engine
	positions []backend.Position
	cancel    context.CancelFunc
	fields    []string
	state     func(msg interface{})
}

func NewUciProtocol(e Engine) *UciProtocol {
	var initPosition = backend.InitialPosition
	e.Update = updateUci
	var uci = &UciProtocol{
		messages:  make(chan interface{}),
		engine:    e,
		positions: []backend.Position{initPosition},
	}
	uci.commands = map[string]func(){
		"uci":        uci.uciCommand,
		"isready":    uci.isReadyCommand,
		"position":   uci.positionCommand,
		"go":         uci.goCommand,
		"ucinewgame": uci.uciNewGameCommand,
		"ponderhit":  uci.ponderhitCommand,
		"stop":       uci.stopCommand,
	}
	return uci
}

func (uci *UciProtocol) Run() {
	var name, version, _ = uci.engine.GetInfo()
	fmt.Printf("%v %v\n", name, version)
	go func() {
		uci.state = uci.idle
		for msg := range uci.messages {
			uci.state(msg)
		}
	}()
	var scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var commandLine = scanner.Text()
		if commandLine == "quit" {
			break
		}
		uci.messages <- commandLine
	}
}

func (uci *UciProtocol) idle(msg interface{}) {
	switch msg := msg.(type) {
	case string:
		var fields = strings.Fields(msg)
		if len(fields) == 0 {
			return
		}
		var commandName = fields[0]
		var cmd, ok = uci.commands[commandName]
		if ok {
			uci.fields = fields[1:]
			cmd()
		} else {
			debugUci("Command not found.")
		}
	case backend.Move:
		debugUci("Unexpected best move.")
	}
}

func (uci *UciProtocol) thinking(msg interface{}) {
	switch msg := msg.(type) {
	case string:
		var fields = strings.Fields(msg)
		if len(fields) == 0 {
			return
		}
		var commandName = fields[0]
		if commandName == "stop" {
			uci.stopCommand()
		} else {
			debugUci("Unexpected command " + commandName + ".")
		}
	case backend.Move:
		fmt.Printf("bestmove %s\n", msg.String())
		uci.state = uci.idle
	}
}

func debugUci(s string) {
	fmt.Println("info string " + s)
}

func (uci *UciProtocol) uciCommand() {
	var name, version, author = uci.engine.GetInfo()
	fmt.Printf("id name %s %s\n", name, version)
	fmt.Printf("id author %s\n", author)
	fmt.Println("uciok")
}

func (uci *UciProtocol) isReadyCommand() {
	fmt.Println("readyok")
}

func (uci *UciProtocol) positionCommand() {
	var args = uci.fields
	var token = args[0]
	var fen string
	var movesIndex = findIndexString(args, "moves")
	if token == "startpos" {
		fen = backend.InitialPositionFen
	} else if token == "fen" {
		if movesIndex == -1 {
			fen = strings.Join(args[1:], " ")
		} else {
			fen = strings.Join(args[1:movesIndex], " ")
		}
	} else {
		debugUci("Wrong position command")
		return
	}
	var p = backend.ParseFen(fen)
	var positions = []backend.Position{p}
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			var newPos, ok = positions[len(positions)-1].MakeMoveLAN(smove)
			if !ok {
				debugUci("Wrong move")
				return
			}
			positions = append(positions, newPos)
		}
	}
	uci.positions = positions
}

func findIndexString(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func (uci *UciProtocol) goCommand() {
	var limits = parseLimits(uci.fields)
	var ctx, cancel = context.WithCancel(context.Background())
	var searchParams = SearchParams{
		Positions: uci.positions,
		Limits:    limits,
	}
	uci.cancel = cancel
	uci.state = uci.thinking
	go func() {
		var searchResult = uci.engine.Search(ctx, searchParams)
		uci.messages <- searchResult
	}()
}

func parseLimits(args []string) (result LimitsType) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "ponder":
			result.Ponder = true
		case "wtime":
			result.WhiteTime, _ = strconv.Atoi(args[i+1])
			i++
		case "btime":
			result.BlackTime, _ = strconv.Atoi(args[i+1])
			i++
		case "winc":
			result.WhiteIncrement, _ = strconv.Atoi(args[i+1])
			i++
		case "binc":
			result.BlackIncrement, _ = strconv.Atoi(args[i+1])
			i++
		case "movestogo":
			result.MovesToGo, _ = strconv.Atoi(args[i+1])
			i++
		case "depth":
			result.Depth, _ = strconv.Atoi(args[i+1])
			i++
		case "nodes":
			result.Nodes, _ = strconv.Atoi(args[i+1])
			i++
		case "mate":
			result.Mate, _ = strconv.Atoi(args[i+1])
			i++
		case "movetime":
			result.MoveTime, _ = strconv.Atoi(args[i+1])
			i++
		case "infinite":
			result.Infinite = true
		}
	}
	return
}

func (uci *UciProtocol) uciNewGameCommand() {
	uci.engine.TransTable.Clear()
}

func (uci *UciProtocol) ponderhitCommand() {
}

func (uci *UciProtocol) stopCommand() {
	if uci.cancel != nil {
		uci.cancel()
	}
}

type Uci struct {
	engine Engine
	cancel context.CancelFunc
}

func updateUci(s SearchInfo) {
	fmt.Printf("info depth %d score cp %d\n", s.Depth, s.Score)
}

func (u *Uci) uciCommand() {
	var name, version, author = u.engine.GetInfo()
	fmt.Printf("id name %s %s\n", name, version)
	fmt.Printf("id author %s\n", author)
	fmt.Println("uciok")
}
