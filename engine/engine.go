package engine

import "context"
import "time"
import "github.com/mhib/combusken/backend"

const MAX_HEIGHT = 128

type Engine struct {
	done <-chan struct{}
	TransTable
	MoveHistory  map[uint64]int
	EvalHistory  [64][64]int
	MovesCount   int
	KillerMoves  [MAX_HEIGHT][2]backend.Move
	CounterMoves [64][64]backend.Move
	Update       func(SearchInfo)
	Nodes        int
}

type SearchInfo struct {
	Score int
	Depth int
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

func NewEngine() (ret Engine) {
	ret.TransTable = NewTransTable()
	return
}

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) backend.Move {
	e.cleanBeforeSearch()
	e.fillMoveHistory(searchParams.Positions)
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
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
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
