package engine

import (
	"context"
	"sync"

	. "github.com/mhib/combusken/backend"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const ValueWin = Mate - 150
const ValueLoss = -ValueWin
const SMPCycles = 16

const WindowDepth = 5
const WindowSize = 20

var SkipSize = []int{1, 1, 1, 2, 2, 2, 1, 3, 2, 2, 1, 3, 3, 2, 2, 1}
var SkipDepths = []int{1, 2, 2, 4, 4, 3, 2, 5, 4, 3, 2, 6, 5, 4, 3, 2}

func lossIn(height int) int {
	return -Mate + height
}

func depthToMate(val int) int {
	if val >= ValueWin {
		return Mate - val
	}
	return val - Mate
}

func (t *thread) quiescence(alpha, beta, height int, inCheck bool) int {
	t.incNodes()
	t.stack[height].PV.clear()
	pos := &t.stack[height].position

	if height >= MAX_HEIGHT || t.isDraw(height) {
		return contempt(pos)
	}

	if hashOk, hashValue, _, _, hashFlag := t.engine.TransTable.Get(pos.Key, height); hashOk {
		tmpHashValue := int(hashValue)
		if hashFlag == TransExact || (hashFlag == TransAlpha && tmpHashValue <= alpha) ||
			(hashFlag == TransBeta && tmpHashValue >= beta) {
			return tmpHashValue
		}
	}

	child := &t.stack[height+1].position

	moveCount := 0

	val := Evaluate(pos)

	var evaled []EvaledMove
	if inCheck {
		evaled = pos.GenerateAllMoves(t.stack[height].moves[:])
	} else {
		if val >= beta {
			return beta
		}
		if alpha < val {
			alpha = val
		}
		evaled = pos.GenerateAllCaptures(t.stack[height].moves[:])
	}

	t.EvaluateQsMoves(pos, evaled, inCheck)

	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if (!inCheck && !seeSign(pos, evaled[i].Move)) || !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		moveCount++
		childInCheck := child.IsInCheck()
		val = -t.quiescence(-beta, -alpha, height+1, childInCheck)
		if val > alpha {
			alpha = val
			t.stack[height].PV.assign(evaled[i].Move, &t.stack[height+1].PV)
			if val >= beta {
				return beta
			}
		}
	}

	if moveCount == 0 && inCheck {
		return lossIn(height)
	}

	return alpha
}

func contempt(pos *Position) int {
	return 0
}

func maxMoveToFirst(moves []EvaledMove) {
	maxIdx := 0
	for i := 1; i < len(moves); i++ {
		if moves[i].Value > moves[maxIdx].Value {
			maxIdx = i
		}
	}
	moves[0], moves[maxIdx] = moves[maxIdx], moves[0]
}

func (t *thread) alphaBeta(depth, alpha, beta, height int, inCheck bool) int {
	t.incNodes()
	t.stack[height].PV.clear()
	rootNode := height == 0
	pvNode := alpha != beta+1

	var pos *Position = &t.stack[height].position

	if !rootNode {
		if t.isDraw(height) {
			return contempt(pos)
		}
	}

	var tmpVal int

	alphaOrig := alpha
	hashOk, hashValue, hashDepth, hashMove, hashFlag := t.engine.TransTable.Get(pos.Key, height)
	if hashOk && !pvNode {
		tmpVal = int(hashValue)
		if hashDepth >= uint8(depth) {
			if hashFlag == TransExact || (hashFlag == TransAlpha && tmpVal <= alpha) || (hashFlag == TransBeta && tmpVal >= beta) {
				return tmpVal
			}
		}
	}

	var child *Position = &t.stack[height+1].position

	if pos.LastMove != NullMove && depth >= 4 && !inCheck && !isLateEndGame(pos) {
		pos.MakeNullMove(child)
		reduction := max(1+depth/3, 3)
		tmpVal = -t.alphaBeta(depth-reduction, -beta, -beta+1, height+1, child.IsInCheck())
		if tmpVal >= beta {
			return beta
		}
	}

	if depth == 0 {
		return t.quiescence(alpha, beta, height, inCheck)
	}

	lazyEval := lazyEval{position: pos}
	val := MinInt

	// Internal iterative deepening
	if hashMove == NullMove && !inCheck && ((pvNode && depth >= 6) || (!pvNode && depth >= 8)) {
		var iiDepth int
		if pvNode {
			iiDepth = depth - depth/4 - 1
		} else {
			iiDepth = (depth - 5) / 2
		}
		t.alphaBeta(iiDepth, alpha, beta, height, inCheck)
		_, _, _, hashMove, _ = t.engine.TransTable.Get(pos.Key, height)
	}

	var evaled []EvaledMove
	if rootNode {
		evaled = t.rootMoves
	} else {
		evaled = pos.GenerateAllMoves(t.stack[height].moves[:])
		t.EvaluateMoves(pos, evaled, hashMove, height)
	}
	quietsSearched := t.stack[height].quietsSearched[:0]
	//t.ResetKillers(height)
	bestMove := NullMove
	moveCount := 0
	for i := range evaled {
		if !rootNode {
			if i < 4 {
				maxMoveToFirst(evaled[i:])
			} else if i == 4 {
				sortMoves(evaled[i:])
			}
		}
		if !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		moveCount++
		childInCheck := child.IsInCheck()
		reduction := 0
		if !inCheck && moveCount > 1 && evaled[i].Value <= MinSpecialMoveValue && !evaled[i].Move.IsCaptureOrPromotion() &&
			!childInCheck {
			if depth >= 3 {
				reduction = lmr(depth, moveCount)
				if pvNode {
					reduction--
				}
				reduction = max(0, min(depth-2, reduction))
			} else {
				if moveCount >= 9+3*depth {
					continue
				}
				if lazyEval.Value()+int(PawnValue.Middle)*depth <= alpha {
					continue
				}
			}
		}
		newDepth := depth - 1
		if inCheck && seeSign(pos, evaled[i].Move) {
			newDepth++
		}

		if !evaled[i].Move.IsCaptureOrPromotion() {
			quietsSearched = append(quietsSearched, evaled[i].Move)
		}
		if reduction > 0 || (!pvNode && moveCount > 1 && evaled[i].Value < MinSpecialMoveValue) {
			tmpVal = -t.alphaBeta(newDepth-reduction, -(alpha + 1), -alpha, height+1, childInCheck)
			if tmpVal <= alpha {
				continue
			}
		}
		tmpVal = -t.alphaBeta(newDepth, -beta, -alpha, height+1, childInCheck)

		if tmpVal > val {
			val = tmpVal
			if val > alpha {
				alpha = val
				bestMove = evaled[i].Move
				t.stack[height].PV.assign(evaled[i].Move, &t.stack[height+1].PV)
				if alpha >= beta {
					break
				}
			}
		}
	}

	if moveCount == 0 {
		if inCheck {
			return lossIn(height)
		}
		return contempt(pos)
	}

	if bestMove != NullMove && !bestMove.IsCaptureOrPromotion() {
		t.Update(pos, quietsSearched, bestMove, depth, height)
	}

	var flag int
	if alpha == alphaOrig {
		flag = TransAlpha
	} else if alpha >= beta {
		flag = TransBeta
	} else {
		flag = TransExact
	}
	t.engine.TransTable.Set(pos.Key, alpha, depth, bestMove, flag, height)
	return alpha
}

func (t *thread) isDraw(height int) bool {
	var pos *Position = &t.stack[height].position

	if t.engine.MovesCount < 50 {
		return false
	}

	if pos.FiftyMove > 100 {
		return true
	}

	if (pos.Pawns|pos.Rooks|pos.Queens) == 0 && !MoreThanOne(pos.Knights|pos.Bishops) {
		return true
	}

	for i := height - 1; i >= 0; i-- {
		descendant := &t.stack[i].position
		if descendant.Key == pos.Key {
			return true
		}
		if descendant.FiftyMove == 0 || descendant.LastMove == NullMove {
			return false
		}
	}

	if t.engine.MoveHistory[pos.Key] >= 2 {
		return true
	}

	return false
}

type result struct {
	Move
	value int
	depth int
	moves []Move
}

func (t *thread) aspirationWindow(depth int, resultChan chan result) {
	var alpha int
	var beta int
	delta := WindowSize
	if depth >= WindowDepth {
		alpha = max(-Mate, t.lastValue-delta)
		beta = min(Mate, t.lastValue+delta)
	} else {
		alpha = -Mate
		beta = Mate
	}

	for {
		t.sortRootMoves()
		value := t.alphaBeta(depth, alpha, beta, 0, t.stack[0].position.IsInCheck())
		if value > alpha && value < beta {
			resultChan <- result{t.stack[0].PV.items[0], value, depth, t.stack[0].PV.Moves()}
			return
		}
		if value <= alpha {
			beta = (alpha + beta) / 2
			alpha = max(-Mate, value-delta)
		} else if value >= beta {
			alpha = (alpha + beta) / 2
			beta = min(Mate, value+delta)
		}
		delta += delta / 2
	}
}

func moveToFirst(moves []EvaledMove, idx int) {
	if idx == 0 {
		return
	}
	move := moves[idx]
	for i := idx; idx > 0; idx-- {
		moves[i] = moves[i-1]
	}
	moves[0] = move
}

func (e *Engine) singleThreadBestMove(ctx context.Context) Move {
	var lastBestMove Move
	thread := e.threads[0]
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		go func(depth int) {
			defer recoverFromTimeout()
			thread.aspirationWindow(depth, resultChan)
		}(i)
		select {
		case <-ctx.Done():
			return lastBestMove
		case res := <-resultChan:
			e.callUpdate(SearchInfo{res.value, i, thread.nodes, res.moves})
			if res.value >= ValueWin && depthToMate(res.value) <= i {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			if i >= MAX_HEIGHT {
				return res.Move
			}
			if e.isSoftTimeout(i, thread.nodes) {
				return res.Move
			}
			lastBestMove = res.Move
		}
	}
}

func (t *thread) iterativeDeepening(resultChan chan result, idx int) {
	mainThread := idx == 0
	cycle := idx % SMPCycles
	for depth := 1; depth < MAX_HEIGHT; depth++ {
		t.aspirationWindow(depth, resultChan)
		if !mainThread && (depth+cycle)%SkipDepths[cycle] == 0 {
			depth += SkipSize[cycle]
		}
	}
}

func (t *thread) sortRootMoves() {
	pos := t.stack[0].position
	ordMove := NullMove
	if hashOk, _, _, hashMove, _ := t.engine.TransTable.Get(pos.Key, 0); hashOk {
		ordMove = hashMove
	}
	t.EvaluateMoves(&t.stack[0].position, t.rootMoves, ordMove, 0)
	sortMoves(t.rootMoves)
}

func (e *Engine) bestMove(ctx context.Context, pos *Position) Move {
	for i := range e.threads {
		e.threads[i].stack[0].position = *pos
		e.threads[i].nodes = 0
		e.threads[i].rootMoves = pos.GenerateAllLegalMoves(e.threads[i].stack[0].moves[:])
	}

	if e.Threads.Val == 1 {
		return e.singleThreadBestMove(ctx)
	}

	var wg = &sync.WaitGroup{}
	resultChan := make(chan result)
	for i := range e.threads {
		wg.Add(1)
		go func(idx int) {
			defer recoverFromTimeout()
			e.threads[idx].iterativeDeepening(resultChan, idx)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	prevDepth := 0
	var lastBestMove Move
	for {
		select {
		case <-e.done:
			return lastBestMove
		case res := <-resultChan:
			if res.depth <= prevDepth {
				continue
			}
			nodes := e.nodes()
			e.callUpdate(SearchInfo{res.value, res.depth, nodes, res.moves})
			if res.value >= ValueWin && depthToMate(res.value) <= res.depth {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			if res.depth >= MAX_HEIGHT {
				return res.Move
			}
			if e.isSoftTimeout(res.depth, nodes) {
				return res.Move
			}
			lastBestMove = res.Move
			prevDepth = res.depth
		}
	}
}

func cloneMoves(src []Move) []Move {
	dst := make([]Move, len(src))
	copy(dst, src)
	return dst
}

func cloneEvaledMoves(src []EvaledMove) []EvaledMove {
	dst := make([]EvaledMove, len(src))
	copy(dst, src)
	return dst
}

func recoverFromTimeout() {
	err := recover()
	if err != nil && err != errTimeout {
		panic(err)
	}
}

type lazyEval struct {
	position *Position
	hasValue bool
	value    int
}

func (le *lazyEval) Value() int {
	if !le.hasValue {
		le.value = Evaluate(le.position)
		le.hasValue = true
	}
	return le.value
}

// Gaps from Best Increments for the Average Case of Shellsort, Marcin Ciura.
var shellSortGaps = [...]int{23, 10, 4, 1}

func sortMoves(moves []EvaledMove) {
	for _, gap := range shellSortGaps {
		for i := gap; i < len(moves); i++ {
			j, t := i, moves[i]
			for ; j >= gap && moves[j-gap].Value < t.Value; j -= gap {
				moves[j] = moves[j-gap]
			}
			moves[j] = t
		}
	}
}

func lmr(d, m int) int {
	switch {
	case d >= 5 && m >= 16:
		return 3
	case d >= 4 && m >= 9:
		return 2
	case d >= 3 && m >= 4:
		return 1
	default:
		return 0
	}
}
