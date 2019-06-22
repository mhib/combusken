package engine

import "time"

type timeManager interface {
	hardTimeout() time.Duration
	isSoftTimeout(depth, nodes int) bool
	updateTime(depth, score int)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

type blankTimeManager struct {
}

func (*blankTimeManager) hardTimeout() time.Duration {
	return 0
}

func (*blankTimeManager) isSoftTimeout(depth, nodes int) bool {
	return false
}

func (*blankTimeManager) updateTime(int, int) {
}

type moveTimeTimeManager struct {
	blankTimeManager
	duration int
}

func (manager *moveTimeTimeManager) hardTimeout() time.Duration {
	return time.Duration(manager.duration) * time.Millisecond
}

type depthTimeManager struct {
	blankTimeManager
	depth int
}

func (manager *depthTimeManager) isSoftTimeout(depth, nodes int) bool {
	return depth >= manager.depth
}

type tournamentTimeManager struct {
	hard      time.Duration
	ideal     time.Duration
	startedAt time.Time
	lastScore int
}

func (manager *tournamentTimeManager) hardTimeout() time.Duration {
	return manager.hard
}

func (manager *tournamentTimeManager) isSoftTimeout(int, int) bool {
	return time.Now().Sub(manager.startedAt) >= manager.ideal
}

func (manager *tournamentTimeManager) updateTime(depth, score int) {
	lastScore := manager.lastScore
	manager.lastScore = score
	if depth < 4 {
		return
	}

	if lastScore > score+10 {
		manager.ideal += manager.ideal / 20
	}
	if lastScore > score+20 {
		manager.ideal += manager.ideal / 20
	}
	if lastScore > score+40 {
		manager.ideal += manager.ideal / 20
	}

	if lastScore+15 < score {
		manager.ideal += manager.ideal / 40
	}
	if lastScore+30 < score {
		manager.ideal += manager.ideal / 20
	}
}

func newTournamentTimeManager(limits LimitsType, overhead int, whiteMove bool) *tournamentTimeManager {
	res := &tournamentTimeManager{startedAt: time.Now()}
	var limit, inc int
	if whiteMove {
		limit, inc = limits.WhiteTime, limits.WhiteIncrement
	} else {
		limit, inc = limits.BlackTime, limits.BlackIncrement
	}
	movesToGo := limits.MovesToGo
	var ideal, hard int

	if movesToGo > 0 {
		ideal = (((limit/movesToGo + 5) + inc) * 3) / 4
		hard = ((limit/movesToGo + 7) + inc) * 4
	} else {
		ideal = (limit + 25*inc) / 50
		hard = 10 * (limit + 25*inc) / 50
	}
	res.ideal = time.Duration(min(ideal, limit-overhead)) * time.Millisecond
	res.hard = time.Duration(min(hard, limit-overhead)) * time.Millisecond
	return res
}

func newTimeManager(limits LimitsType, overhead int, whiteMove bool) timeManager {
	if limits.WhiteTime > 0 || limits.BlackTime > 0 {
		return newTournamentTimeManager(limits, overhead, whiteMove)
	} else if limits.MoveTime > 0 {
		return &moveTimeTimeManager{duration: limits.MoveTime}
	} else if limits.Depth > 0 {
		return &depthTimeManager{depth: limits.Depth}
	} else if limits.Infinite {
		return &blankTimeManager{}
	} else {
		return &moveTimeTimeManager{duration: 1000}
	}
}
