package engine

import "context"
import "time"
import "github.com/mhib/combusken/backend"

type Engine struct {
	done <-chan struct{}
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
	return "Combusken", "0.0.1", "Marcin Henryk Bartkowiak"
}

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) backend.Move {
	if searchParams.Limits.MoveTime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(searchParams.Limits.MoveTime)*time.Millisecond)
		defer cancel()
		e.done = ctx.Done()
		return TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	} else if searchParams.Limits.Depth > 0 {
		e.done = ctx.Done()
		return DepthSearch(&searchParams.Positions[len(searchParams.Positions)-1], searchParams.Limits.Depth)
	} else if searchParams.Limits.Infinite {
		e.done = ctx.Done()
		return TimeSearch(ctx, &searchParams.Positions[len(searchParams.Positions)-1])
	}
	return backend.Move(0)
}
