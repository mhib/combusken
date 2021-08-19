// Based on Counter's implementation
package uci

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/engine"
)

type UciProtocol struct {
	writeBuffer  bytes.Buffer
	commands     map[string]func(args ...string)
	messages     chan interface{}
	engine       *Engine
	positions    []backend.Position
	cancel       context.CancelFunc
	ponderCancel context.CancelFunc
	state        func(msg interface{})
}

type searchResult struct {
	bestMove   backend.Move
	ponderMove backend.Move
}

func NewUciProtocol(e *Engine) *UciProtocol {
	var updateBufferArray [1000]byte
	uci := &UciProtocol{
		messages:    make(chan interface{}),
		engine:      e,
		positions:   []backend.Position{backend.InitialPosition},
		writeBuffer: *bytes.NewBuffer(updateBufferArray[:]),
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
	e.SetUpdate(func(s *SearchInfo) {
		uci.messages <- s
	})
	return uci
}

func (uci *UciProtocol) Run() {
	name, version, _ := uci.engine.GetInfo()
	uci.printLn(name, " ", version)
	go func() {
		uci.state = uci.idle
		for msg := range uci.messages {
			uci.state(msg)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		commandLine := scanner.Text()
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
		fields := strings.Fields(msg)
		if len(fields) == 0 {
			return
		}
		commandName := fields[0]
		cmd, ok := uci.commands[commandName]
		if ok {
			cmd(fields[1:]...)
		} else {
			uci.debug("Command not found.")
		}
	case searchResult:
		uci.debug("Unexpected best move.")
	case *SearchInfo:
		uci.debug("Unexpected searchInfo.")
	}
}

func (uci *UciProtocol) thinking(msg interface{}) {
	switch msg := msg.(type) {
	case string:
		fields := strings.Fields(msg)
		if len(fields) == 0 {
			return
		}
		commandName := fields[0]
		switch commandName {
		case "stop":
			uci.stopCommand()
		case "ponderhit":
			uci.ponderhitCommand()
		default:
			uci.debug("Unexpected command " + commandName + ".")
		}
	case searchResult:
		uci.writeBuffer.WriteString("bestmove ")
		uci.writeBuffer.WriteString(msg.bestMove.String())
		if msg.ponderMove != backend.NullMove {
			uci.writeBuffer.WriteString(" ")
			uci.writeBuffer.WriteString(msg.ponderMove.String())
		}
		uci.writeBuffer.WriteString("\n")
		uci.writeBuffer.WriteTo(os.Stdout)
		uci.state = uci.idle
	case *SearchInfo:
		uci.update(msg)
	}
}

func (uci *UciProtocol) debug(s string) {
	uci.writeBuffer.WriteString("info string ")
	uci.writeBuffer.WriteString(s)
	uci.writeBuffer.WriteString("\n")
	uci.writeBuffer.WriteTo(os.Stdout)
}

func (uci *UciProtocol) printLn(strings ...string) {
	for idx := range strings {
		uci.writeBuffer.WriteString(strings[idx])
	}
	uci.writeBuffer.WriteString("\n")
	uci.writeBuffer.WriteTo(os.Stdout)
}

func (uci *UciProtocol) uciCommand(...string) {
	uci.engine.NewGame()
	name, version, author := uci.engine.GetInfo()
	uci.printLn("id name ", name, " ", version)
	uci.printLn("id author ", author)
	for _, option := range uci.engine.GetOptions() {
		uci.printLn(option.ToUci())
	}
	uci.printLn("uciok")
}

func (uci *UciProtocol) isReadyCommand(...string) {
	uci.printLn("readyok")
}

func (uci *UciProtocol) positionCommand(args ...string) {
	var fen string
	token := args[0]
	movesIndex := findIndexString(args, "moves")
	if token == "startpos" {
		fen = backend.InitialPositionFen
	} else if token == "fen" {
		if movesIndex == -1 {
			fen = strings.Join(args[1:], " ")
		} else {
			fen = strings.Join(args[1:movesIndex], " ")
		}
	} else {
		uci.debug("Wrong position command")
		return
	}
	p := backend.ParseFen(fen)
	positions := []backend.Position{p}
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			newPos, ok := positions[len(positions)-1].MakeMoveLAN(smove)
			if !ok {
				uci.debug("Wrong move")
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
	limits := parseLimits(fields)
	ctx, cancel := context.WithCancel(context.Background())
	ponderCtx, ponderCancel := context.WithCancel(ctx)
	searchParams := SearchParams{
		Positions: uci.positions,
		Limits:    limits,
	}
	uci.cancel = cancel
	uci.ponderCancel = ponderCancel
	if !limits.Ponder {
		ponderCancel()
	}
	uci.state = uci.thinking
	go func() {
		bestMove, ponderMove := uci.engine.Search(ctx, ponderCtx, searchParams)
		cancel()
		uci.messages <- searchResult{bestMove, ponderMove}
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
	if uci.ponderCancel != nil {
		uci.ponderCancel()
	}
}

func (uci *UciProtocol) stopCommand(...string) {
	if uci.cancel != nil {
		uci.cancel()
	}
}

func (uci *UciProtocol) update(s *SearchInfo) {
	uci.writeBuffer.WriteString(fmt.Sprintf("info depth %d seldepth %d nodes %d score ", s.Depth, s.SelDepth, s.Nodes))
	if s.Score.Mate != 0 {
		uci.writeBuffer.WriteString(fmt.Sprintf("mate %d ", s.Score.Mate))
	} else {
		uci.writeBuffer.WriteString(fmt.Sprintf("cp %d ", s.Score.Centipawn))
	}
	uci.writeBuffer.WriteString(fmt.Sprintf("nps %d ", s.Nps))
	uci.writeBuffer.WriteString(fmt.Sprintf("time %d ", s.Duration))
	uci.writeBuffer.WriteString(fmt.Sprintf("tbhits %d ", s.Tbhits))

	uci.writeBuffer.WriteString("pv ")
	for _, move := range s.Moves {
		uci.writeBuffer.WriteString(move.String())
		uci.writeBuffer.WriteString(" ")
	}
	uci.writeBuffer.WriteString("\n")
	uci.writeBuffer.WriteTo(os.Stdout)
}

func (uci *UciProtocol) setOptionCommand(fields ...string) {
	if len(fields) < 4 {
		uci.debug("invalid setoption arguments")
		return
	}

	valIdx := findIndexString(fields, "value")
	name := strings.Join(fields[1:valIdx], " ")
	value := fields[valIdx+1]

	for _, option := range uci.engine.GetOptions() {
		if strings.EqualFold(option.GetName(), name) {
			err := option.SetValue(value)
			if err != nil {
				uci.debug(err.Error())
			}
			return
		}
	}
	uci.debug("unhandled option")
}
