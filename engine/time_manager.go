package engine

import (
	"math"
	"time"

	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

type timeElapser struct {
	startedAt time.Time
}

func (elapser *timeElapser) getElapsedTime() time.Duration {
	return time.Since(elapser.startedAt)
}

type timeManager interface {
	hardTimeout() time.Duration
	isSoftTimeout(depth, nodes int) bool
	updateTime(depth, score int)
	getElapsedTime() time.Duration
}

type depthMoveTimeManager struct {
	timeElapser
	duration int
	depth    int
}

func (manager *depthMoveTimeManager) hardTimeout() time.Duration {
	if manager.duration > 0 {
		return time.Duration(manager.duration) * time.Millisecond
	}
	return 0
}

func (manager *depthMoveTimeManager) isSoftTimeout(depth, nodes int) bool {
	return manager.depth > 0 && depth >= manager.depth
}

func (manager *depthMoveTimeManager) updateTime(int, int) {
}

type tournamentTimeManager struct {
	timeElapser
	hard      time.Duration
	ideal     time.Duration
	lastScore int
}

func (manager *tournamentTimeManager) hardTimeout() time.Duration {
	return manager.hard
}

func (manager *tournamentTimeManager) isSoftTimeout(int, int) bool {
	return time.Since(manager.startedAt) >= manager.ideal
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

}

func newTournamentTimeManager(startedAt time.Time, limits LimitsType, overhead, sideToMove int) *tournamentTimeManager {
	res := &tournamentTimeManager{timeElapser: timeElapser{startedAt: startedAt}}
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
	res.ideal = time.Duration(Min(ideal, limit-overhead)) * time.Millisecond
	res.hard = time.Duration(Min(hard, limit-overhead)) * time.Millisecond
	return res
}

func newTimeManager(limits LimitsType, overhead int, sideToMove int) timeManager {
	startedAt := time.Now()
	if limits.WhiteTime > 0 || limits.BlackTime > 0 {
		return newTournamentTimeManager(startedAt, limits, overhead, sideToMove)
	} else {
		return &depthMoveTimeManager{timeElapser{startedAt: startedAt}, limits.MoveTime, limits.Depth}
	}
}
