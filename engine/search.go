package engine

import (
	"context"
	"math/rand"
	"sync"

	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/evaluation"
	"github.com/mhib/combusken/fathom"
	"github.com/mhib/combusken/transposition"
	. "github.com/mhib/combusken/utils"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const ValueWin = Mate - 150
const ValueLoss = -ValueWin

const seePruningDepth = 8
const seeQuietMargin = -80
const seeNoisyMargin = -18

const reverseFutilityPruningDepth = 6
const reverseFutilityPruningMargin = 90

const moveCountPruningDepth = 8
const futilityPruningDepth = 8
const probCutDepth = 5
const probCutMargin = 189

const SMPCycles = 16

const WindowSize = 25
const WindowDepth = 6

const QSDepthChecks = 0
const QSDepthNoChecks = -1

var SkipSize = []int{1, 1, 1, 2, 2, 2, 1, 3, 2, 2, 1, 3, 3, 2, 2, 1}
var SkipDepths = []int{1, 2, 2, 4, 4, 3, 2, 5, 4, 3, 2, 6, 5, 4, 3, 2}

var PawnValueMiddle = PawnValue.Middle()

func lossIn(height int) int {
	return -Mate + height
}

func winIn(height int) int {
	return Mate - height
}

func depthToMate(val int) int {
	if val >= ValueWin {
		return Mate - val
	}
	return val - Mate
}

func (t *thread) quiescence(depth, alpha, beta, height int, inCheck bool) int {
	t.incNodes()
	t.stack[height].PV.clear()
	pos := &t.stack[height].position
	alphaOrig := alpha

	if height >= MAX_HEIGHT || t.isDraw(height) {
		return t.contempt(pos, depth)
	}

	var ttDepth int
	if inCheck || depth >= QSDepthChecks {
		ttDepth = QSDepthChecks
	} else {
		ttDepth = QSDepthNoChecks
	}
	hashOk, hashValue, hashDepth, hashMove, hashFlag := transposition.GlobalTransTable.Get(pos.Key)
	if hashOk && int(hashDepth) >= ttDepth {
		tmpHashValue := int(transposition.ValueFromTrans(hashValue, height))
		if hashFlag == TransExact || (hashFlag == TransAlpha && tmpHashValue <= alpha) ||
			(hashFlag == TransBeta && tmpHashValue >= beta) {
			return tmpHashValue
		}
	}

	child := &t.stack[height+1].position

	bestMove := NullMove

	moveCount := 0

	val := Evaluate(pos)

	var evaled []EvaledMove
	if inCheck {
		evaled = pos.GenerateAllMoves(t.stack[height].moves[:])
	} else {
		// Early return if not in check and evaluation exceeded beta
		if val >= beta {
			return beta
		}
		if alpha < val {
			alpha = val
		}
		evaled = pos.GenerateAllCaptures(t.stack[height].moves[:])
	}

	t.EvaluateQsMoves(pos, evaled, hashMove, inCheck)

	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		// Ignore move with negative SEE unless in check
		if (!inCheck && !SeeSign(pos, evaled[i].Move)) || !pos.MakeMove(evaled[i].Move, child) {
			continue
		}

		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		t.SetCurrentMove(height, evaled[i].Move)
		moveCount++
		childInCheck := child.IsInCheck()
		val = -t.quiescence(depth-1, -beta, -alpha, height+1, childInCheck)
		if val > alpha {
			alpha = val
			bestMove = evaled[i].Move
			if val >= beta {
				break
			}
			t.stack[height].PV.assign(evaled[i].Move, &t.stack[height+1].PV)
		}
	}

	if moveCount == 0 && inCheck {
		return lossIn(height)
	}

	var flag int
	if alpha == alphaOrig {
		flag = TransAlpha
	} else if alpha >= beta {
		flag = TransBeta
	} else {
		flag = TransExact
	}

	transposition.GlobalTransTable.Set(pos.Key, transposition.ValueToTrans(alpha, height), ttDepth, bestMove, flag)

	return alpha
}

// Currently draws are scored as 0 +/- 1 randomly
func (t *thread) contempt(pos *Position, depth int) int {
	if depth < 4 {
		return 0
	}
	return 2*(t.nodes&1) - 1
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

func moveToFirst(moves []EvaledMove, move Move) {
	for i := 0; i < len(moves); i++ {
		if moves[i].Move == move {
			moves[0], moves[i] = moves[i], moves[0]
			return
		}
	}
}

func moveCountPruning(improving, depth int) int {
	return (5+depth*depth)*(1+improving)/2 - 1
}

func (t *thread) alphaBeta(depth, alpha, beta, height int, inCheck bool) int {
	t.incNodes()
	t.stack[height].PV.clear()

	var pos *Position = &t.stack[height].position

	if height >= MAX_HEIGHT || t.isDraw(height) {
		return t.contempt(pos, depth)
	}

	var tmpVal int

	// Node is not pv if it is searched with null window
	pvNode := alpha != beta-1

	// Mate distance pruning
	alpha = Max(lossIn(height), alpha)
	beta = Min(winIn(height+1), beta)
	if alpha >= beta {
		return alpha
	}

	alphaOrig := alpha
	hashOk, hashValue, hashDepth, hashMove, hashFlag := transposition.GlobalTransTable.Get(pos.Key)
	hashValue = transposition.ValueFromTrans(hashValue, height)
	if hashOk {
		tmpVal = int(hashValue)
		// Hash pruning
		if hashDepth >= int16(depth) && (depth == 0 || !pvNode) {
			if hashFlag == TransExact {
				return tmpVal
			}
			if hashFlag == TransAlpha && tmpVal <= alpha {
				return alpha
			}
			if hashFlag == TransBeta && tmpVal >= beta {
				return beta
			}
		}
	}

	// Probe tablebase
	if fathom.IsWDLProbeable(pos, depth) {
		if tbResult := fathom.ProbeWDL(pos, depth); tbResult != fathom.TB_RESULT_FAILED {
			var ttBound int
			if tbResult == fathom.TB_LOSS {
				tmpVal = ValueLoss + height + 1
				ttBound = TransAlpha
			} else if tbResult == fathom.TB_WIN {
				tmpVal = ValueWin - height - 1
				ttBound = TransBeta
			} else {
				tmpVal = 0
				ttBound = TransExact
			}
			if ttBound == TransExact || ttBound == TransBeta && tmpVal >= beta || ttBound == TransAlpha && tmpVal <= alpha {
				transposition.GlobalTransTable.Set(pos.Key, int16(tmpVal), MAX_HEIGHT-1, NullMove, ttBound)
				return tmpVal
			}
		}
	}

	var child *Position = &t.stack[height+1].position

	if depth <= 0 {
		return t.quiescence(0, alpha, beta, height, inCheck)
	}

	t.stack[height].InvalidateEvaluation()
	t.MoveEvaluator.ResetKillers(height + 1)

	// Reverse futility pruning
	if !pvNode && depth < reverseFutilityPruningDepth && int(t.stack[height].Evaluation())-reverseFutilityPruningMargin*depth >= beta && int(t.stack[height].Evaluation()) < ValueWin {
		return int(t.stack[height].Evaluation())
	}

	// Null move pruning
	if pos.LastMove != NullMove && depth >= 2 && !inCheck && (!hashOk || (hashFlag&TransAlpha == 0) || int(hashValue) >= beta) && !IsLateEndGame(pos) && int(t.stack[height].Evaluation()) >= beta {
		pos.MakeNullMove(child)
		t.MoveEvaluator.CurrentMove[height] = NullMove
		reduction := depth/4 + 3 + Min(int(t.stack[height].Evaluation())-beta, 384)/128
		tmpVal = -t.alphaBeta(depth-reduction, -beta, -beta+1, height+1, child.IsInCheck())
		if tmpVal >= beta {
			return beta
		}
	}

	var evaled []EvaledMove

	// Probcut pruning
	// If we have a good enough capture and a reduced search returns a value
	// much above beta, we can (almost) safely prune the previous move.
	if !pvNode && depth >= probCutDepth && Abs(beta) < ValueWin {
		evaluation := int(t.stack[height].Evaluation())
		improving := BoolToInt(height < 2 || evaluation > int(t.stack[height-2].Evaluation()))
		rBeta := Min(beta+probCutMargin-45*improving, ValueWin-1)
		//Idea from stockfish
		if !(hashMove != NullMove && int(hashDepth) >= depth-4 && int(hashValue) < rBeta) {
			evaled = pos.GenerateAllCaptures(t.stack[height].moves[:])
			t.EvaluateQsMoves(pos, evaled, hashMove, false)
			probCutCount := 0
			for i := 0; i < len(evaled) && probCutCount < 3; i++ {
				maxMoveToFirst(evaled[i:])
				if !SeeAbove(pos, evaled[i].Move, rBeta-evaluation) {
					continue
				}
				if !pos.MakeMove(evaled[i].Move, child) {
					continue
				}
				probCutCount++
				t.SetCurrentMove(height, evaled[i].Move)
				isChildInCheck := child.IsInCheck()
				tmpVal = -t.quiescence(0, -rBeta, -rBeta+1, height+1, isChildInCheck)
				if tmpVal >= rBeta {
					tmpVal = -t.alphaBeta(depth-4, -rBeta, -rBeta+1, height+1, isChildInCheck)
				}
				if tmpVal >= rBeta {
					return tmpVal
				}
			}
		}
	}

	val := MinInt

	// Internal iterative deepening
	// https://www.chessprogramming.org/Internal_Iterative_Deepening
	// Values taken from Laser
	if hashMove == NullMove && !inCheck && ((pvNode && depth >= 6) || (!pvNode && depth >= 8)) {
		var iiDepth int
		if pvNode {
			iiDepth = depth - depth/4 - 1
		} else {
			iiDepth = (depth - 5) / 2
		}
		t.alphaBeta(iiDepth, alpha, beta, height, inCheck)
		_, _, _, hashMove, _ = transposition.GlobalTransTable.Get(pos.Key)
	}

	// Quiet moves are stored in order to reduce their history value at the end of search
	quietsSearched := t.stack[height].quietsSearched[:0]
	bestMove := NullMove
	moveCount := 0
	movesSorted := false
	hashMoveChecked := false
	seeMargins := [2]int{seeQuietMargin * depth, seeNoisyMargin * depth * depth}

	// Check hashMove before move generation
	if pos.IsMovePseudoLegal(hashMove) {
		hashMoveChecked = true
		if pos.MakeMove(hashMove, child) {
			moveCount++
			childInCheck := child.IsInCheck()
			newDepth := depth - 1

			// Prefetch as early as possible
			transposition.GlobalTransTable.Prefetch(child.Key)

			// Castling extension
			if hashMove.IsCastling() {
				newDepth++
				// Check extension
				// Moves with positive SEE and gives check are searched with increased depth
			} else if inCheck && SeeSign(pos, hashMove) {
				newDepth++
				// Singular extension
				// https://www.chessprogramming.org/Singular_Extensions
			} else if depth >= 8 && int(hashDepth) >= depth-2 && hashFlag != TransAlpha {
				evaled = pos.GenerateAllMoves(t.stack[height].moves[:])
				t.EvaluateMoves(pos, evaled, hashMove, height, depth)
				sortMoves(evaled)
				movesSorted = true
				evaled = evaled[1:]
				if t.isMoveSingular(depth, height, hashMove, int(hashValue), evaled) {
					newDepth++
				}
			}

			// Store move if it is quiet
			if !hashMove.IsCaptureOrPromotion() {
				quietsSearched = append(quietsSearched, hashMove)
			}

			t.SetCurrentMove(height, hashMove)

			tmpVal = -t.alphaBeta(newDepth, -beta, -alpha, height+1, childInCheck)

			if tmpVal > val {
				val = tmpVal
				if val > alpha {
					alpha = val
					bestMove = hashMove
					if alpha >= beta {
						goto afterLoop
					}
					t.stack[height].PV.assign(hashMove, &t.stack[height+1].PV)
				}
			}
		}
	}

	// Generate moves if not generated in hashMove Check
	if !movesSorted {
		evaled = pos.GenerateAllMoves(t.stack[height].moves[:])
		if hashMoveChecked {
			moveToFirst(evaled, hashMove)
			evaled = evaled[1:] // Ignore hash move
		}
		t.EvaluateMoves(pos, evaled, hashMove, height, depth)
	}

	for i := range evaled {
		// Move might have been already sorted if singularity have been checked
		if !movesSorted {
			// Sort first few moves with selection sort
			if i < 3 || len(evaled)-i < 3 {
				maxIdx := i
				for idx := i + 1; idx < len(evaled); idx++ {
					if evaled[idx].Value > evaled[maxIdx].Value {
						maxIdx = idx
					}
				}
				evaled[i], evaled[maxIdx] = evaled[maxIdx], evaled[i]
			} else {
				// Sort rest of moves with shell sort
				sortMoves(evaled[i:])
				movesSorted = true
			}
		}
		isNoisy := evaled[i].Move.IsCaptureOrPromotion()

		if val > ValueLoss && !inCheck && moveCount > 0 && evaled[i].Value < MinSpecialMoveValue && !isNoisy {
			if depth <= futilityPruningDepth && int(t.stack[height].Evaluation())+int(PawnValueMiddle)*depth <= alpha {
				continue
			}
			if depth <= moveCountPruningDepth && moveCount >= moveCountPruning(BoolToInt(height <= 2 || t.stack[height].Evaluation() >= t.stack[height-2].Evaluation()), depth) {
				continue
			}
		}

		if val > ValueLoss &&
			depth <= seePruningDepth &&
			moveCount > 0 &&
			evaled[i].Value < MinGoodCapture &&
			!SeeAbove(pos, evaled[i].Move, seeMargins[BoolToInt(isNoisy)]) {
			continue
		}

		if !pos.MakeMove(evaled[i].Move, child) {
			continue
		}

		t.SetCurrentMove(height, evaled[i].Move)

		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		moveCount++
		childInCheck := child.IsInCheck()
		reduction := 0

		// Late Move Reduction
		// https://www.chessprogramming.org/Late_Move_Reductions
		if depth >= 3 && !inCheck && moveCount > 1 && !isNoisy && !childInCheck {
			reduction = lmr(depth, moveCount)
			reduction += BoolToInt(!pvNode)

			// less reudction for special moves
			reduction -= BoolToInt(evaled[i].Value >= MinSpecialMoveValue)

			// Increase reduction if not improving
			reduction += BoolToInt(height <= 2 || t.stack[height].Evaluation() < t.stack[height-2].Evaluation())
			reduction = Max(0, Min(depth-2, reduction))
		}
		newDepth := depth - 1

		// Castling extension
		if evaled[i].IsCastling() {
			newDepth++

			// Check extension
			// Moves with positive SEE and gives check are searched with increased depth
		} else if inCheck && SeeSign(pos, evaled[i].Move) {
			newDepth++
		}

		// Store move if it is quiet
		if !isNoisy {
			quietsSearched = append(quietsSearched, evaled[i].Move)
		}

		// Search conditions as in Ethereal
		// Search with null window and reduced depth if lmr
		if reduction > 0 {
			tmpVal = -t.alphaBeta(newDepth-reduction, -(alpha + 1), -alpha, height+1, childInCheck)
		}
		// Search with null window without reduced depth if
		// search with lmr null window exceeded alpha or
		// not in pv (this is the same as normal search as non pv nodes are searched with null window anyway)
		// pv and not first move
		if (reduction > 0 && tmpVal > alpha) || (reduction == 0 && !(pvNode && moveCount == 1)) {
			tmpVal = -t.alphaBeta(newDepth, -(alpha + 1), -alpha, height+1, childInCheck)
		}
		// If pvNode and first move or search with null window exceeded alpha, search with full window
		if pvNode && (moveCount == 1 || tmpVal > alpha) {
			tmpVal = -t.alphaBeta(newDepth, -beta, -alpha, height+1, childInCheck)
		}

		if tmpVal > val {
			val = tmpVal
			if val > alpha {
				alpha = val
				bestMove = evaled[i].Move
				if alpha >= beta {
					break
				}
				t.stack[height].PV.assign(evaled[i].Move, &t.stack[height+1].PV)
			}
		}
	}

	if moveCount == 0 {
		if inCheck {
			return lossIn(height)
		}
		return t.contempt(pos, depth)
	}

afterLoop:
	if alpha >= beta && bestMove != NullMove && !bestMove.IsCaptureOrPromotion() {
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
	transposition.GlobalTransTable.Set(pos.Key, transposition.ValueToTrans(alpha, height), depth, bestMove, flag)
	return alpha
}

func (t *thread) isMoveSingular(depth, height int, hashMove Move, hashValue int, moves []EvaledMove) bool {
	var pos *Position = &t.stack[height].position
	var child *Position = &t.stack[height+1].position
	// Store child as we already made a move into it in alphaBeta
	oldChild := *child
	val := -Mate
	rBeta := Max(hashValue-depth, -Mate)
	quiets := 0
	for i := range moves {
		if !pos.MakeMove(moves[i].Move, child) {
			continue
		}
		t.SetCurrentMove(height, moves[i].Move)
		val = -t.alphaBeta(depth/2-1, -rBeta-1, -rBeta, height+1, child.IsInCheck())
		if val > rBeta {
			break
		}
		if !moves[i].Move.IsCaptureOrPromotion() {
			quiets++
			if quiets >= 6 {
				break
			}
		} else if moves[i].Value < MaxBadCapture {
			break
		}
	}
	// restore child
	*child = oldChild
	return val <= rBeta
}

func (t *thread) isDraw(height int) bool {
	var pos *Position = &t.stack[height].position

	// Fifty move rule
	if pos.FiftyMove > 100 {
		return true
	}

	// Cannot mate with only one minor piece and no pawns
	if (pos.Pieces[Pawn]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 &&
		!MoreThanOne(pos.Pieces[Knight]|pos.Pieces[Bishop]) {
		return true
	}

	// Look for repetitoin in current search stack
	for i := height - 1; i >= 0; i-- {
		descendant := &t.stack[i].position
		if descendant.Key == pos.Key {
			return true
		}
		if descendant.FiftyMove == 0 || descendant.LastMove == NullMove {
			return false
		}
	}

	// Check for repetition in already played positions
	if _, found := t.engine.RepeatedPositions[pos.Key]; found {
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

// https://www.chessprogramming.org/Aspiration_Windows
// After a lot of tries ELO gain have been accomplished only with relatively large window(50 cp)
func (t *thread) aspirationWindow(depth, lastValue int, moves []EvaledMove, resultChan chan result) int {
	var alpha, beta int
	delta := WindowSize
	searchDepth := depth
	if depth >= WindowDepth {
		alpha = Max(-Mate, lastValue-delta)
		beta = Min(Mate, lastValue+delta)
	} else {
		// Search with [-Mate, Mate] in shallow depths
		alpha = -Mate
		beta = Mate
	}
	for {
		res := t.depSearch(Max(1, searchDepth), alpha, beta, moves)
		if res.value > alpha && res.value < beta {
			resultChan <- res
			return res.value
		}
		if res.value <= alpha {
			beta = (alpha + beta) / 2
			alpha = Max(-Mate, alpha-delta)
			searchDepth = depth
		}
		if res.value >= beta {
			beta = Min(Mate, beta+delta)
			searchDepth--
		}
		delta += delta/2 + 5
	}
}

// depSearch is special case of alphaBeta function for root node
func (t *thread) depSearch(depth, alpha, beta int, moves []EvaledMove) result {
	var pos *Position = &t.stack[0].position
	var child *Position = &t.stack[1].position
	var bestMove Move = NullMove
	inCheck := pos.IsInCheck()
	moveCount := 0
	t.stack[0].PV.clear()
	t.stack[0].InvalidateEvaluation()
	t.MoveEvaluator.ResetKillers(1)
	quietsSearched := t.stack[0].quietsSearched[:0]

	for i := range moves {
		pos.MakeLegalMove(moves[i].Move, child)
		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		t.SetCurrentMove(0, moves[i].Move)

		moveCount++
		if !moves[i].IsCaptureOrPromotion() {
			quietsSearched = append(quietsSearched, moves[i].Move)
		}
		reduction := 0
		childInCheck := child.IsInCheck()
		if !inCheck && moveCount > 1 && moves[i].Value < MinSpecialMoveValue && !moves[i].Move.IsCaptureOrPromotion() &&
			!childInCheck {
			if depth <= moveCountPruningDepth && moveCount >= moveCountPruning(1, depth) {
				continue
			}
			if depth >= 3 {
				reduction = lmr(depth, moveCount) - 1
				reduction = Max(0, Min(depth-2, reduction))
			}
		}
		var val int
		newDepth := depth - 1
		if moves[i].IsCastling() {
			newDepth++
		} else if inCheck && SeeSign(pos, moves[i].Move) {
			newDepth++
		}
		if reduction > 0 {
			val = -t.alphaBeta(newDepth-reduction, -(alpha + 1), -alpha, 1, childInCheck)
			if val <= alpha {
				continue
			}
		}
		val = -t.alphaBeta(newDepth, -beta, -alpha, 1, childInCheck)
		if val > alpha {
			alpha = val
			bestMove = moves[i].Move
			if alpha >= beta {
				break
			}
			t.stack[0].PV.assign(moves[i].Move, &t.stack[1].PV)
		}
	}
	if moveCount == 0 {
		if inCheck {
			alpha = lossIn(0)
		} else {
			alpha = t.contempt(pos, depth)
		}
	}
	if alpha >= beta && bestMove != NullMove && !bestMove.IsCaptureOrPromotion() {
		t.Update(pos, quietsSearched, bestMove, depth, 0)
	}
	t.EvaluateMoves(pos, moves, bestMove, 0, depth)
	sortMoves(moves)
	return result{bestMove, alpha, depth, cloneMoves(t.stack[0].PV.items[:t.stack[0].PV.size])}
}

func (e *Engine) singleThreadBestMove(ctx context.Context, rootMoves []EvaledMove) Move {
	var lastBestMove Move
	thread := e.threads[0]
	lastValue := -Mate
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		go func(depth int) {
			defer recoverFromTimeout()
			lastValue = thread.aspirationWindow(depth, lastValue, rootMoves, resultChan)
		}(i)
		select {
		case <-ctx.Done():
			return lastBestMove
		case res := <-resultChan:
			e.Update(SearchInfo{newUciScore(res.value), i, thread.nodes, res.moves})
			if res.value >= ValueWin && depthToMate(res.value) <= i {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			if i >= MAX_HEIGHT {
				return res.Move
			}
			e.updateTime(res.depth, res.value)
			if e.isSoftTimeout(i, thread.nodes) {
				return res.Move
			}
			lastBestMove = res.Move
		}
	}
}

func (t *thread) iterativeDeepening(moves []EvaledMove, resultChan chan result, idx int) {
	mainThread := idx == 0
	lastValue := -Mate
	// I do not think this matters much, but at the beginning only thread with id 0 have sorted moves list
	if !mainThread {
		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
	}
	// Depth skipping pattern taken from Ethereal
	cycle := idx % SMPCycles
	for depth := 1; depth <= MAX_HEIGHT; depth++ {
		lastValue = t.aspirationWindow(depth, lastValue, moves, resultChan)
		if !mainThread && (depth+cycle)%SkipDepths[cycle] == 0 {
			depth += SkipSize[cycle]
		}
	}
}

func (e *Engine) bestMove(ctx context.Context, pos *Position) Move {
	for i := range e.threads {
		e.threads[i].stack[0].position = *pos
		e.threads[i].nodes = 0
	}

	rootMoves := pos.GenerateAllLegalMoves()
	if len(rootMoves) == 1 {
		return rootMoves[0].Move
	}

	if fathom.IsDTZProbeable(pos) {
		if ok, bestMove, wdl, dtz := fathom.ProbeDTZ(pos, rootMoves); ok {
			var score int
			if wdl == fathom.TB_LOSS {
				score = ValueLoss + dtz + 1
			} else if wdl == fathom.TB_WIN {
				score = ValueWin - dtz - 1
			} else {
				score = 0
			}
			e.Update(SearchInfo{newUciScore(score), MAX_HEIGHT - 1, 0, []Move{bestMove}})
			return bestMove
		}
	}

	ordMove := NullMove
	if hashOk, _, _, hashMove, _ := transposition.GlobalTransTable.Get(pos.Key); hashOk {
		ordMove = hashMove
	}
	e.threads[0].EvaluateMoves(pos, rootMoves, ordMove, 0, 127)

	sortMoves(rootMoves)

	if e.Threads.Val == 1 {
		return e.singleThreadBestMove(ctx, rootMoves)
	}

	var wg = &sync.WaitGroup{}
	resultChan := make(chan result)
	for i := range e.threads {
		wg.Add(1)
		// Start parallel searching
		go func(idx int) {
			defer recoverFromTimeout()
			e.threads[idx].iterativeDeepening(cloneEvaledMoves(rootMoves), resultChan, idx)
			wg.Done()
		}(i)
	}

	// Wait for closing
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	prevDepth := 0
	var lastBestMove Move
	for {
		select {
		case <-e.done:
			// Hard timeout
			return lastBestMove
		case res := <-resultChan:
			// If thread reports result for depth that is lower than already calculated one, ignore results
			if res.depth <= prevDepth {
				continue
			}
			nodes := e.nodes()
			e.Update(SearchInfo{newUciScore(res.value), res.depth, nodes, res.moves})
			if res.value >= ValueWin && depthToMate(res.value) <= res.depth {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			if res.depth >= MAX_HEIGHT {
				return res.Move
			}
			e.updateTime(res.depth, res.value)
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
