package engine

import (
	"context"
	"math"
	"time"

	. "github.com/mhib/combusken/chess"
	. "github.com/mhib/combusken/utils"
)

type timeElapser struct {
	startedAt time.Time
}

func (elapser *timeElapser) getElapsedTime() time.Duration {
	return time.Since(elapser.startedAt)
}

func newTimeElapser() timeElapser {
	return timeElapser{startedAt: time.Now()}
}

type timeManager interface {
	updateTime(depth, score int)
	getElapsedTime() time.Duration
}

type depthMoveTimeManager struct {
	timeElapser
	duration time.Duration
	depth    int
	cancel   context.CancelFunc
}

func (manager *depthMoveTimeManager) updateTime(depth, nodes int) {
	if manager.depth > 0 && depth >= manager.depth {
		manager.cancel()
	}
}

func newDepthMoveTimeManager(ctx context.Context, limits *LimitsType) (context.Context, *depthMoveTimeManager) {
	var cancel context.CancelFunc
	duration := time.Duration(limits.MoveTime) * time.Millisecond
	if duration > 0 {
		ctx, cancel = context.WithTimeout(ctx, duration)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	return ctx, &depthMoveTimeManager{newTimeElapser(), duration, limits.Depth, cancel}
}

type tournamentTimeManager struct {
	timeElapser
	hard      time.Duration
	ideal     time.Duration
	lastScore int
	cancel    context.CancelFunc
}

func (manager *tournamentTimeManager) updateTime(depth, score int) {
	lastScore := manager.lastScore
	manager.lastScore = score
	if depth < 4 {
		return
	}
	if lastScore > score {
		multiplier := math.Min(float64(lastScore-score)*0.042709881489053105, 3.525357845675907)
		manager.ideal += time.Duration(float64(manager.ideal) * multiplier / 16)
	}

	since := time.Since(manager.startedAt)
	if since >= manager.ideal {
		manager.cancel()
	}
}

func newTournamentTimeManager(ctx context.Context, limits *LimitsType, overhead int, ponder bool, sideToMove int) (context.Context, *tournamentTimeManager) {
	res := &tournamentTimeManager{timeElapser: newTimeElapser()}
	var limit, inc int
	if sideToMove == White {
		limit, inc = limits.WhiteTime, limits.WhiteIncrement
	} else {
		limit, inc = limits.BlackTime, limits.BlackIncrement
	}
	movesToGo := limits.MovesToGo
	var ideal, hard int

	if movesToGo > 0 {
		ideal = ((((limit - overhead) / (movesToGo + 4)) + inc) * 14) / 16
		hard = (((limit - overhead) / (movesToGo + 6)) + inc) * 4
	} else {
		ideal = ((limit - overhead) + 25*inc) / 45
		hard = 5 * ((limit - overhead) + 25*inc) / 45
	}
	if ponder {
		ideal += ideal / 4
		hard += hard / 16
	}
	res.ideal = time.Duration(Min(ideal, limit-overhead)) * time.Millisecond
	res.hard = time.Duration(Min(hard, limit-overhead)) * time.Millisecond
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, res.hard)
	res.cancel = cancel
	return ctx, res
}

type infiniteTimeManager struct {
	timeElapser
}

func (infiniteTimeManager) updateTime(int, int) {
}

func newTimeManager(ctx context.Context, limits *LimitsType, overhead int, ponder bool, sideToMove int) (context.Context, timeManager) {
	if limits.Infinite {
		return ctx, &infiniteTimeManager{newTimeElapser()}
	}
	if limits.WhiteTime > 0 || limits.BlackTime > 0 {
		return newTournamentTimeManager(ctx, limits, overhead, ponder, sideToMove)
	} else {
		return newDepthMoveTimeManager(ctx, limits)
	}
}
