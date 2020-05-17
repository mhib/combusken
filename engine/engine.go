package engine

import (
	"context"
	"errors"
	"runtime"

	"github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/fathom"
	"github.com/mhib/combusken/transposition"

	. "github.com/mhib/combusken/utils"
)

const MAX_HEIGHT = 127
const STACK_SIZE = MAX_HEIGHT + 1
const MAX_MOVES = 256

var errTimeout = errors.New("Search timeout")

type Engine struct {
	Hash              IntOption
	Threads           IntOption
	MoveOverhead      IntOption
	PawnHash          IntOption
	SyzygyPath        StringOption
	SyzygyProbeDepth  IntOption
	done              <-chan struct{}
	RepeatedPositions map[uint64]interface{}
	MovesCount        int
	Update            func(SearchInfo)
	timeManager
	threads []thread
}

type thread struct {
	engine *Engine
	MoveHistory
	nodes int
	stack [STACK_SIZE]StackEntry
}

type UciScore struct {
	Mate      int
	Centipawn int
}

func newUciScore(score int) UciScore {
	if score >= ValueWin {
		return UciScore{Mate: (Mate - score + 1) / 2}
	} else if score <= ValueLoss {
		return UciScore{Mate: (-Mate - score) / 2}
	} else {
		return UciScore{Centipawn: score}
	}
}

type SearchInfo struct {
	Score    UciScore
	Depth    int
	Nodes    int
	Nps      int
	Duration int
	Moves    []backend.Move
}

type StackEntry struct {
	MoveProvider
	PV
	quietsSearched       [MAX_MOVES]backend.Move
	position             backend.Position
	evaluation           int16
	evaluationCalculated bool
}

func (se *StackEntry) InvalidateEvaluation() {
	se.evaluationCalculated = false
}

func (se *StackEntry) Evaluation() int16 {
	if !se.evaluationCalculated {
		se.evaluation = int16(evaluation.Evaluate(&se.position))
		se.evaluationCalculated = true
	}
	return se.evaluation
}

type PV struct {
	size  int
	items [STACK_SIZE]backend.Move
}

type LimitsType struct {
	Ponder         bool
	Infinite       bool
	WhiteTime      int
	BlackTime      int
	WhiteIncrement int
	BlackIncrement int
	MoveTime       int
	MovesToGo      int
	Depth          int
	Nodes          int
	Mate           int
}

type SearchParams struct {
	Positions []backend.Position
	Limits    LimitsType
}

func (e *Engine) GetInfo() (name, version, author string) {
	return "Combusken", "1.1.1", "Marcin Henryk Bartkowiak"
}

func (e *Engine) GetOptions() []EngineOption {
	return []EngineOption{&e.Hash, &e.Threads, &e.PawnHash, &e.MoveOverhead, &e.SyzygyPath, &e.SyzygyProbeDepth}
}

func NewEngine() (ret Engine) {
	ret.Hash = IntOption{"Hash", 4, 1024 * 256, 256}
	ret.Threads = IntOption{"Threads", 1, runtime.NumCPU(), 1}
	ret.PawnHash = IntOption{"PawnHash", 0, 8, 2}
	ret.MoveOverhead = IntOption{"Move Overhead", 0, 10000, 50}
	ret.SyzygyPath = StringOption{"SyzygyPath", "", false}
	ret.SyzygyProbeDepth = IntOption{"SyzygyProbeDepth", 0, 100, 0}
	ret.threads = make([]thread, 1)
	ret.Update = func(SearchInfo) {}
	return
}

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) backend.Move {
	e.fillMoveHistory(searchParams.Positions)
	e.timeManager = newTimeManager(searchParams.Limits, e.MoveOverhead.Val, searchParams.Positions[len(searchParams.Positions)-1].SideToMove)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	if e.hardTimeout() > 0 {
		ctx, cancel = context.WithTimeout(ctx, e.hardTimeout())
	}
	defer cancel()
	e.done = ctx.Done()
	return e.bestMove(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
}

func (e *Engine) fillMoveHistory(positions []backend.Position) {
	e.MovesCount = len(positions) - 1
	moveHistory := make(map[uint64]int)
	for i := len(positions) - 1; i >= 0; i-- {
		moveHistory[positions[i].Key]++
		if positions[i].FiftyMove == 0 {
			break
		}
	}
	e.RepeatedPositions = make(map[uint64]interface{})
	for key, count := range moveHistory {
		if count >= 2 {
			e.RepeatedPositions[key] = struct{}{}
		}
	}
}

func (e *Engine) NewGame() {
	transposition.GlobalTransTable = transposition.NewTransTable(e.Hash.Val)
	e.threads = make([]thread, e.Threads.Val)
	for i := range e.threads {
		e.threads[i].MoveHistory = MoveHistory{}
		e.threads[i].engine = e
	}
	transposition.GlobalPawnKingTable = transposition.NewPKTable(e.PawnHash.Val)
	fathom.MIN_PROBE_DEPTH = e.SyzygyProbeDepth.Val
	if e.SyzygyPath.Dirty {
		fathom.SetPath(e.SyzygyPath.Val)
		e.SyzygyPath.Clean()
	}
	runtime.GC()
}

func (e *Engine) nodes() (sum int) {
	for i := range e.threads {
		sum += e.threads[i].nodes
	}
	return
}

func (t *thread) incNodes() {
	t.nodes++
	if (t.nodes % 255) == 0 {
		select {
		case <-t.engine.done:
			panic(errTimeout)
		default:
		}
	}
}

func (t *thread) getNextMove(pos *backend.Position, depth, height int) backend.Move {
	return t.stack[height].GetNextMove(pos, &t.MoveHistory, depth, height)
}

func (pv *PV) clear() {
	pv.size = 0
}

func (pv *PV) assign(m backend.Move, child *PV) {
	pv.size = 1 + child.size
	pv.items[0] = m
	copy(pv.items[1:], child.Moves())
}

func (pv *PV) Moves() []backend.Move {
	return pv.items[:pv.size]
}
