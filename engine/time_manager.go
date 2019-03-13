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

type timeoutStrategy func(*timeManager) bool

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

func (manager *timeManager) isSoftTimeout() bool {
	return manager.timeoutStrategy(manager)
}

func compareSoftTimeout(manager *timeManager) bool {
	return time.Now().Sub(manager.startedAt) >= manager.softTimeout
}

func newTimeManager(limits LimitsType, whiteMove bool) (res timeManager) {
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

func newBlankTimeManager() (res timeManager) {
	res.timeoutStrategy = func(*timeManager) bool {
		return false
	}
	return
}
