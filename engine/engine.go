package engine

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/mhib/combusken/chess"
	"github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/fathom"
	"github.com/mhib/combusken/transposition"

	. "github.com/mhib/combusken/utils"
)

const MaxHeight = 127
const StackSize = MaxHeight + 1
const MaxMoves = 256
const MaxMultiPV = 128

var errTimeout = errors.New("Search timeout")

type Engine struct {
	Hash              IntOption
	Threads           IntOption
	MoveOverhead      IntOption
	PawnHash          IntOption
	SyzygyPath        StringOption
	SyzygyProbeDepth  IntOption
	MultiPV           IntOption
	Ponder            CheckOption
	done              <-chan struct{}
	RepeatedPositions map[uint64]interface{}
	MovesCount        int
	Update            func(*SearchInfo)
	timeManager
	threads         []thread
	multiPVExcluded [MaxMultiPV]chess.Move
	lastValues      [MaxMultiPV]int16
}

type thread struct {
	engine *Engine
	MoveHistory
	evaluation.EvaluationContext
	index           int
	nodes           int
	tbhits          int
	disableNmpColor int
	seldepth        int
	stack           [StackSize]StackEntry
}

type ReportScore struct {
	Mate      int
	Centipawn int
}

func newReportScore(score int) ReportScore {
	if score >= ValueWin {
		return ReportScore{Mate: (Mate - score + 1) / 2}
	} else if score <= ValueLoss {
		return ReportScore{Mate: (-Mate - score) / 2}
	} else {
		return ReportScore{Centipawn: score}
	}
}

type SearchInfo struct {
	Score    ReportScore
	Depth    int
	SelDepth int
	MultiPV  int
	Nodes    int
	Nps      int
	Duration int
	Tbhits   int
	Moves    []chess.Move
}

type StackEntry struct {
	MoveProvider
	PV
	quietsSearched [MaxMoves]chess.Move
	noisySearched  [MaxMoves]chess.Move
	position       chess.Position
	evaluation     int16
}

func (t *thread) getEvaluation(height int) int16 {
	return t.stack[height].evaluation
}

func (t *thread) setEvaluation(height int, eval int16) {
	t.stack[height].evaluation = eval
}

type PV struct {
	size  int
	items [StackSize]chess.Move
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
	Positions []chess.Position
	Limits    LimitsType
}

func (e *Engine) GetInfo() (name, version, author string) {
	return "Combusken", "1.9.9_dev", "Marcin Henryk Bartkowiak"
}

func (e *Engine) GetOptions() []EngineOption {
	return []EngineOption{&e.Hash, &e.Threads, &e.PawnHash, &e.MoveOverhead, &e.SyzygyPath, &e.SyzygyProbeDepth, &e.Ponder, &e.MultiPV}
}

func NewEngine() (ret Engine) {
	ret.Hash = IntOption{"Hash", 4, 1024 * 256, 16, 16}
	ret.Threads = IntOption{"Threads", 1, runtime.NumCPU(), 1, 1}
	ret.PawnHash = IntOption{"PawnHash", 0, 8, 2, 2}
	ret.MoveOverhead = IntOption{"Move Overhead", 0, 10000, 50, 50}
	ret.SyzygyPath = StringOption{"SyzygyPath", "", "", false}
	ret.SyzygyProbeDepth = IntOption{"SyzygyProbeDepth", 0, 100, 0, 0}
	ret.MultiPV = IntOption{"MultiPV", 1, MaxMultiPV, 1, 1}
	ret.Ponder = CheckOption{"Ponder", false, false}
	ret.threads = make([]thread, 1)
	ret.Update = func(*SearchInfo) {}
	return
}

func (e *Engine) SetUpdate(update func(*SearchInfo)) {
	e.Update = update
}

func (e *Engine) IsMoveExcluded(move chess.Move) bool {
	for i := 0; e.multiPVExcluded[i] != chess.NullMove && i < MaxMultiPV; i++ {
		if e.multiPVExcluded[i] == move {
			return true
		}
	}
	return false
}

func (e *Engine) Search(ctx context.Context, ponderCtx context.Context, searchParams SearchParams) (bestMove, ponderMove chess.Move) {
	e.fillMoveHistory(searchParams.Positions)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	var timeoutCtx context.Context
	timeoutCtx, e.timeManager = newTimeManager(ctx, &searchParams.Limits, e.MoveOverhead.Val, e.Ponder.Val, searchParams.Positions[len(searchParams.Positions)-1].SideToMove)
	if searchParams.Limits.Ponder {
		definePonderCancellation(ponderCtx, timeoutCtx, cancel)
	} else {
		ctx = timeoutCtx
	}
	e.done = ctx.Done()
	var wg sync.WaitGroup
	wg.Add(len(e.threads))
	bestMove, ponderMove = e.bestMove(ctx, ponderCtx, &wg, &searchParams.Positions[len(searchParams.Positions)-1])
	cancel()
	wg.Wait()
	return
}

func definePonderCancellation(ponderCtx context.Context, timeoutCtx context.Context, cancel context.CancelFunc) {
	go func() {
		<-ponderCtx.Done()
		<-timeoutCtx.Done()
		cancel()
	}()
}

func (e *Engine) fillMoveHistory(positions []chess.Position) {
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
		e.threads[i].index = i
		e.threads[i].MoveHistory = MoveHistory{}
		e.threads[i].engine = e
	}
	evaluation.GlobalPawnKingTable = evaluation.NewPawnKingTable(e.PawnHash.Val)
	fathom.MinProbeDepth = e.SyzygyProbeDepth.Val
	if e.SyzygyPath.Dirty {
		fathom.SetPath(e.SyzygyPath.Val)
		e.SyzygyPath.Clean()
	}
	runtime.GC()
}

func (e *Engine) aggregatesInfo() (nodes, tbhits int) {
	for i := range e.threads {
		nodes += e.threads[i].nodes
		tbhits += e.threads[i].tbhits
	}
	return
}

func (t *thread) incNodes() {
	t.nodes++
	if (t.nodes % 511) == 0 {
		select {
		case <-t.engine.done:
			panic(errTimeout)
		default:
		}
	}
}

func (t *thread) getNextMove(pos *chess.Position, depth, height int) chess.Move {
	return t.stack[height].GetNextMove(pos, &t.MoveHistory, depth, height)
}

func (t *thread) skipQuiets(height int) {
	t.stack[height].skipQuiets()
}

func (t *thread) isMainThread() bool {
	return t.index == 0
}

func (t *thread) getRootMovesBuffer() []chess.EvaledMove {
	return t.stack[0].MoveProvider.Moves[:]
}

func (pv *PV) clear() {
	pv.size = 0
}

func (pv *PV) assign(m chess.Move, child *PV) {
	pv.size = 1 + child.size
	pv.items[0] = m
	copy(pv.items[1:], child.Moves())
}

func (pv *PV) Moves() []chess.Move {
	return pv.items[:pv.size]
}
