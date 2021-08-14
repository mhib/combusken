// Based on Counter's implementation
package uci

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/engine"
)

type UciProtocol struct {
	commands     map[string]func(args ...string)
	messages     chan interface{}
	engine       Engine
	positions    []backend.Position
	cancel       context.CancelFunc
	ponderCancel context.CancelFunc
	state        func(msg interface{})
	waitChan     chan interface{}
}

type searchResult struct {
	bestMove   backend.Move
	ponderMove backend.Move
}

func NewUciProtocol(e Engine) *UciProtocol {
	e.SetUpdate(updateUci)
	uci := &UciProtocol{
		messages:  make(chan interface{}),
		waitChan:  make(chan interface{}),
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
	close(uci.waitChan)
	return uci
}

func (uci *UciProtocol) Run() {
	name, version, _ := uci.engine.GetInfo()
	fmt.Printf("%v %v\n", name, version)
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
			debugUci("Command not found.")
		}
	case searchResult:
		debugUci("Unexpected best move.")
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
			debugUci("Unexpected command " + commandName + ".")
		}
	case searchResult:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("bestmove %s", msg.bestMove.String()))
		if msg.ponderMove != backend.NullMove {
			sb.WriteString(fmt.Sprintf(" ponder %s", msg.ponderMove.String()))
		}
		sb.WriteString("\n")
		fmt.Print(sb.String())
		uci.state = uci.idle
	}
}

func debugUci(s string) {
	fmt.Println("info string " + s)
}

func (uci *UciProtocol) uciCommand(...string) {
	uci.engine.NewGame()
	name, version, author := uci.engine.GetInfo()
	fmt.Printf("id name %s %s\n", name, version)
	fmt.Printf("id author %s\n", author)
	for _, option := range uci.engine.GetOptions() {
		fmt.Println(option.ToUci())
	}
	fmt.Println("uciok")
}

func (uci *UciProtocol) isReadyCommand(...string) {
	<-uci.waitChan
	fmt.Println("readyok")
}

func (uci *UciProtocol) positionCommand(args ...string) {
	uci.waitChan = make(chan interface{})
	defer close(uci.waitChan)
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
		debugUci("Wrong position command")
		return
	}
	p := backend.ParseFen(fen)
	positions := []backend.Position{p}
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			newPos, ok := positions[len(positions)-1].MakeMoveLAN(smove)
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
	uci.waitChan = make(chan interface{})
	defer close(uci.waitChan)
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

func updateUci(s SearchInfo) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("info depth %d nodes %d score ", s.Depth, s.Nodes))
	if s.Score.Mate != 0 {
		sb.WriteString(fmt.Sprintf("mate %d ", s.Score.Mate))
	} else {
		sb.WriteString(fmt.Sprintf("cp %d ", s.Score.Centipawn))
	}
	sb.WriteString(fmt.Sprintf("nps %d ", s.Nps))
	sb.WriteString(fmt.Sprintf("time %d ", s.Duration))
	sb.WriteString(fmt.Sprintf("tbhits %d ", s.Tbhits))

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

	valIdx := findIndexString(fields, "value")
	name := strings.Join(fields[1:valIdx], " ")
	value := fields[valIdx+1]

	for _, option := range uci.engine.GetOptions() {
		if strings.EqualFold(option.GetName(), name) {
			err := option.SetValue(value)
			if err != nil {
				debugUci(err.Error())
			}
			return
		}
	}
	debugUci("unhandled option")
}
