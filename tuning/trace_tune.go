package tuning

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
	MIDDLE = iota
	END
)

const learningRate = 10.0

type coefficient struct {
	value int
	idx   int
}

type traceEntry struct {
	result       float64
	eval         float64
	phase        int
	factors      [2]float64
	coefficients []coefficient
	evalDiff     float64
}

type weight [2]float64

type traceTuner struct {
	k                       float64
	weights                 []weight
	bestWeights             []weight
	entries                 []traceEntry
	bestError               float64
	bestErrorRegularization float64
	done                    bool
	batchSize               int
}

func printWeights(weights []weight) {
	for _, weight := range weights {
		fmt.Printf("Score(%d, %d), ", int(math.Round(weight[0])), int(math.Round(weight[1])))
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
				diff := entry.result - sigmoid(t.k, entry.evalDiff+t.linearEvaluation(&entry))

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

func TraceTune() {
	t := &traceTuner{done: false, batchSize: 16384 * 2}
	t.weights = loadWeights()
	t.bestWeights = make([]weight, len(t.weights))
	copy(t.bestWeights, t.weights)

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
	t.calculateOptimalK()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("\nBest values; error: %.17g", t.bestError)
		printWeights(t.bestWeights)
		t.done = true
	}()

	iteration := 0
	t.bestError = 1e10
	for !t.done {
		rand.Shuffle(len(t.entries), func(i, j int) {
			t.entries[i], t.entries[j] = t.entries[j], t.entries[i]
		})

		for batchStart := 0; batchStart < len(t.entries); batchStart += t.batchSize {
			batch := t.entries[batchStart:Min(len(t.entries)-1, batchStart+t.batchSize)]
			gradient := t.calculateGradient(batch)
			for idx := range t.weights {
				for i := MIDDLE; i <= END; i++ {
					t.weights[idx][i] -= (learningRate / float64(t.batchSize)) * gradient[idx][i]
				}
			}
		}
		currentError := t.computeLinearError()
		if currentError < t.bestError {
			t.bestError = currentError
			copy(t.bestWeights, t.weights)
			fmt.Printf("Iteration %d error: %.17g regularization: %.17g\n", iteration, t.bestError, t.regularization())
			printWeights(t.bestWeights)
		} else {
			break
		}

		iteration++
	}
}

const regularizationWeight = 0.2e-7

func (t *traceTuner) regularization() float64 {
	sum := 0.0
	for _, weight := range t.weights {
		sum += math.Abs(weight[0])
		sum += math.Abs(weight[1])
	}
	return sum * regularizationWeight
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
	res.eval = float64(Evaluate(&board))

	// Do not care about scaled positions
	if ScaleFactor(&board, int16(res.eval)) != SCALE_NORMAL {
		return res, false
	}

	if board.SideToMove == Black {
		res.eval *= -1
	}

	for idx, val := range loadTrace() {
		if val != 0 {
			res.coefficients = append(res.coefficients, coefficient{idx: idx, value: val})
		}
	}

	res.phase = (TotalPhase - QueenPhase*PopCount(board.Pieces[Queen]) -
		RookPhase*PopCount(board.Pieces[Rook]) -
		BishopPhase*PopCount(board.Pieces[Bishop]) -
		KnightPhase*PopCount(board.Pieces[Knight]))

	if res.phase < 0 {
		res.phase = 0
	}

	res.factors[MIDDLE] = 1.0 - float64(res.phase)/float64(TotalPhase)
	res.factors[END] = float64(res.phase) / float64(TotalPhase)
	res.phase = (res.phase*256 + (TotalPhase / 2)) / TotalPhase

	res.evalDiff = res.eval - tuner.linearEvaluation(&res)

	return res, true
}

func (t *traceTuner) calculateGradient(entries []traceEntry) []weight {
	numCPU := runtime.NumCPU()
	res := make([]weight, len(t.weights))

	resultChan := make(chan []weight)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()
			localRes := make([]weight, len(t.weights))
			for y := idx; y < len(entries); y += numCPU {
				entry := entries[y]
				derivative := t.singleLinearDerivative(&entry)
				for _, coef := range entry.coefficients {
					for i := MIDDLE; i <= END; i++ {
						localRes[coef.idx][i] += derivative * entry.factors[i] * float64(coef.value)
					}
				}
			}
			resultChan <- localRes
		}(i)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for threadResult := range resultChan {
		for idx := range res {
			for i := MIDDLE; i <= END; i++ {
				res[idx][i] += threadResult[idx][i]
			}
		}
	}

	for idx := range res {
		for i := MIDDLE; i <= END; i++ {
			res[idx][i] += sign(t.weights[idx][i]) * regularizationWeight
		}
	}
	return res
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

func (t *traceTuner) linearEvaluation(entry *traceEntry) float64 {
	var middle, end float64
	for _, coeff := range entry.coefficients {
		middle += t.weights[coeff.idx][MIDDLE] * float64(coeff.value)
		end += t.weights[coeff.idx][END] * float64(coeff.value)
	}
	return (middle*(256.0-float64(entry.phase)) + end*float64(entry.phase)) / 256.0
}

// Optimization taken from Ethereal
// d/dx (y - sigmoid(k, x))^2 = C * (y - sigmoid(k, x)) * sigmoid(k, x) * (1 - sigmoid(k, x))
// where C is a negative contant
// https://www.wolframalpha.com/input/?i=d%2Fdx+%28y-+%281%2F%281+%2B+10%5E%28-400kx%29%29%29%29%5E2
// https://www.wolframalpha.com/input/?i=%281%2F%281+%2B+10%5E%28-400kx%29%29%29+*+%281+-+1%2F%281+%2B+10%5E%28-400kx%29%29+%29
// https://www.wolframalpha.com/input/?i=%28d%2Fdx+%28y-+%281%2F%281+%2B+10%5E%28-400kx%29%29%29%29%5E2%29+%2F+%28%281%2F%281+%2B+10%5E%28-400kx%29%29%29+*+%281+-+1%2F%281+%2B+10%5E%28-400kx%29%29+%29+*+%28y+-+1%2F%281+%2B+10%5E%28-400kx%29%29%29%29
func (t *traceTuner) singleLinearDerivative(entry *traceEntry) float64 {
	sigma := sigmoid(t.k, entry.evalDiff+t.linearEvaluation(entry))
	sigmaPrim := sigma * (1 - sigma)
	return -((entry.result - sigma) * sigmaPrim)
}

func loadTrace() (res []int) {
	res = append(res, T.PawnValue)
	res = append(res, T.KnightValue)
	res = append(res, T.BishopValue)
	res = append(res, T.RookValue)
	res = append(res, T.QueenValue)

	for i := Knight; i <= King; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				res = append(res, T.PieceScores[i][y][x])
			}
		}
	}
	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			res = append(res, T.PawnScores[y][x])
		}
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 4; x++ {
			res = append(res, T.PawnsConnected[y][x])
		}
	}
	for y := 0; y < 9; y++ {
		res = append(res, T.MobilityBonus[0][y])
	}
	for y := 0; y < 14; y++ {
		res = append(res, T.MobilityBonus[1][y])
	}
	for y := 0; y < 15; y++ {
		res = append(res, T.MobilityBonus[2][y])
	}
	for y := 0; y < 28; y++ {
		res = append(res, T.MobilityBonus[3][y])
	}
	for y := 0; y < 8; y++ {
		res = append(res, T.PassedFriendlyDistance[y])
	}
	for y := 0; y < 8; y++ {
		res = append(res, T.PassedEnemyDistance[y])
	}
	for y := 0; y < 7; y++ {
		res = append(res, T.PassedRank[y])
	}
	for y := 0; y < 8; y++ {
		res = append(res, T.PassedFile[y])
	}
	for y := 0; y < 8; y++ {
		res = append(res, T.PassedStacked[y])
	}
	res = append(res, T.Isolated)
	res = append(res, T.Doubled)
	res = append(res, T.Backward)
	res = append(res, T.BackwardOpen)
	res = append(res, T.BishopPair)
	res = append(res, T.BishopRammedPawns)
	res = append(res, T.BishopOutpostUndefendedBonus)
	res = append(res, T.BishopOutpostDefendedBonus)
	res = append(res, T.LongDiagonalBishop)
	res = append(res, T.KnightOutpostUndefendedBonus)
	res = append(res, T.KnightOutpostDefendedBonus)
	res = append(res, T.PotentialKnightOutpostUndefendedBonus)
	res = append(res, T.PotentialKnightOutpostDefendedBonus)
	for y := 0; y < 4; y++ {
		res = append(res, T.DistantKnight[y])
	}
	res = append(res, T.MinorBehindPawn)
	res = append(res, T.RookOnFile[0])
	res = append(res, T.RookOnFile[1])
	res = append(res, T.RookOnQueenFile)
	for y := 0; y < 12; y++ {
		res = append(res, T.KingDefenders[y])
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, T.KingShelter[x][y][z])
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, T.KingStorm[x][y][z])
			}
		}
	}
	res = append(res, T.KingOnPawnlessFlank)

	res = append(res, T.Hanging)
	res = append(res, T.ThreatByKing)
	for i := 0; i <= King; i++ {
		res = append(res, T.ThreatByMinor[i])
	}
	for i := 0; i <= King; i++ {
		res = append(res, T.ThreatByRook[i])
	}

	return
}

func scoreToWeight(s Score) weight {
	return weight{float64(s.Middle()), float64(s.End())}
}

func loadWeights() []weight {
	var tmp []Score
	tmp = append(tmp, PawnValue)
	tmp = append(tmp, KnightValue)
	tmp = append(tmp, BishopValue)
	tmp = append(tmp, RookValue)
	tmp = append(tmp, QueenValue)

	for i := Knight; i <= King; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				tmp = append(tmp, PieceScores[i][y][x])
			}
		}
	}
	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			tmp = append(tmp, PawnScores[y][x])
		}
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 4; x++ {
			tmp = append(tmp, PawnsConnected[y][x])
		}
	}
	for y := 0; y < 9; y++ {
		tmp = append(tmp, MobilityBonus[0][y])
	}
	for y := 0; y < 14; y++ {
		tmp = append(tmp, MobilityBonus[1][y])
	}
	for y := 0; y < 15; y++ {
		tmp = append(tmp, MobilityBonus[2][y])
	}
	for y := 0; y < 28; y++ {
		tmp = append(tmp, MobilityBonus[3][y])
	}
	for y := 0; y < 8; y++ {
		tmp = append(tmp, PassedFriendlyDistance[y])
	}
	for y := 0; y < 8; y++ {
		tmp = append(tmp, PassedEnemyDistance[y])
	}
	for y := 0; y < 7; y++ {
		tmp = append(tmp, PassedRank[y])
	}
	for y := 0; y < 8; y++ {
		tmp = append(tmp, PassedFile[y])
	}
	for y := 0; y < 8; y++ {
		tmp = append(tmp, PassedStacked[y])
	}
	tmp = append(tmp, Isolated)
	tmp = append(tmp, Doubled)
	tmp = append(tmp, Backward)
	tmp = append(tmp, BackwardOpen)
	tmp = append(tmp, BishopPair)
	tmp = append(tmp, BishopRammedPawns)
	tmp = append(tmp, BishopOutpostUndefendedBonus)
	tmp = append(tmp, BishopOutpostDefendedBonus)
	tmp = append(tmp, LongDiagonalBishop)
	tmp = append(tmp, KnightOutpostUndefendedBonus)
	tmp = append(tmp, KnightOutpostDefendedBonus)
	tmp = append(tmp, PotentialKnightOutpostUndefendedBonus)
	tmp = append(tmp, PotentialKnightOutpostDefendedBonus)
	for y := 0; y < 4; y++ {
		tmp = append(tmp, DistantKnight[y])
	}
	tmp = append(tmp, MinorBehindPawn)
	tmp = append(tmp, RookOnFile[0])
	tmp = append(tmp, RookOnFile[1])
	tmp = append(tmp, RookOnQueenFile)
	for y := 0; y < 12; y++ {
		tmp = append(tmp, KingDefenders[y])
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				tmp = append(tmp, KingShelter[x][y][z])
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				tmp = append(tmp, KingStorm[x][y][z])
			}
		}
	}
	tmp = append(tmp, KingOnPawnlessFlank)
	tmp = append(tmp, Hanging)
	tmp = append(tmp, ThreatByKing)
	for x := Pawn; x <= King; x++ {
		tmp = append(tmp, ThreatByMinor[x])
	}
	for x := Pawn; x <= King; x++ {
		tmp = append(tmp, ThreatByRook[x])
	}

	res := make([]weight, 0, len(tmp))
	for _, s := range tmp {
		res = append(res, scoreToWeight(s))
	}

	fmt.Println(res)

	return res
}
