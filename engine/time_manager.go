package engine

import "time"

const defaultMovesToGo = 35
const buffer = 50

type timeManager struct {
	timeoutStrategy
	startedAt   time.Time
	softTimeout time.Duration
	hardTimeout time.Duration
}

type timeoutStrategy func(manager *timeManager, depth, nodes int) bool

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

func (manager *timeManager) isSoftTimeout(depth, nodes int) bool {
	return manager.timeoutStrategy(manager, depth, nodes)
}

func compareSoftTimeout(manager *timeManager, depth, nodes int) bool {
	return time.Now().Sub(manager.startedAt) >= manager.softTimeout
}

func neverSoftTimeout(manager *timeManager, depth, nodes int) bool {
	return false
}

func newTimeManager(limits LimitsType, whiteMove bool) timeManager {
	if limits.MoveTime > 0 {
		return newMoveTimeManager(limits.MoveTime)
	} else if limits.Depth > 0 {
		return newDepthTimeManager(limits.Depth)
	} else if limits.Infinite {
		return timeManager{timeoutStrategy: neverSoftTimeout}
	} else if limits.WhiteTime > 0 || limits.BlackTime > 0 {
		return newTournamentTimeManager(limits, whiteMove)
	} else {
		return newMoveTimeManager(1000)
	}
}

func newDepthTimeManager(limitsDepth int) (res timeManager) {
	res.timeoutStrategy = func(manager *timeManager, depth, nodes int) bool {
		return depth >= limitsDepth
	}
	return
}

func newMoveTimeManager(duration int) (res timeManager) {
	res.startedAt = time.Now()
	res.hardTimeout = time.Duration(duration) * time.Millisecond
	res.timeoutStrategy = neverSoftTimeout
	return
}

func newTournamentTimeManager(limits LimitsType, whiteMove bool) (res timeManager) {
	var limit, inc int
	if whiteMove {
		limit, inc = limits.WhiteTime, limits.WhiteIncrement
	} else {
		limit, inc = limits.BlackTime, limits.BlackIncrement
	}
	movesToGo := limits.MovesToGo
	if movesToGo == 0 {
		movesToGo = defaultMovesToGo
	}
	ideal := limit / movesToGo
	ideal += inc
	ensureNoFlag := max(limit-buffer, 0)
	res.hardTimeout = time.Duration(min(ideal*2, ensureNoFlag)) * time.Millisecond
	res.softTimeout = res.hardTimeout / 4
	res.startedAt = time.Now()
	res.timeoutStrategy = compareSoftTimeout
	return
}
