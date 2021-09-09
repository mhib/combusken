package tuning

// Tuner is based on this paper by Andrew Grant:
// https://github.com/AndyGrant/Ethereal/blob/master/Tuning.pdf

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/evaluation"
	. "github.com/mhib/combusken/utils"
)

const (
	Middle = iota
	End
)

const learningRate = 0.1

type linearCoefficient struct {
	idx   int
	value int
}

type safetyCoefficient struct {
	idx        int
	blackValue int
	whiteValue int
}

type traceEntry struct {
	result                 float64
	eval                   float64
	evalDiff               float64
	phase                  int
	factors                [2]float64
	linearCoefficients     []linearCoefficient
	safetyCoefficients     []safetyCoefficient
	complexityCoefficients []linearCoefficient
	scale                  int
	whiteMove              bool
}

type weight [2]float64

type traceTuner struct {
	k                     float64
	linearWeights         []weight
	safetyWeights         []weight
	safetyWeightsLen      int
	complexityWeights     []weight
	adagrad               []weight
	bestLinearWeights     []weight
	bestSafetyWeights     []weight
	bestComplexityWeights []weight
	entries               []traceEntry
	bestError             float64
	done                  bool
	batchSize             int
}

func printWeights(weightSlices ...[]weight) {
	for _, weights := range weightSlices {
		for _, weight := range weights {
			fmt.Printf("Score(%d, %d), ", int(math.Round(weight[0])), int(math.Round(weight[1])))
		}
	}
	fmt.Println()
}

func (t *traceTuner) computeEvalError() float64 {
	numCPU := runtime.NumCPU()
	results := make([]float64, numCPU)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var c, sum float64
			for y := idx; y < len(t.entries); y += numCPU {
				entry := t.entries[y]
				diff := entry.result - sigmoid(t.k, entry.eval)

				// Kahan summation
				y := (diff * diff) - c
				t := sum + y
				c = (t - sum) - y
				sum = t
			}
			results[idx] = sum
		}(i)
	}
	wg.Wait()
	var sum, c float64
	for _, tResult := range results {
		y := tResult - c
		t := sum + y
		c = (t - sum) - y
		sum = t
	}
	return sum / float64(len(t.entries))
}

func (t *traceTuner) computeLinearError() float64 {
	numCPU := runtime.NumCPU()
	results := make([]float64, numCPU)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var c, sum float64
			for y := idx; y < len(t.entries); y += numCPU {
				entry := t.entries[y]
				diff := entry.result - sigmoid(t.k, entry.evalDiff+t.linearEvaluation(&entry).cp)

				// Kahan summation
				y := (diff * diff) - c
				t := sum + y
				c = (t - sum) - y
				sum = t
			}
			results[idx] = sum
		}(i)
	}
	wg.Wait()
	var sum, c float64
	for _, tResult := range results {
		y := tResult - c
		t := sum + y
		c = (t - sum) - y
		sum = t
	}
	return sum / float64(len(t.entries))
}

func (t *traceTuner) calculateOptimalK() {
	start := -10.0
	end := 10.0
	delta := 1.0
	t.k = start
	best := t.computeEvalError()
	for i := 0; i < 15; i++ {
		t.k = start - delta
		for t.k < end {
			t.k += delta
			err := t.computeEvalError()
			if err <= best {
				best = err
				start = t.k
			}
		}
		end = start + delta
		start = start - delta
		delta /= 10
		fmt.Printf("Optimal k after %d steps: %.17g\n", i+1, start)
	}
	t.k = start
}

func (t *traceTuner) copyCurrentWeightsToBest() {
	copy(t.bestLinearWeights, t.linearWeights)
	copy(t.bestSafetyWeights, t.safetyWeights)
	copy(t.bestComplexityWeights, t.complexityWeights)
}

func TraceTune() {
	t := &traceTuner{done: false}
	t.linearWeights, t.safetyWeights, t.complexityWeights = loadWeights()
	t.safetyWeightsLen = len(t.safetyWeights)
	t.adagrad = make([]weight, len(t.linearWeights)+len(t.safetyWeights)+len(t.complexityWeights))
	t.bestLinearWeights = make([]weight, len(t.linearWeights))
	t.bestSafetyWeights = make([]weight, len(t.safetyWeights))
	t.bestComplexityWeights = make([]weight, len(t.complexityWeights))
	t.copyCurrentWeightsToBest()

	inputChan := make(chan string)
	go loadEntries(inputChan)
	var thread thread
	for fen := range inputChan {
		if entry, ok := t.parseTraceEntry(&thread, fen); ok {
			t.entries = append(t.entries, entry)
		}
	}
	fmt.Println("Number of entries:")
	fmt.Println(len(t.entries))
	t.batchSize = len(t.entries) / 10
	t.calculateOptimalK()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("\nBest values; error: %.17g", t.bestError)
		printWeights(t.bestLinearWeights, t.bestSafetyWeights)
		t.done = true
	}()

	iteration := 0
	iterationsSinceImprovement := 0
	t.bestError = 1e10
	for !t.done {
		rand.Shuffle(len(t.entries), func(i, j int) {
			t.entries[i], t.entries[j] = t.entries[j], t.entries[i]
		})

		for batchStart := 0; batchStart < len(t.entries); batchStart += t.batchSize {
			batch := t.entries[batchStart:Min(len(t.entries)-1, batchStart+t.batchSize)]
			linearGradient, safetyGradient, complexityGradient := t.calculateGradient(batch)
			t.applyGradient(t.linearWeights, linearGradient, t.adagrad)
			t.applyGradient(t.safetyWeights, safetyGradient, t.adagrad[len(t.linearWeights):])
			t.applyGradient(t.complexityWeights, complexityGradient, t.adagrad[len(t.linearWeights)+len(t.safetyWeights):])
		}
		currentError := t.computeLinearError()
		if currentError < t.bestError {
			t.bestError = currentError
			t.copyCurrentWeightsToBest()
			fmt.Printf("Iteration %d error: %.17g\n", iteration, t.bestError)
			printWeights(t.bestLinearWeights, t.bestSafetyWeights, t.bestComplexityWeights)
			iterationsSinceImprovement = 0
		} else {
			iterationsSinceImprovement++
			if iterationsSinceImprovement > 50 {
				break
			}
		}

		iteration++
	}
}

func (t *traceTuner) applyGradient(weights, gradient, adagrad []weight) {
	for idx := range weights {
		for i := Middle; i <= End; i++ {
			adagrad[idx][i] += math.Pow(t.k*gradient[idx][i]/float64(t.batchSize), 2)
			weights[idx][i] -= (t.k / float64(t.batchSize)) * gradient[idx][i] * (learningRate / math.Sqrt(1e-8+adagrad[idx][i]))
		}
	}
}

func (tuner *traceTuner) parseTraceEntry(t *thread, fen string) (traceEntry, bool) {
	var res traceEntry
	sepIdx := strings.Index(fen, ";")
	boardFen := fen[:sepIdx]
	score := fen[sepIdx+1:]
	if strings.Contains(score, "1-0") {
		res.result = 1.0
	} else if strings.Contains(score, "0-1") {
		res.result = 0.0
	} else {
		res.result = 0.5
	}
	board := ParseFen(boardFen)
	t.stack[0].position = board
	t.quiescence(-Mate, Mate, 0, board.IsInCheck())
	for _, move := range t.stack[0].pv.Moves() {
		var child Position
		board.MakeMove(move, &child)
		board = child
	}
	T = Trace{}
	res.eval = float64((&EvaluationContext{}).Evaluate(&board))

	res.scale = T.Scale
	if board.SideToMove == Black {
		res.eval *= -1
	}

	res.whiteMove = board.SideToMove == White
	linearTrace, safetyTrace, complexityTrace := loadTrace()
	for idx, val := range linearTrace {
		if val != 0 {
			res.linearCoefficients = append(res.linearCoefficients, linearCoefficient{idx: idx, value: val})
		}
	}
	for idx, val := range safetyTrace {
		if val[White] != 0 || val[Black] != 0 {
			res.safetyCoefficients = append(res.safetyCoefficients, safetyCoefficient{idx: idx, blackValue: val[Black], whiteValue: val[White]})
		}
	}
	for idx, val := range complexityTrace {
		if val != 0 {
			res.complexityCoefficients = append(res.complexityCoefficients, linearCoefficient{idx: idx, value: val})
		}
	}

	res.phase = (TotalPhase - QueenPhase*PopCount(board.Pieces[Queen]) -
		RookPhase*PopCount(board.Pieces[Rook]) -
		BishopPhase*PopCount(board.Pieces[Bishop]) -
		KnightPhase*PopCount(board.Pieces[Knight]))

	if res.phase < 0 {
		res.phase = 0
	}

	res.factors[Middle] = 1.0 - float64(res.phase)/float64(TotalPhase)
	res.factors[End] = float64(res.phase) / float64(TotalPhase)
	res.phase = (res.phase*256 + (TotalPhase / 2)) / TotalPhase

	res.evalDiff = res.eval - tuner.linearEvaluation(&res).cp

	// if math.Abs(res.evalDiff) > 1 {
	// 	fmt.Println("Problem with evaluation", res.evalDiff, fen)
	// }

	return res, true
}

func (t *traceTuner) calculateGradient(entries []traceEntry) ([]weight, []weight, []weight) {
	numCPU := runtime.NumCPU()
	type weightTuple struct {
		liner      []weight
		safety     []weight
		complexity []weight
	}
	linearRes := make([]weight, len(t.linearWeights))
	safetyRes := make([]weight, len(t.safetyWeights))
	complexityRes := make([]weight, len(t.complexityWeights))

	resultChan := make(chan weightTuple)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()
			localLinearRes := make([]weight, len(t.linearWeights))
			localSafetyRes := make([]weight, len(t.safetyWeights))
			localComplexityRes := make([]weight, len(t.complexityWeights))
			for y := idx; y < len(entries); y += numCPU {
				entry := entries[y]
				evaluationResult := t.linearEvaluation(&entry)
				derivative := t.singleLinearDerivative(&entry, evaluationResult.cp)
				middleMultiplier := entry.factors[Middle] * derivative
				endMultiplier := entry.factors[End] * (float64(entry.scale) / float64(ScaleNormal)) * derivative
				canUpdateEndgame := evaluationResult.endGameEval == 0 || evaluationResult.complexity >= -Abs(int(evaluationResult.endGameEval))
				complexitySign := BoolToInt(evaluationResult.endGameEval > 0) - BoolToInt(evaluationResult.endGameEval < 0)
				for _, coef := range entry.linearCoefficients {
					localLinearRes[coef.idx][Middle] += float64(coef.value) * middleMultiplier
					if canUpdateEndgame {
						localLinearRes[coef.idx][End] += float64(coef.value) * endMultiplier
					}
				}
				for coefIdx, coef := range entry.safetyCoefficients {
					// King safety attack value
					if coef.idx == t.safetyWeightsLen-1 {
						whiteScale := float64(coef.whiteValue) / float64(entry.safetyCoefficients[coefIdx+1].whiteValue)
						if entry.safetyCoefficients[coefIdx+1].whiteValue == 0 {
							whiteScale = 0.0
						}
						blackScale := float64(coef.blackValue) / float64(entry.safetyCoefficients[coefIdx+1].blackValue)
						if entry.safetyCoefficients[coefIdx+1].blackValue == 0 {
							blackScale = 0.0
						}
						localSafetyRes[coef.idx][Middle] += (math.Max(float64(evaluationResult.safetyBlack.Middle()), 0)*blackScale -
							math.Max(float64(evaluationResult.safetyWhite.Middle()), 0)*whiteScale) * (middleMultiplier / 360)
						if canUpdateEndgame {
							localSafetyRes[coef.idx][End] += (sign(float64(evaluationResult.safetyBlack.End()))*blackScale -
								sign(float64(evaluationResult.safetyWhite.End()))*whiteScale) * (endMultiplier / 20)

						}
						break
					}
					localSafetyRes[coef.idx][Middle] += (math.Max(float64(evaluationResult.safetyBlack.Middle()), 0)*float64(coef.blackValue) -
						math.Max(float64(evaluationResult.safetyWhite.Middle()), 0)*float64(coef.whiteValue)) * (middleMultiplier / 360)
					if canUpdateEndgame {
						localSafetyRes[coef.idx][End] += (sign(float64(evaluationResult.safetyBlack.End()))*float64(coef.blackValue) -
							sign(float64(evaluationResult.safetyWhite.End()))*float64(coef.whiteValue)) * (endMultiplier / 20)
					}
				}
				for _, coef := range entry.complexityCoefficients {
					if canUpdateEndgame && evaluationResult.endGameEval != 0 {
						localComplexityRes[coef.idx][End] += float64(coef.value) * endMultiplier * float64(complexitySign)
					}
				}
			}
			resultChan <- weightTuple{localLinearRes, localSafetyRes, localComplexityRes}
		}(i)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for threadResult := range resultChan {
		for idx := range linearRes {
			for i := Middle; i <= End; i++ {
				linearRes[idx][i] += threadResult.liner[idx][i]
			}
		}
		for idx := range safetyRes {
			for i := Middle; i <= End; i++ {
				safetyRes[idx][i] += threadResult.safety[idx][i]
			}
		}
		for idx := range complexityRes {
			for i := Middle; i <= End; i++ {
				safetyRes[idx][i] += threadResult.complexity[idx][i]
			}
		}
	}
	return linearRes, safetyRes, complexityRes
}

func sign(x float64) float64 {
	if x == 0 {
		return 0
	} else if x > 0 {
		return 1
	} else {
		return -1
	}
}

type linearEvaluationResult struct {
	cp          float64
	safetyBlack Score
	safetyWhite Score
	endGameEval int16
	complexity  int
}

func (t *traceTuner) linearEvaluation(entry *traceEntry) linearEvaluationResult {
	var middle, end int
	var safetyWhite, safetyBlack Score
	for _, coeff := range entry.linearCoefficients {
		middle += int(math.Round(t.linearWeights[coeff.idx][Middle])) * coeff.value
		end += int(math.Round(t.linearWeights[coeff.idx][End])) * coeff.value
	}
	for traceIdx, coeff := range entry.safetyCoefficients {
		multiplier := S(int16(math.Round(t.safetyWeights[coeff.idx][Middle])), int16(math.Round(t.safetyWeights[coeff.idx][End])))
		if coeff.idx != t.safetyWeightsLen-1 {
			safetyBlack += multiplier * Score(coeff.blackValue)
			safetyWhite += multiplier * Score(coeff.whiteValue)
			continue
		}

		// King safety attack value
		if entry.safetyCoefficients[traceIdx+1].blackValue != 0 {
			safetyBlack += S(
				int16(int(multiplier.Middle())*coeff.blackValue/entry.safetyCoefficients[traceIdx+1].blackValue),
				int16(int(multiplier.End())*coeff.blackValue/entry.safetyCoefficients[traceIdx+1].blackValue),
			)
		}
		if entry.safetyCoefficients[traceIdx+1].whiteValue != 0 {
			safetyWhite += S(
				int16(int(multiplier.Middle())*coeff.whiteValue/entry.safetyCoefficients[traceIdx+1].whiteValue),
				int16(int(multiplier.End())*coeff.whiteValue/entry.safetyCoefficients[traceIdx+1].whiteValue),
			)
		}
		break
	}

	var complexity int
	for _, coeff := range entry.complexityCoefficients {
		complexity += int(math.Round(t.complexityWeights[coeff.idx][End])) * coeff.value
	}
	score := S(int16(middle), int16(end))
	middleWhite := int(safetyWhite.Middle())
	endWhite := int(safetyWhite.End())
	middleBlack := int(safetyBlack.Middle())
	endBlack := int(safetyBlack.End())
	score += S(
		int16(-middleWhite*Max(middleWhite, 0)/720),
		-int16(Max(endWhite, 0)/20),
	)
	score -= S(
		int16(-middleBlack*Max(middleBlack, 0)/720),
		-int16(Max(endBlack, 0)/20),
	)

	endGameEval := score.End()

	sign := BoolToInt(endGameEval > 0) - BoolToInt(endGameEval < 0)
	score += S(0, int16(sign*Max(-Abs(int(score.End())), complexity)))

	phased := (int(score.Middle())*(256-entry.phase) + (int(score.End())*entry.phase*entry.scale)/ScaleNormal) / 256
	if entry.whiteMove {
		return linearEvaluationResult{float64(phased + int(Tempo)), safetyBlack, safetyWhite, endGameEval, complexity}
	} else {
		return linearEvaluationResult{float64(phased - int(Tempo)), safetyBlack, safetyWhite, endGameEval, complexity}
	}

}

func (t *traceTuner) singleLinearDerivative(entry *traceEntry, linearEvaluation float64) float64 {
	sigma := sigmoid(t.k, entry.evalDiff+linearEvaluation)
	sigmaPrim := sigma * (1 - sigma)
	return -((entry.result - sigma) * sigmaPrim)
}

func loadTrace() (linearRes []int, safetyRes [][2]int, complexityRes []int) {
	linearRes = append(linearRes, T.PawnValue)
	linearRes = append(linearRes, T.KnightValue)
	linearRes = append(linearRes, T.BishopValue)
	linearRes = append(linearRes, T.RookValue)
	linearRes = append(linearRes, T.QueenValue)

	for flag := 0; flag <= 15; flag++ {
		for y := 1; y < 7; y++ {
			for x := 0; x < 8; x++ {
				linearRes = append(linearRes, T.PawnScores[flag][y][x])
			}
		}
	}
	for i := Knight; i <= King; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				linearRes = append(linearRes, T.PieceScores[i][y][x])
			}
		}
	}

	for y := 0; y < 7; y++ {
		for x := 0; x < 4; x++ {
			linearRes = append(linearRes, T.PawnsConnected[y][x])
		}
	}
	for y := 0; y < 9; y++ {
		linearRes = append(linearRes, T.MobilityBonus[0][y])
	}
	for y := 0; y < 14; y++ {
		linearRes = append(linearRes, T.MobilityBonus[1][y])
	}
	for y := 0; y < 15; y++ {
		linearRes = append(linearRes, T.MobilityBonus[2][y])
	}
	for y := 0; y < 28; y++ {
		linearRes = append(linearRes, T.MobilityBonus[3][y])
	}
	for y := 0; y < 8; y++ {
		linearRes = append(linearRes, T.PassedFriendlyDistance[y])
	}
	for y := 0; y < 8; y++ {
		linearRes = append(linearRes, T.PassedEnemyDistance[y])
	}
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for c := 0; c < 2; c++ {
				for y := 0; y < 7; y++ {
					linearRes = append(linearRes, T.PassedRank[a][b][c][y])
				}
			}
		}
	}
	for y := 0; y < 8; y++ {
		linearRes = append(linearRes, T.PassedFile[y])
	}
	for y := 0; y < 7; y++ {
		linearRes = append(linearRes, T.PassedStacked[y])
	}
	for y := 0; y < 6; y++ {
		linearRes = append(linearRes, T.PassedUncontested[y])
	}
	for y := 0; y < 6; y++ {
		linearRes = append(linearRes, T.PassedPushDefended[y])
	}
	for y := 0; y < 6; y++ {
		linearRes = append(linearRes, T.PassedPushUncontestedDefended[y])
	}
	linearRes = append(linearRes, T.Isolated)
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for file := FileA; file <= FileH; file++ {
				linearRes = append(linearRes, T.StackedPawns[a][b][file])
			}
		}
	}
	linearRes = append(linearRes, T.AttackedBySafePawn[:]...)
	linearRes = append(linearRes, T.Backward)
	linearRes = append(linearRes, T.BackwardOpen)
	linearRes = append(linearRes, T.BishopPair)
	linearRes = append(linearRes, T.BishopRammedPawns)
	linearRes = append(linearRes, T.BishopOutpostUndefendedBonus)
	linearRes = append(linearRes, T.BishopOutpostDefendedBonus)
	linearRes = append(linearRes, T.LongDiagonalBishop)
	for y := 0; y < 4; y++ {
		linearRes = append(linearRes, T.DistantBishop[y])
	}
	linearRes = append(linearRes, T.KnightOutpostUndefendedBonus)
	linearRes = append(linearRes, T.KnightOutpostDefendedBonus)
	for y := 0; y < 4; y++ {
		linearRes = append(linearRes, T.DistantKnight[y])
	}
	linearRes = append(linearRes, T.MinorBehindPawn)
	linearRes = append(linearRes, T.RookOnFile[0])
	linearRes = append(linearRes, T.RookOnFile[1])
	linearRes = append(linearRes, T.RookOnQueenFile)
	linearRes = append(linearRes, T.TrappedRook)
	for y := 0; y < 12; y++ {
		linearRes = append(linearRes, T.KingDefenders[y])
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				linearRes = append(linearRes, T.KingShelter[x][y][z])
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				linearRes = append(linearRes, T.KingStorm[x][y][z])
			}
		}
	}
	linearRes = append(linearRes, T.KingOnPawnlessFlank)

	linearRes = append(linearRes, T.Hanging)
	linearRes = append(linearRes, T.ThreatByKing)
	for i := Pawn; i <= King; i++ {
		linearRes = append(linearRes, T.ThreatByMinor[i])
	}
	for i := Pawn; i <= King; i++ {
		linearRes = append(linearRes, T.ThreatByRook[i])
	}

	for flag := 0; flag <= 15; flag++ {
		linearRes = append(linearRes, T.RookBishopExistence[flag])
	}
	for flag := 0; flag <= 15; flag++ {
		linearRes = append(linearRes, T.QueenBishopExistence[flag])
	}

	for flag := 0; flag <= 15; flag++ {
		linearRes = append(linearRes, T.KingBishopExistence[flag])
	}

	//
	// King Safety
	//
	for x := Pawn; x <= Queen; x++ {
		safetyRes = append(safetyRes, [2]int{T.KingSafetyAttacksWeights[Black][x], T.KingSafetyAttacksWeights[White][x]})
	}
	safetyRes = append(safetyRes, T.KingSafetyWeakSquares)
	safetyRes = append(safetyRes, T.KingSafetyFriendlyPawns)
	safetyRes = append(safetyRes, T.KingSafetyNoEnemyQueens)
	safetyRes = append(safetyRes, T.KingSafetySafeQueenCheck)
	safetyRes = append(safetyRes, T.KingSafetySafeRookCheck)
	safetyRes = append(safetyRes, T.KingSafetySafeBishopCheck)
	safetyRes = append(safetyRes, T.KingSafetySafeKnightCheck)
	safetyRes = append(safetyRes, T.KingSafetyAdjustment)

	safetyRes = append(safetyRes, T.KingSafetyAttackValueNumerator)
	safetyRes = append(safetyRes, T.KingSafetyAttackValueDenumerator)

	//
	// Complexity
	//
	complexityRes = append(complexityRes, T.ComplexityTotalPawns)
	complexityRes = append(complexityRes, T.ComplexityPawnEndgame)
	complexityRes = append(complexityRes, T.ComplexityPawnBothFlanks)
	complexityRes = append(complexityRes, T.ComplexityInfiltration)
	complexityRes = append(complexityRes, T.ComplexityAdjustment)

	return
}

func scoreToWeight(s Score) weight {
	return weight{float64(s.Middle()), float64(s.End())}
}

func scoresToWeights(scores []Score) []weight {
	res := make([]weight, 0, len(scores))
	for _, s := range scores {
		res = append(res, scoreToWeight(s))
	}
	return res
}

func loadWeights() ([]weight, []weight, []weight) {
	var linearScores []Score
	linearScores = append(linearScores, PawnValue)
	linearScores = append(linearScores, KnightValue)
	linearScores = append(linearScores, BishopValue)
	linearScores = append(linearScores, RookValue)
	linearScores = append(linearScores, QueenValue)

	for flag := 0; flag <= 15; flag++ {
		for y := 1; y < 7; y++ {
			for x := 0; x < 8; x++ {
				linearScores = append(linearScores, PawnScores[flag][y][x])
			}
		}
	}
	for i := Knight; i <= King; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				linearScores = append(linearScores, PieceScores[i][y][x])
			}
		}
	}

	for y := 0; y < 7; y++ {
		for x := 0; x < 4; x++ {
			linearScores = append(linearScores, PawnsConnected[y][x])
		}
	}
	for y := 0; y < 9; y++ {
		linearScores = append(linearScores, MobilityBonus[0][y])
	}
	for y := 0; y < 14; y++ {
		linearScores = append(linearScores, MobilityBonus[1][y])
	}
	for y := 0; y < 15; y++ {
		linearScores = append(linearScores, MobilityBonus[2][y])
	}
	for y := 0; y < 28; y++ {
		linearScores = append(linearScores, MobilityBonus[3][y])
	}
	for y := 0; y < 8; y++ {
		linearScores = append(linearScores, PassedFriendlyDistance[y])
	}
	for y := 0; y < 8; y++ {
		linearScores = append(linearScores, PassedEnemyDistance[y])
	}
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for c := 0; c < 2; c++ {
				for y := 0; y < 7; y++ {
					linearScores = append(linearScores, PassedRank[a][b][c][y])
				}
			}
		}
	}
	for y := 0; y < 8; y++ {
		linearScores = append(linearScores, PassedFile[y])
	}
	for y := 0; y < 7; y++ {
		linearScores = append(linearScores, PassedStacked[y])
	}
	for y := 0; y < 6; y++ {
		linearScores = append(linearScores, PassedUncontested[y])
	}
	for y := 0; y < 6; y++ {
		linearScores = append(linearScores, PassedPushDefended[y])
	}
	for y := 0; y < 6; y++ {
		linearScores = append(linearScores, PassedPushUncontestedDefended[y])
	}
	linearScores = append(linearScores, Isolated)
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for file := FileA; file <= FileH; file++ {
				linearScores = append(linearScores, StackedPawns[a][b][file])
			}
		}
	}
	linearScores = append(linearScores, AttackedBySafePawn[:]...)
	linearScores = append(linearScores, Backward)
	linearScores = append(linearScores, BackwardOpen)
	linearScores = append(linearScores, BishopPair)
	linearScores = append(linearScores, BishopRammedPawns)
	linearScores = append(linearScores, BishopOutpostUndefendedBonus)
	linearScores = append(linearScores, BishopOutpostDefendedBonus)
	linearScores = append(linearScores, LongDiagonalBishop)
	for y := 0; y < 4; y++ {
		linearScores = append(linearScores, DistantBishop[y])
	}
	linearScores = append(linearScores, KnightOutpostUndefendedBonus)
	linearScores = append(linearScores, KnightOutpostDefendedBonus)
	for y := 0; y < 4; y++ {
		linearScores = append(linearScores, DistantKnight[y])
	}
	linearScores = append(linearScores, MinorBehindPawn)
	linearScores = append(linearScores, RookOnFile[0])
	linearScores = append(linearScores, RookOnFile[1])
	linearScores = append(linearScores, RookOnQueenFile)
	linearScores = append(linearScores, TrappedRook)
	for y := 0; y < 12; y++ {
		linearScores = append(linearScores, KingDefenders[y])
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				linearScores = append(linearScores, KingShelter[x][y][z])
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				linearScores = append(linearScores, KingStorm[x][y][z])
			}
		}
	}
	linearScores = append(linearScores, KingOnPawnlessFlank)
	linearScores = append(linearScores, Hanging)
	linearScores = append(linearScores, ThreatByKing)
	for x := Pawn; x <= King; x++ {
		linearScores = append(linearScores, ThreatByMinor[x])
	}
	for x := Pawn; x <= King; x++ {
		linearScores = append(linearScores, ThreatByRook[x])
	}

	for flag := 0; flag <= 15; flag++ {
		linearScores = append(linearScores, RookBishopExistence[flag])
	}

	for flag := 0; flag <= 15; flag++ {
		linearScores = append(linearScores, QueenBishopExistence[flag])
	}

	for flag := 0; flag <= 15; flag++ {
		linearScores = append(linearScores, KingBishopExistence[flag])
	}

	//
	// Safety
	//

	var safetyScores []Score
	for x := Pawn; x <= Queen; x++ {
		safetyScores = append(safetyScores, KingSafetyAttacksWeights[x])
	}
	safetyScores = append(safetyScores, KingSafetyWeakSquares)
	safetyScores = append(safetyScores, KingSafetyFriendlyPawns)
	safetyScores = append(safetyScores, KingSafetyNoEnemyQueens)
	safetyScores = append(safetyScores, KingSafetySafeQueenCheck)
	safetyScores = append(safetyScores, KingSafetySafeRookCheck)
	safetyScores = append(safetyScores, KingSafetySafeBishopCheck)
	safetyScores = append(safetyScores, KingSafetySafeKnightCheck)
	safetyScores = append(safetyScores, KingSafetyAdjustment)

	safetyScores = append(safetyScores, KingSafetyAttackValue)

	//
	// Complexity
	//

	var complexityScores []Score
	complexityScores = append(complexityScores, ComplexityTotalPawns)
	complexityScores = append(complexityScores, ComplexityPawnEndgame)
	complexityScores = append(complexityScores, ComplexityPawnBothFlanks)
	complexityScores = append(complexityScores, ComplexityInfiltration)
	complexityScores = append(complexityScores, ComplexityAdjustment)

	return scoresToWeights(linearScores), scoresToWeights(safetyScores), scoresToWeights(complexityScores)
}
