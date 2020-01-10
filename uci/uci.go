// Based on Counter's implementation
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
	commands  map[string]func(args ...string)
	messages  chan interface{}
	engine    Engine
	positions []backend.Position
	cancel    context.CancelFunc
	state     func(msg interface{})
}

func NewUciProtocol(e Engine) *UciProtocol {
	e.Update = updateUci
	var uci = &UciProtocol{
		messages:  make(chan interface{}),
		engine:    e,
		positions: []backend.Position{backend.InitialPosition},
	}
	uci.commands = map[string]func(args ...string){
		"uci":        uci.uciCommand,
		"isready":    uci.isReadyCommand,
		"position":   uci.positionCommand,
		"go":         uci.goCommand,
		"ucinewgame": uci.uciNewGameCommand,
		"ponderhit":  uci.ponderhitCommand,
		"stop":       uci.stopCommand,
		"setoption":  uci.setOptionCommand,
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
			uci.stopCommand()
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
			cmd(fields[1:]...)
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

func (uci *UciProtocol) uciCommand(...string) {
	var name, version, author = uci.engine.GetInfo()
	fmt.Printf("id name %s %s\n", name, version)
	fmt.Printf("id author %s\n", author)
	for _, option := range uci.engine.GetOptions() {
		fmt.Printf("option name %v type %v default %v min %v max %v\n",
			option.Name, "spin", option.Val, option.Min, option.Max)
	}
	fmt.Println("uciok")
}

func (uci *UciProtocol) isReadyCommand(...string) {
	uci.engine.NewGame()
	fmt.Println("readyok")
}

func (uci *UciProtocol) positionCommand(args ...string) {
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

func (uci *UciProtocol) goCommand(fields ...string) {
	var limits = parseLimits(fields)
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

func (uci *UciProtocol) uciNewGameCommand(...string) {
	uci.engine.NewGame()
}

func (uci *UciProtocol) ponderhitCommand(...string) {
	debugUci("Not implemented")
}

func (uci *UciProtocol) stopCommand(...string) {
	if uci.cancel != nil {
		uci.cancel()
	}
}

func updateUci(s SearchInfo) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("info depth %d nodes %d score ", s.Depth, s.Nodes))
	if s.Score.Mate != 0 {
		sb.WriteString(fmt.Sprintf("mate %d ", s.Score.Mate))
	} else {
		sb.WriteString(fmt.Sprintf("cp %d ", s.Score.Centipawn))
	}

	sb.WriteString("pv ")
	for _, move := range s.Moves {
		sb.WriteString(move.String())
		sb.WriteString(" ")
	}
	sb.WriteString("\n")
	fmt.Print(sb.String())
}

func (uci *UciProtocol) setOptionCommand(fields ...string) {
	if len(fields) < 4 {
		debugUci("invalid setoption arguments")
		return
	}

	var valIdx = findIndexString(fields, "value")
	var name = strings.Join(fields[1:valIdx], " ")
	var value = fields[valIdx+1]

	for _, option := range uci.engine.GetOptions() {
		if strings.EqualFold(option.Name, name) {
			v, err := strconv.Atoi(value)
			if err != nil {
				debugUci("invalid setoption arguments")
				return
			}
			if v < option.Min || v > option.Max {
				debugUci("argument out of range")
				return
			}
			option.Val = v
			return
		}
	}
	debugUci("unhandled option")
}
