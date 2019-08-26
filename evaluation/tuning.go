package evaluation

import (
	"bufio"
	"fmt"
	. "github.com/mhib/combusken/backend"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type tuneEntry struct {
	Position
	result float64
}

type thread struct {
	stack [127]stackEntry
}

type stackEntry struct {
	position Position
	pv
	moves [256]EvaledMove
}

type pv struct {
	size  int
	items [127]Move
}

func (pv *pv) clear() {
	pv.size = 0
}

func (pv *pv) assign(m Move, child *pv) {
	pv.size = 1 + child.size
	pv.items[0] = m
	copy(pv.items[1:], child.Moves())
}

func (pv *pv) Moves() []Move {
	return pv.items[:pv.size]
}

type emptyPKTableType struct {
}

func (t *emptyPKTableType) Get(uint64) (bool, Score) {
	return false, Score{}
}

func (t *emptyPKTableType) Set(uint64, Score) {
}

func (t *emptyPKTableType) Clear() {
}

var emptyPKTable = emptyPKTableType{}

// Copy if quiescence search to extract quiet position
func (t *thread) quiescence(alpha, beta, height int, inCheck bool) int {
	t.stack[height].pv.clear()
	pos := &t.stack[height].position

	if height >= 127 {
		return 0
	}

	child := &t.stack[height+1].position

	moveCount := 0

	val := Evaluate(pos, &emptyPKTable)

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

	for i := range evaled {
		if (!inCheck && !SeeSign(pos, evaled[i].Move)) || !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		moveCount++
		childInCheck := child.IsInCheck()
		val = -t.quiescence(-beta, -alpha, height+1, childInCheck)
		if val > alpha {
			alpha = val
			if val >= beta {
				return beta
			}
			t.stack[height].pv.assign(evaled[i].Move, &t.stack[height+1].pv)
		}
	}

	if moveCount == 0 && inCheck {
		return -Mate + height
	}
	return alpha
}

type tuner struct {
	k       float64
	weights []EvaluationValue
	entries []tuneEntry
}

func Tune() {
	inputChan := make(chan string)
	go loadEntries(inputChan)
	wg := &sync.WaitGroup{}
	resultChan := make(chan tuneEntry)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var t thread
			for fen := range inputChan {
				parseEntry(&t, fen, resultChan)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	t := &tuner{}
	for entry := range resultChan {
		t.entries = append(t.entries, entry)
	}
	fmt.Println("Number of entries:")
	fmt.Println(len(t.entries))
	t.calculateOptimalK()
	fmt.Printf("Optimal k: %.17g\n", t.k)
	t.weights = loadScoresToSlice()
	t.coordinateDescent()
	fmt.Println("Generated weights in coordinate descent:")
	fmt.Println(t.weights)
}

func (t *tuner) computeError() float64 {
	numCPU := runtime.NumCPU()
	results := make([]float64, numCPU)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for y := idx; y < len(t.entries); y += numCPU {
				entry := t.entries[y]
				evaluation := float64(Evaluate(&entry.Position, &emptyPKTable))
				if !entry.Position.WhiteMove {
					evaluation *= -1
				}
				diff := entry.result - sigmoid(t.k, evaluation)
				results[idx] += diff * diff
			}
		}(i)
	}
	wg.Wait()
	sum := 0.0
	for _, tResult := range results {
		sum += tResult
	}
	return sum / float64(len(t.entries))
}

func (t *tuner) calculateOptimalK() {
	start := -10.0
	end := 10.0
	delta := 1.0
	t.k = start
	best := t.computeError()
	for i := 0; i < 10; i++ {
		t.k = start - delta
		for t.k < end {
			t.k += delta
			err := t.computeError()
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

func parseEntry(t *thread, fen string, resultChan chan tuneEntry) {
	var res tuneEntry
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
	res.Position = board

	resultChan <- res
}

func loadEntries(inputChan chan string) {
	defer close(inputChan)
	absPath, _ := filepath.Abs("./games.fen")
	file, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputChan <- scanner.Text()
	}
	return
}

func sigmoid(K, S float64) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, -K*S/400.0))
}

func absScore(num int16) int {
	if num < 0 {
		return -int(num)
	}
	return int(num)
}

func (t *tuner) regularization() float64 {
	alpha := 0.2e-7
	sum := 0
	for _, score := range t.weights {
		if score.regularized() {
			for i := 0; i < score.phaseCount(); i++ {
				sum += absScore(score.get(i))
			}
		}
	}
	return alpha * float64(sum)
}

// I've tried annealing but, I was unable to make it work
//func (t *tuner) annealing() (improved bool) {
//energy := func() float64 {
//return (t.computeError() - t.regularization()) * 5000000
//}
//bestEnergy := energy()
//prevEnergy := bestEnergy
//fmt.Println(bestEnergy)

//bestWeights := make([]Score, len(t.weights))
//prevWeights := make([]Score, len(t.weights))
//for i := range t.weights {
//bestWeights[i] = *t.weights[i]
//prevWeights[i] = *t.weights[i]
//}

//maxTemperature := 1000.0
//minTemperature := 2.5
//tFactor := -math.Log(maxTemperature / minTemperature)

//steps := 50000
//currentEnergy := bestEnergy

//for i := 1; i <= steps; i++ {
//for i := range t.weights {
//t.weights[i].Middle += int16(rand.NormFloat64() * 1)
//t.weights[i].End += int16(rand.NormFloat64() * 1)
//}
//loadScoresToPieceSquares()
//temperature := maxTemperature * math.Exp(tFactor*float64(i)/float64(steps))
//currentEnergy = energy()
//dEnergy := currentEnergy - prevEnergy
//fmt.Printf("D: %f Temperature: %f Error: %f Exp: %f\r", dEnergy, temperature, currentEnergy, math.Exp(-dEnergy/temperature))
//if dEnergy > 0 && math.Exp(-dEnergy/temperature) < rand.Float64() {
//// Previous state
//for i := range t.weights {
//*t.weights[i] = prevWeights[i]
//}
//} else {
//for i := range t.weights {
//prevWeights[i] = *t.weights[i]
//}
//prevEnergy = currentEnergy
//if currentEnergy < bestEnergy {
//fmt.Println("Yay!")
//fmt.Println(currentEnergy)
//fmt.Println(t.weights)
//improved = true
//for i := range t.weights {
//bestWeights[i] = *t.weights[i]
//}
//}
//}
//}
//fmt.Printf("\nlastEnergy: %f\n", currentEnergy)

//for i := range t.weights {
//*t.weights[i] = bestWeights[i]
//}
//if improved {
//fmt.Println(t.weights)
//}
//return
//}

func (t *tuner) coordinateDescent() {
	bestError := t.computeError()
	bestErrorWithRegularization := bestError + t.regularization()
	fmt.Printf("Initial values; error: %.17g; regularization: %.17g\n", bestError, t.regularization())
	fmt.Println(t.weights)

	indexes := make([]int, len(t.weights))
	for idx := range t.weights {
		indexes[idx] = idx
	}
	for iter := 0; ; iter++ {
		improved := false
		rand.Shuffle(len(indexes), func(i, j int) {
			indexes[i], indexes[j] = indexes[j], indexes[i]
		})
		for _, idx := range indexes {
			score := t.weights[idx]
			for phase := 0; phase < score.phaseCount(); phase++ {
				oldValue := score.get(phase)
				bestValue := oldValue
				// try increasing
				for i := int16(1); i <= 64; i *= 2 {
					score.set(phase, oldValue+i)
					newError := t.computeError()
					newErrorWithRegularization := newError + t.regularization()
					// First compare to prevent decreasing parameter just to lower regularization(as some parameters may be irrelevant in test positions)
					if newError < bestError && newErrorWithRegularization < bestErrorWithRegularization {
						bestValue = oldValue + i
						bestError = newError
						bestErrorWithRegularization = newErrorWithRegularization
						improved = true
					} else {
						score.set(phase, bestValue)
						break
					}
				}
				// try decreasing
				if bestValue == oldValue {
					for i := int16(1); i <= 64; i *= 2 {
						score.set(phase, oldValue-i)
						newError := t.computeError()
						newErrorWithRegularization := newError + t.regularization()
						if newError < bestError && newErrorWithRegularization < bestErrorWithRegularization {
							bestValue = oldValue - i
							bestError = newError
							bestErrorWithRegularization = newErrorWithRegularization
							improved = true
						} else {
							score.set(phase, bestValue)
							break
						}
					}
				}
			}
		}
		fmt.Printf("Iteration %d; error: %.17g; regularization: %.17g\n", iter+1, bestError, t.regularization())
		fmt.Println(t.weights)
		if !improved {
			break
		}
	}
}

const releventToPieceSquaresCount = 5 + 5*8*4 + 6*8 + 8*4

type EvaluationValue interface {
	phaseCount() int
	set(phase int, value int16)
	get(phase int) int16
	regularized() bool
}

type ScoreValue struct {
	*Score
}

func (ScoreValue) phaseCount() int {
	return 2
}

func (ScoreValue) regularized() bool {
	return true
}

func (sv ScoreValue) set(phase int, value int16) {
	if phase == 0 {
		sv.Score.Middle = value
	} else {
		sv.Score.End = value
	}
}

func (sv ScoreValue) get(phase int) int16 {
	if phase == 0 {
		return sv.Score.Middle
	}
	return sv.Score.End
}

type SingleValue struct {
	*int16
}

func (SingleValue) phaseCount() int {
	return 1
}

func (sv SingleValue) set(phase int, value int16) {
	*sv.int16 = value
}

func (sv SingleValue) get(phase int) int16 {
	return *sv.int16
}

func (SingleValue) regularized() bool {
	return false
}

func (sv SingleValue) String() string {
	return fmt.Sprintf("%d", *sv.int16)
}

func loadScoresToSlice() (res []EvaluationValue) {
	res = append(res, ScoreValue{&PawnValue})
	res = append(res, ScoreValue{&KnightValue})
	res = append(res, ScoreValue{&BishopValue})
	res = append(res, ScoreValue{&RookValue})
	res = append(res, ScoreValue{&QueenValue})

	for i := 2; i <= 6; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				res = append(res, ScoreValue{&pieceScores[i][y][x]})
			}
		}
	}
	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			res = append(res, ScoreValue{&pawnScores[y][x]})
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			res = append(res, ScoreValue{&pawnsConnected[y][x]})
		}
	}
	for y := 0; y < 9; y++ {
		res = append(res, ScoreValue{&mobilityBonus[0][y]})
	}
	for y := 0; y < 14; y++ {
		res = append(res, ScoreValue{&mobilityBonus[1][y]})
	}
	for y := 0; y < 15; y++ {
		res = append(res, ScoreValue{&mobilityBonus[2][y]})
	}
	for y := 0; y < 28; y++ {
		res = append(res, ScoreValue{&mobilityBonus[3][y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&passedFriendlyDistance[y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&passedEnemyDistance[y]})
	}
	for y := 0; y < 7; y++ {
		res = append(res, ScoreValue{&passedRank[y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&passedFile[y]})
	}
	res = append(res, ScoreValue{&isolated})
	res = append(res, ScoreValue{&doubled})
	res = append(res, ScoreValue{&backward})
	res = append(res, ScoreValue{&backwardOpen})
	res = append(res, ScoreValue{&bishopPair})
	res = append(res, ScoreValue{&bishopRammedPawns})
	res = append(res, ScoreValue{&bishopOutpostUndefendedBonus})
	res = append(res, ScoreValue{&bishopOutpostDefendedBonus})
	res = append(res, ScoreValue{&knightOutpostUndefendedBonus})
	res = append(res, ScoreValue{&knightOutpostDefendedBonus})
	res = append(res, ScoreValue{&minorBehindPawn})
	res = append(res, ScoreValue{&tempo})
	res = append(res, ScoreValue{&rookOnFile[0]})
	res = append(res, ScoreValue{&rookOnFile[1]})
	for y := 0; y < 12; y++ {
		res = append(res, ScoreValue{&kingDefenders[y]})
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, ScoreValue{&kingShelter[x][y][z]})
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, ScoreValue{&kingStorm[x][y][z]})
			}
		}
	}
	for x := Knight; x <= Queen; x++ {
		res = append(res, SingleValue{&kingSafetyAttacksWeights[x]})
	}
	res = append(res, SingleValue{&kingSafetyAttackValue})
	res = append(res, SingleValue{&kingSafetyWeakSquares})
	res = append(res, SingleValue{&kingSafetyFriendlyPawns})
	res = append(res, SingleValue{&kingSafetyNoEnemyQueens})
	res = append(res, SingleValue{&kingSafetySafeQueenCheck})
	res = append(res, SingleValue{&kingSafetySafeRookCheck})
	res = append(res, SingleValue{&kingSafetySafeBishopCheck})
	res = append(res, SingleValue{&kingSafetySafeKnightCheck})
	res = append(res, SingleValue{&kingSafetyAdjustment})
	return
}
