package engine

import (
	"context"
	"math"
	"math/rand"
	"sync"

	. "github.com/mhib/combusken/chess"
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
const ValueTbWinInMaxDepth = ValueWin - MaxHeight - 1

const seePruningDepth = 10
const seeQuietMargin = -100
const seeNoisyMargin = -28

const reverseFutilityPruningDepth = 7
const reverseFutilityPruningMargin = 100

const moveCountPruningDepth = 8
const futilityPruningDepth = 8

const counterMovePruningDepth = 2
const counterMovePruningVal = -785
const secondCounterMovePruningDepth = 8
const secondCounterMovePruningVal = -7340

const probCutDepth = 6
const probCutMargin = 80

const WindowSize = 18
const WindowDepth = 5
const WindowDeltaInc = 10

const qSDepthChecks = 0
const qSDepthNoChecks = -1

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

func (t *thread) quiescence(depth, alpha, beta, height int) int {
	t.incNodes()
	t.stack[height].PV.clear()
	t.seldepth = Max(t.seldepth, height)

	pos := &t.stack[height].position
	alphaOrig := alpha
	pvNode := alpha != beta-1
	inCheck := pos.IsInCheck()

	if height >= MaxHeight || t.isDraw(height) {
		return t.contempt(pos, depth)
	}

	var ttDepth int
	if inCheck || depth >= qSDepthChecks {
		ttDepth = qSDepthChecks
	} else {
		ttDepth = qSDepthNoChecks
	}
	hashOk, hashValue, hashEval, hashDepth, hashMove, hashFlag, _ := transposition.GlobalTransTable.Get(pos.Key)
	if hashOk && hashValue != UnknownValue && int(hashDepth) >= ttDepth {
		hashValue = transposition.ValueFromTrans(hashValue, height)
		if hashFlag == TransExact || (hashFlag == TransAlpha && int(hashValue) <= alpha) ||
			(hashFlag == TransBeta && int(hashValue) >= beta) {
			return int(hashValue)
		}
	}

	child := &t.stack[height+1].position

	bestMove := NullMove

	moveCount := 0

	var bestVal int
	var eval int16

	if inCheck {
		bestVal = MinInt
		t.stack[height].InitNormal(pos, &t.MoveHistory, height, hashMove)
		eval = UnknownValue
	} else {
		if hashOk && hashEval != UnknownValue {
			eval = hashEval
			bestVal = int(eval)
			if hashValue != UnknownValue {
				requiredFlag := uint8(TransAlpha)
				if int(hashValue) > bestVal {
					requiredFlag = uint8(TransBeta)
				}
				if requiredFlag&hashFlag != 0 {
					bestVal = int(hashValue)
				}
			}
		} else {
			if pos.LastMove != NullMove {
				eval = int16(t.Evaluate(pos))
			} else {
				eval = -t.getEvaluation(height-1) + 2*Tempo
			}
			bestVal = int(eval)
			transposition.GlobalTransTable.Set(pos.Key, UnknownValue, eval, transposition.NoneDepth, NullMove, TransNone, pvNode)
		}
		// Early return if not in check and evaluation exceeded beta
		if bestVal >= beta {
			return beta
		}
		if alpha < bestVal {
			alpha = bestVal
		}
		t.stack[height].InitQs()
	}

	for {
		move := t.getNextMove(pos, 0, height)
		if move == NullMove {
			break
		}
		if !pos.MakeMove(move, child) {
			continue
		}

		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		t.SetCurrentMove(height, move)
		moveCount++

		t.ApplyMove(move, pos, child)

		val := -t.quiescence(depth-1, -beta, -alpha, height+1)

		t.RevertMove(move, pos, child)

		if val > bestVal {
			bestVal = val
			bestMove = move
			if val > alpha {
				alpha = val
				if pvNode {
					t.stack[height].PV.assign(move, &t.stack[height+1].PV)
				}
				if val >= beta {
					break
				}
			}

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

	transposition.GlobalTransTable.Set(pos.Key, transposition.ValueToTrans(alpha, height), eval, ttDepth, bestMove, flag, pvNode)

	return alpha
}

// Currently draws are scored as 0 +/- 1 randomly
func (t *thread) contempt(pos *Position, depth int) int {
	if depth < 4 {
		return 0
	}
	return 2*(t.nodes&1) - 1
}

func moveCountPruning(improving, depth int) int {
	return (2+depth*depth)*(1+improving)/2 - 1
}

func (t *thread) alphaBeta(depth, alpha, beta, height int, cutNode bool) int {
	t.incNodes()
	t.stack[height].PV.clear()
	t.seldepth = Max(t.seldepth, height)

	var pos *Position = &t.stack[height].position

	if height >= MaxHeight || t.isDraw(height) {
		return t.contempt(pos, depth)
	}

	// Node is not pv if it is searched with null window
	pvNode := alpha != beta-1
	inCheck := pos.IsInCheck()

	// Mate distance pruning
	alpha = Max(lossIn(height+1), alpha)
	beta = Min(winIn(height+2), beta)
	if alpha >= beta {
		return alpha
	}

	alphaOrig := alpha
	hashOk, hashValue, hashEval, hashDepth, hashMove, hashFlag, hashPv := transposition.GlobalTransTable.Get(pos.Key)
	hashPv = hashPv || pvNode
	var val int
	if hashOk && hashValue != UnknownValue {
		hashValue = transposition.ValueFromTrans(hashValue, height)
		val := int(hashValue)
		// Hash pruning
		if hashDepth >= int16(depth) && (depth == 0 || !pvNode) {
			if hashFlag == TransExact {
				return val
			}
			if hashFlag == TransAlpha && val <= alpha {
				return alpha
			}
			if hashFlag == TransBeta && val >= beta {
				return beta
			}
		}
	}

	// Probe tablebase
	if fathom.IsWDLProbeable(pos, depth) {
		if tbResult := fathom.ProbeWDL(pos, depth); tbResult != fathom.TbResultFailed {
			t.tbhits++
			var ttBound int
			if tbResult == fathom.TbLoss {
				val = ValueLoss + height + 1
				ttBound = TransAlpha
			} else if tbResult == fathom.TbWin {
				val = ValueWin - height - 1
				ttBound = TransBeta
			} else {
				val = 0
				ttBound = TransExact
			}
			if ttBound == TransExact || ttBound == TransBeta && val >= beta || ttBound == TransAlpha && val <= alpha {
				transposition.GlobalTransTable.Set(pos.Key, int16(val), UnknownValue, MaxHeight, NullMove, ttBound, pvNode)
				return val
			}
		}
	}

	if depth <= 0 {
		return t.quiescence(0, alpha, beta, height)
	}

	var child *Position = &t.stack[height+1].position

	t.ResetKillers(height + 1)
	var improving bool
	var eval int16
	if inCheck {
		improving = false
		eval = UnknownValue
		t.setEvaluation(height, eval)
		goto afterPreMovesPruning
	}
	if hashOk && hashEval != UnknownValue {
		eval = hashEval
		t.setEvaluation(height, hashEval)
		// Idea from stockfish
		// Use hashValue as better position evaluation
		if hashValue != UnknownValue {
			requiredFlag := TransAlpha
			if hashValue > eval {
				requiredFlag = TransBeta
			}
			if requiredFlag&int(hashFlag) != 0 {
				eval = hashValue
			}
		}
	} else {
		if pos.LastMove != NullMove {
			eval = int16(t.Evaluate(pos))
		} else {
			eval = -t.getEvaluation(height-1) + 2*Tempo
		}
		t.setEvaluation(height, eval)
		transposition.GlobalTransTable.Set(pos.Key, UnknownValue, eval, transposition.NoneDepth, NullMove, TransNone, false)
	}

	if height > 1 {
		if t.getEvaluation(height-2) != UnknownValue {
			improving = t.getEvaluation(height) > t.getEvaluation(height-2)
		} else {
			improving = height <= 3 ||
				t.getEvaluation(height-4) == UnknownValue ||
				t.getEvaluation(height) > t.getEvaluation(height-4)
		}
	} else {
		improving = true
	}

	// Reverse futility pruning
	if !pvNode && depth < reverseFutilityPruningDepth && int(eval)-reverseFutilityPruningMargin*(depth-BoolToInt(improving)) >= beta && int(eval) < ValueTbWinInMaxDepth {
		return int(eval)
	}

	// Null move pruning
	if !pvNode && pos.LastMove != NullMove && t.disableNmpColor != pos.SideToMove && depth >= 2 && t.GetPreviousMoveFromCurrentSide(height) != NullMove && (!hashOk || (hashFlag&TransAlpha == 0) || int(hashValue) >= beta) && !isPiecelessEndGame(pos) && int(eval) >= beta {
		pos.MakeNullMove(child)
		t.CurrentMove[height] = NullMove
		reduction := depth/4 + 3 + Min(int(eval)-beta, 384)/128
		val = -t.alphaBeta(depth-reduction, -beta, -beta+1, height+1, !cutNode)
		if val >= beta {
			if depth < 10 || t.disableNmpColor != ColourNone {
				if val >= ValueTbWinInMaxDepth {
					return beta
				} else {
					return val
				}
			} else {
				// Null move pruning verification search.
				// Idea from stockfish
				t.disableNmpColor = pos.SideToMove
				val = t.alphaBeta(depth-reduction, beta-1, beta, height, false)
				t.disableNmpColor = ColourNone
				if val >= beta {
					if val >= ValueTbWinInMaxDepth {
						return beta
					} else {
						return val
					}
				}
			}
		}
	}

	// Probcut pruning
	// If we have a good enough capture and a reduced search returns a value
	// much above beta, we can (almost) safely prune the previous move.
	if !pvNode && depth >= probCutDepth && Abs(beta) < ValueTbWinInMaxDepth {
		rBeta := Min(beta+probCutMargin, ValueTbWinInMaxDepth-1)
		//Idea from stockfish
		if !(hashMove != NullMove && int(hashDepth) >= depth-4 && int(hashValue) < rBeta) {
			t.stack[height].InitQs()
			probCutCount := 0
			for probCutCount < 3 {

				move := t.getNextMove(pos, depth, height)
				if move == NullMove {
					break
				}
				if !pos.MakeMove(move, child) {
					continue
				}

				probCutCount++
				t.SetCurrentMove(height, move)

				t.ApplyMove(move, pos, child)

				val = -t.quiescence(0, -rBeta, -rBeta+1, height+1)
				if val >= rBeta {
					val = -t.alphaBeta(depth-4, -rBeta, -rBeta+1, height+1, !cutNode)
				}

				t.RevertMove(move, pos, child)

				if val >= rBeta {
					return val
				}
			}
		}
	}

afterPreMovesPruning:
	bestVal := MinInt

	// IID alternative by Ed Schroeder
	if pvNode && depth >= 3 && hashMove == NullMove {
		depth--
	}

	noisySearched := t.stack[height].noisySearched[:0]

	// Quiet moves are stored in order to reduce their history value at the end of search
	quietsSearched := t.stack[height].quietsSearched[:0]
	bestMove := NullMove
	moveCount := 0
	t.stack[height].InitNormal(pos, &t.MoveHistory, height, hashMove)

	for {
		move := t.getNextMove(pos, depth, height)
		if move == NullMove {
			break
		}
		isNoisy := move.IsCaptureOrPromotion()

		if bestVal > ValueLoss && !inCheck && moveCount > 0 && t.stack[height].GetStage() > StageGenerateQuiet && !isNoisy {
			if depth <= futilityPruningDepth && int(eval)+int(PawnValueMiddle)*depth <= alpha {
				t.skipQuiets(height)
				continue
			}
			if depth <= moveCountPruningDepth && moveCount >= moveCountPruning(BoolToInt(improving), depth) {
				t.skipQuiets(height)
				continue
			}
			if pos.LastMove != NullMove {
				if depth <= counterMovePruningDepth {
					if t.CounterHistoryValue(pos.LastMove, move) < counterMovePruningVal {
						continue
					}
				} else if depth <= secondCounterMovePruningDepth && t.CounterHistoryValue(pos.LastMove, move) < int32(depth*secondCounterMovePruningVal) {
					continue
				}
			}
		}

		if !pos.MakeMove(move, child) {
			continue
		}

		t.SetCurrentMove(height, move)

		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		moveCount++

		reduction := 0
		// Late Move Reduction
		// https://www.chessprogramming.org/Late_Move_Reductions
		if moveCount > 1 && (!isNoisy || cutNode) {
			reduction = lmr(depth, moveCount)

			// less reduction for special moves
			reduction -= BoolToInt(t.stack[height].GetStage() < StageGenerateQuiet)

			reduction += BoolToInt(!hashPv)

			reduction -= BoolToInt(child.IsInCheck())
			if isNoisy {
				reduction += BoolToInt(t.getCaptureHistory(move) < 0)
			} else {
				reduction += BoolToInt(cutNode) * 2
				// Increase reduction if not improving
				reduction += BoolToInt(!improving)

				reduction -= (t.HistoryValue(pos, move, t.GetPreviousMoveFromCurrentSide(height)) - 2746) / 12124
			}
			reduction = Max(0, Min(depth-2, reduction))
		} else if isNoisy && !child.IsInCheck() && moveCount > 1 && depth >= 2 {
			reduction = BoolToInt(t.getCaptureHistory(move) < 0)
		}

		if bestVal > ValueLoss && depth <= seePruningDepth && t.stack[height].GetStage() > StageGoodNoisy {
			reducedDepth := depth - reduction
			if (isNoisy && !SeeAbove(pos, move, seeNoisyMargin*reducedDepth*reducedDepth)) ||
				(!isNoisy && !SeeAbove(pos, move, seeQuietMargin*reducedDepth)) {
				continue
			}
		}

		var extension int
		if move == hashMove && depth >= 8 && int(hashDepth) >= depth-2 && hashFlag != TransAlpha {
			singularValue, singularBeta := t.singularSearch(depth, height, hashMove, int(hashValue), cutNode)
			if singularValue <= singularBeta {
				extension = 1
			} else if singularBeta >= beta {
				// Multi-cut pruning
				// Idea from stockfish
				return singularBeta
			}
		} else {
			extension = BoolToInt(move.IsCastling() || (inCheck && SeeSign(pos, move)))
		}

		newDepth := depth - 1 + extension

		// Store move if it is quiet
		if isNoisy {
			noisySearched = append(noisySearched, move)
		} else {
			quietsSearched = append(quietsSearched, move)
		}

		t.ApplyMove(move, pos, child)

		// Search conditions as in Ethereal
		// Search with null window and reduced depth if lmr
		if reduction > 0 {
			val = -t.alphaBeta(newDepth-reduction, -(alpha + 1), -alpha, height+1, true)
		}
		// Search with null window without reduced depth if
		// search with lmr null window exceeded alpha or
		// not in pv (this is the same as normal search as non pv nodes are searched with null window anyway)
		// pv and not first move
		if (reduction > 0 && val > alpha) || (reduction == 0 && !(pvNode && moveCount == 1)) {
			val = -t.alphaBeta(newDepth, -(alpha + 1), -alpha, height+1, !cutNode)
		}
		// If pvNode and first move or search with null window exceeded alpha, search with full window
		if pvNode && (moveCount == 1 || val > alpha) {
			val = -t.alphaBeta(newDepth, -beta, -alpha, height+1, false)
		}
		t.RevertMove(move, pos, child)

		if val > bestVal {
			bestVal = val
			bestMove = move
			if val > alpha {
				alpha = val
				if pvNode {
					t.stack[height].PV.assign(move, &t.stack[height+1].PV)
				}
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
		return t.contempt(pos, depth)
	}

	if alpha >= beta && bestMove != NullMove {
		t.UpdateNoisy(pos, noisySearched, bestMove, depth)

		if !bestMove.IsCaptureOrPromotion() {
			t.UpdateQuiet(pos, quietsSearched, bestMove, depth, height)
		}
	}

	var flag int
	if alpha == alphaOrig {
		flag = TransAlpha
	} else if alpha >= beta {
		flag = TransBeta
	} else {
		flag = TransExact
	}
	transposition.GlobalTransTable.Set(pos.Key, transposition.ValueToTrans(alpha, height), t.getEvaluation(height), depth, bestMove, flag, pvNode)
	return alpha
}

func (t *thread) singularSearch(depth, height int, hashMove Move, hashValue int, cutNode bool) (val, rBeta int) {
	var pos *Position = &t.stack[height].position
	var child *Position = &t.stack[height+1].position
	// Store child as we already made a move into it in alphaBeta
	oldChild := *child
	val = -Mate
	rBeta = Max(hashValue-depth, -Mate)
	quiets := 0
	t.stack[height].InitSingular()
	for {
		move := t.getNextMove(pos, depth, height)
		if move == NullMove || t.stack[height].GetStage() >= StageNoisy {
			break
		}
		if !pos.MakeMove(move, child) {
			continue
		}
		t.SetCurrentMove(height, move)

		t.ApplyMove(move, pos, child)

		val = -t.alphaBeta(depth/2-1, -rBeta-1, -rBeta, height+1, cutNode)

		t.RevertMove(move, pos, child)

		if val > rBeta {
			break
		}
		if !move.IsCaptureOrPromotion() {
			quiets++
			if quiets >= 4 {
				break
			}
		}
	}
	// restore child
	*child = oldChild
	t.stack[height].RestoreFromSingular()
	return val, rBeta
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

	// Look for repetition in the current search stack
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

type searchResult struct {
	value    int
	depth    int
	seldepth int
	moves    []Move
}

type aspirationWindowResult struct {
	*searchResult
	requestedDepth int
	multiPV        int
}

func (t *thread) aspirationWindow(depth, lastValue int, moves []EvaledMove, multiPV int) aspirationWindowResult {
	var alpha, beta int
	delta := WindowSize
	searchDepth := depth
	if depth >= WindowDepth {
		alpha = Max(-Mate, lastValue-delta)
		beta = Min(Mate, lastValue+delta)

		tempo := lastValue / 4
		tempo = Min(30, Max(tempo, -30))
		if t.stack[0].position.SideToMove == White {
			t.SetContempt(S(int16(tempo), int16(tempo/2)))
		} else {
			t.SetContempt(-S(int16(tempo), int16(tempo/2)))
		}
	} else {
		// Search with [-Mate, Mate] in shallow depths
		alpha = -Mate
		beta = Mate
		t.SetContempt(Score(0))
	}
	for {
		res := t.depSearch(Max(1, searchDepth), alpha, beta, moves)
		if res.value > alpha && res.value < beta {
			return aspirationWindowResult{&res, depth, multiPV}
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
		delta += delta/2 + WindowDeltaInc
	}
}

// depSearch is special case of alphaBeta function for root node
func (t *thread) depSearch(depth, alpha, beta int, moves []EvaledMove) searchResult {
	t.seldepth = 0

	var pos *Position = &t.stack[0].position
	var child *Position = &t.stack[1].position
	var bestMove Move = NullMove
	alphaOrig := alpha
	inCheck := pos.IsInCheck()
	moveCount := 0
	eval := int16(t.Evaluate(pos))
	t.setEvaluation(0, eval)
	t.stack[0].PV.clear()
	t.ResetKillers(1)
	multiPV := t.isMainThread() && t.engine.MultiPV.Val != 1
	quietsSearched := t.stack[0].quietsSearched[:0]
	noisySearched := t.stack[0].noisySearched[:0]
	bestVal := MinInt
	var val int

	for i := range moves {
		if multiPV && t.engine.IsMoveExcluded(moves[i].Move) {
			continue
		}
		pos.MakeLegalMove(moves[i].Move, child)
		// Prefetch as early as possible
		transposition.GlobalTransTable.Prefetch(child.Key)

		t.SetCurrentMove(0, moves[i].Move)

		moveCount++
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
		newDepth := depth - 1
		if moves[i].IsCastling() {
			newDepth++
		} else if inCheck && SeeSign(pos, moves[i].Move) {
			newDepth++
		}

		if moves[i].IsCaptureOrPromotion() {
			noisySearched = append(noisySearched, moves[i].Move)
		} else {
			quietsSearched = append(quietsSearched, moves[i].Move)
		}

		t.ApplyMove(moves[i].Move, pos, child)

		if reduction > 0 {
			val = -t.alphaBeta(newDepth-reduction, -(alpha + 1), -alpha, 1, true)
			if val <= alpha {
				t.RevertMove(moves[i].Move, pos, child)
				continue
			}
		}
		val = -t.alphaBeta(newDepth, -beta, -alpha, 1, false)
		t.RevertMove(moves[i].Move, pos, child)
		if val > bestVal {
			bestVal = val
			bestMove = moves[i].Move
			if val > alpha {
				alpha = val
				if alpha >= beta {
					break
				}
				t.stack[0].PV.assign(moves[i].Move, &t.stack[1].PV)
			}
		}
	}
	if moveCount == 0 {
		if inCheck {
			alpha = lossIn(0)
		} else {
			alpha = t.contempt(pos, depth)
		}
	}
	if alpha >= beta && bestMove != NullMove {
		t.UpdateNoisy(pos, noisySearched, bestMove, depth)
		if !bestMove.IsCaptureOrPromotion() {
			t.UpdateQuiet(pos, quietsSearched, bestMove, depth, 0)
		}
	}
	t.EvaluateMoves(pos, moves, bestMove, 0, depth)
	sortMoves(moves)
	var flag int
	if alpha == alphaOrig {
		flag = TransAlpha
	} else if alpha >= beta {
		flag = TransBeta
	} else {
		flag = TransExact
	}
	transposition.GlobalTransTable.Set(pos.Key, transposition.ValueToTrans(alpha, 0), eval, depth, bestMove, flag, true)
	return searchResult{alpha, depth, t.seldepth, cloneMoves(t.stack[0].PV.Moves())}
}

func (t *thread) iterativeDeepening(moves []EvaledMove, resultChan chan aspirationWindowResult) {
	var res aspirationWindowResult
	lastValue := -Mate
	// I do not think this matters much, but at the beginning only thread with id 0 have sorted moves list
	if !t.isMainThread() {
		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
	}

	for depth := 1; depth <= MaxHeight; depth++ {
		res = t.aspirationWindow(depth, lastValue, moves, 1)
		select {
		case resultChan <- res:
			lastValue = res.value
		case <-t.engine.done:
			return
		}
	}
}

func (t *thread) multiPVIterativeDeepening(moves []EvaledMove, resultChan chan aspirationWindowResult) {
	var res aspirationWindowResult
	multiPV := Min(t.engine.MultiPV.Val, len(moves))
	for depth := 1; depth <= MaxHeight; depth++ {
		for moveIdx := 1; moveIdx <= multiPV; moveIdx++ {
			t.engine.multiPVExcluded[moveIdx-1] = NullMove
			res = t.aspirationWindow(depth, int(t.engine.lastValues[moveIdx-1]), moves, moveIdx)
			select {
			case resultChan <- res:
				t.engine.multiPVExcluded[moveIdx-1] = res.moves[0]
				t.engine.lastValues[moveIdx-1] = int16(res.value)
			case <-t.engine.done:
				return
			}
		}
	}
}

func (t *thread) silentIterativeDeepening(moves []EvaledMove) {
	lastValue := -Mate
	rand.Shuffle(len(moves), func(i, j int) {
		moves[i], moves[j] = moves[j], moves[i]
	})
	for depth := 1; depth <= MaxHeight; depth++ {
		lastValue = t.aspirationWindow(depth, lastValue, moves, 1).value
	}
}

func (e *Engine) bestMove(ctx, ponderCtx context.Context, wg *sync.WaitGroup, pos *Position) (Move, Move) {
	isMultiPV := e.MultiPV.Val != 1
	e.multiPVExcluded[0] = NullMove
	for i := range e.threads {
		e.threads[i].stack[0].position = *pos
		e.threads[i].nodes = 0
		e.threads[i].tbhits = 0
		e.threads[i].disableNmpColor = ColourNone
		e.threads[i].Initialize(pos)
	}

	rootMoves := GenerateAllLegalMoves(pos)

	if !isMultiPV && fathom.IsDTZProbeable(pos) {
		if ok, bestMove, wdl, dtz := fathom.ProbeDTZ(pos, rootMoves); ok {
			var score int
			if wdl == fathom.TbLoss {
				score = ValueLoss + dtz + 1
			} else if wdl == fathom.TbWin {
				score = ValueWin - dtz - 1
			} else {
				score = 0
			}
			e.Update(&SearchInfo{newReportScore(score), MaxHeight - 1, MaxHeight - 1, 1, 0, 1, 0, 1, []Move{bestMove}})
			for range e.threads {
				wg.Done()
			}
			return bestMove, NullMove
		}
	}

	ordMove := NullMove
	if hashOk, _, _, _, hashMove, _, _ := transposition.GlobalTransTable.Get(pos.Key); hashOk {
		ordMove = hashMove
	}
	e.threads[0].EvaluateMoves(pos, rootMoves, ordMove, 0, 127)

	sortMoves(rootMoves)

	resultChan := make(chan aspirationWindowResult)
	if isMultiPV {
		for i := range e.threads {
			if i == 0 {
				go func(idx int) {
					defer recoverFromTimeout(wg)
					e.threads[idx].multiPVIterativeDeepening(cloneEvaledMoves(e.threads[idx].getRootMovesBuffer(), rootMoves), resultChan)
				}(i)
				continue
			}
			go func(idx int) {
				defer recoverFromTimeout(wg)
				e.threads[idx].silentIterativeDeepening(cloneEvaledMoves(e.threads[idx].getRootMovesBuffer(), rootMoves))
			}(i)
		}
	} else {
		for i := range e.threads {
			go func(idx int) {
				defer recoverFromTimeout(wg)
				e.threads[idx].iterativeDeepening(cloneEvaledMoves(e.threads[idx].getRootMovesBuffer(), rootMoves), resultChan)
			}(i)
		}

	}

	prevDepth := 0
	var lastBestMove, lastPonderMove Move
	for {
		select {
		case <-e.done:
			// Hard timeout
			return lastBestMove, lastPonderMove
		case res := <-resultChan:
			// If thread reports result for depth that is lower than already calculated one, ignore results
			if res.requestedDepth <= prevDepth && !isMultiPV {
				continue
			}
			nodes, tbhits := e.aggregatesInfo()
			timeSinceStart := e.getElapsedTime()
			e.Update(&SearchInfo{newReportScore(res.value), res.requestedDepth, res.seldepth, res.multiPV, nodes, int(float64(nodes) / timeSinceStart.Seconds()), int(timeSinceStart.Milliseconds()), tbhits, res.moves})
			if res.multiPV != 1 {
				continue
			}
			e.updateTime(res.depth, res.value)
			prevDepth = res.requestedDepth
			lastBestMove = res.moves[0]
			lastPonderMove = NullMove
			if len(res.moves) > 1 {
				lastPonderMove = res.moves[1]
			}
			if res.depth >= MaxHeight {
				return lastBestMove, lastPonderMove
			}
			// Do not stop searching even when found a mate when in MultiPV or ponder
			if isMultiPV || isContextActive(ponderCtx) {
				continue
			}
			if res.value >= ValueWin && depthToMate(res.value) <= res.depth {
				return lastBestMove, lastPonderMove
			}
		}
	}
}

func isContextActive(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

func cloneMoves(src []Move) []Move {
	dst := make([]Move, len(src))
	copy(dst, src)
	return dst
}

func cloneEvaledMoves(dst, src []EvaledMove) []EvaledMove {
	copy(dst, src)
	return dst[:len(src)]
}

func recoverFromTimeout(wg *sync.WaitGroup) {
	wg.Done()
	err := recover()
	if err != nil && err != errTimeout {
		panic(err)
	}
}

func isPiecelessEndGame(pos *Position) bool {
	return ((pos.Pieces[Rook] | pos.Pieces[Queen] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[pos.SideToMove]) == 0
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
	return lmrTable[Min(d, 63)][Min(m, 63)]
}

var lmrTable [64][64]int

func init() {
	base := -0.7675250065028055
	movePower := 1.1798078422389596
	depthPower := 0.3744551389641595

	lmrFormula := func(d, m int) int {
		return Max(int(math.Round(math.Log(math.Pow(float64(d), depthPower))*math.Log(math.Pow(float64(m), movePower))+base)), 0)
	}
	for d := 1; d < 64; d++ {
		for m := 1; m < 64; m++ {
			lmrTable[d][m] = lmrFormula(d, m)
		}
	}
}
