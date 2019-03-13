package engine

import "context"
import "time"
import "errors"
import "github.com/mhib/combusken/backend"

const MAX_HEIGHT = 127
const STACK_SIZE = MAX_HEIGHT + 1

var errTimeout = errors.New("Search timeout")

type IntUciOption struct {
	Name string
	Min  int
	Max  int
	Val  int
}

type Engine struct {
	Hash     IntUciOption
	done     <-chan struct{}
	timedOut chan bool
	TransTable
	MoveHistory  map[uint64]int
	EvalHistory  [64][64]int
	MovesCount   int
	KillerMoves  [STACK_SIZE][2]backend.Move
	CounterMoves [64][64]backend.Move
	Update       func(SearchInfo)
	Nodes        int
	Stack        [STACK_SIZE]StackEntry
	timeManager
}

type SearchInfo struct {
	Score int
	Depth int
	Nodes int
	PV
}

type StackEntry struct {
	position backend.Position
	PV
	moves [256]backend.EvaledMove
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
	return "Combusken", "0.0.2", "Marcin Henryk Bartkowiak"
}

func (e *Engine) GetOptions() []*IntUciOption {
	return []*IntUciOption{&e.Hash}
}

func NewEngine() (ret Engine) {
	ret.Hash = IntUciOption{"Hash", 4, 1024, 64}
	ret.TransTable = NewTransTable(ret.Hash.Val)
	return
}

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) backend.Move {
	e.cleanBeforeSearch()
	e.timedOut = make(chan bool, 1)
	e.fillMoveHistory(searchParams.Positions)
	e.timeManager = newBlankTimeManager()
	if searchParams.Limits.MoveTime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(searchParams.Limits.MoveTime)*time.Millisecond)
		defer cancel()
		e.done = ctx.Done()
		return e.TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	} else if searchParams.Limits.Depth > 0 {
		e.done = ctx.Done()
		return e.DepthSearch(&searchParams.Positions[len(searchParams.Positions)-1], searchParams.Limits.Depth)
	} else if searchParams.Limits.Infinite {
		e.done = ctx.Done()
		return e.TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	} else if searchParams.Limits.WhiteTime > 0 {
		var cancel context.CancelFunc
		e.timeManager = newTimeManager(searchParams.Limits, searchParams.Positions[len(searchParams.Positions)-1].WhiteMove)
		ctx, cancel = context.WithTimeout(ctx, e.hardTimeout)
		defer cancel()
		e.done = ctx.Done()
		return e.TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		e.done = ctx.Done()
		return e.TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	}
}

func (e *Engine) fillMoveHistory(positions []backend.Position) {
	e.MovesCount = len(positions) - 1
	e.MoveHistory = make(map[uint64]int)
	for i := len(positions) - 1; i >= 0; i-- {
		e.MoveHistory[positions[i].Key]++
		if positions[i].FiftyMove == 0 {
			break
		}
	}
}

func (e *Engine) cleanBeforeSearch() {
	//e.TransTable.Clear()
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			e.EvalHistory[y][x] = 0
		}
	}
	for y := 0; y < MAX_HEIGHT; y++ {
		for x := 0; x < 2; x++ {
			e.KillerMoves[y][x] = backend.NullMove
		}
	}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			e.CounterMoves[y][x] = backend.NullMove
		}
	}
}

func (e *Engine) callUpdate(s SearchInfo) {
	if e.Update != nil {
		e.Update(s)
	}
}

func (e *Engine) incNodes() {
	e.Nodes++
	if (e.Nodes % 255) == 0 {
		select {
		case <-e.timedOut:
			panic(errTimeout)
		default:
		}
	}
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
