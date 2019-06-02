package evaluation

import (
	"bufio"
	"fmt"
	. "github.com/mhib/combusken/backend"
	"math"
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

func (t *thread) quiescence(alpha, beta, height int, inCheck bool) int {
	t.stack[height].pv.clear()
	pos := &t.stack[height].position

	if height >= 127 {
		return 0
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
	var entries []tuneEntry
	for entry := range resultChan {
		entries = append(entries, entry)
	}
	fmt.Println("Number of entries:")
	fmt.Println(len(entries))
	k := calculateOptimalK(entries)
	fmt.Printf("Optimal k: %.17g\n", k)
	weights := loadScoresToSlice()
	coordinateDescent(weights, entries, k)
	fmt.Println("Generated weights:")
	fmt.Println(weights)
}

func computeError(entries []tuneEntry, k float64) float64 {
	numCPU := runtime.NumCPU()
	results := make([]float64, numCPU)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for y := idx; y < len(entries); y += numCPU {
				entry := entries[y]
				evaluation := float64(Evaluate(&entry.Position))
				if !entry.Position.WhiteMove {
					evaluation *= -1
				}
				diff := entry.result - sigmoid(k, evaluation)
				results[idx] += diff * diff
			}
		}(i)
	}
	wg.Wait()
	sum := 0.0
	for _, tResult := range results {
		sum += tResult
	}
	return sum / float64(len(entries))
}

func calculateOptimalK(entries []tuneEntry) float64 {
	start := -10.0
	end := 10.0
	delta := 1.0
	best := computeError(entries, start)
	for i := 0; i < 10; i++ {
		current := start - delta
		for current < end {
			current += delta
			err := computeError(entries, current)
			if err <= best {
				best = err
				start = current
			}
		}
		end = start + delta
		start = start - delta
		delta /= 10
		fmt.Printf("Optimal k after %d steps: %.17g\n", i+1, start)
	}
	return start
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

func regularization(weights []*Score) float64 {
	alpha := 0.2e-7
	sum := 0
	for _, score := range weights {
		sum += absScore(score.Middle)
		sum += absScore(score.End)
	}
	return alpha * float64(sum)
}

func coordinateDescent(weights []*Score, entries []tuneEntry, k float64) {
	getter := func(score *Score, phase int) int16 {
		if phase == 0 {
			return score.Middle
		}
		return score.End
	}
	setter := func(score *Score, phase int, value int16, idx int) {
		if phase == 0 {
			score.Middle = value
		} else {
			score.End = value
		}
		if idx < releventToPieceSquaresCount {
			loadScoresToPieceSquares()
		}
	}
	bestError := computeError(entries, k)
	bestErrorWithRegularization := bestError + regularization(weights)
	fmt.Printf("Inital values; error: %.17g; regularization: %.17g\n", bestError, regularization(weights))
	fmt.Println(weights)
	for iter := 0; ; iter++ {
		improved := false
		for idx, score := range weights {
			for phase := 0; phase < 2; phase++ {
				oldValue := getter(score, phase)
				bestValue := oldValue
				for i := int16(1); i <= 64; i *= 2 {
					setter(score, phase, oldValue+i, idx)
					newError := computeError(entries, k)
					newErrorWithRegularization := newError + regularization(weights)
					// First compare to prevent decreasing parameter just to lower regularization(as some parameters may be irrelevant in test positions)
					if newError < bestError && newErrorWithRegularization < bestErrorWithRegularization {
						bestValue = oldValue + i
						bestError = newError
						bestErrorWithRegularization = newErrorWithRegularization
						improved = true
					} else {
						setter(score, phase, bestValue, idx)
						break
					}
				}
				if bestValue == oldValue {
					for i := int16(1); i <= 64; i *= 2 {
						setter(score, phase, oldValue-i, idx)
						newError := computeError(entries, k)
						newErrorWithRegularization := newError + regularization(weights)
						if newError < bestError && newErrorWithRegularization < bestErrorWithRegularization {
							bestValue = oldValue - i
							bestError = newError
							bestErrorWithRegularization = newErrorWithRegularization
							improved = true
						} else {
							setter(score, phase, bestValue, idx)
							break
						}
					}
				}
			}
		}
		fmt.Printf("Iteration %d; error: %.17g; regularization: %.17g\n", iter+1, bestError, regularization(weights))
		fmt.Println(weights)
		if !improved {
			break
		}
	}
}

const releventToPieceSquaresCount = 5 + 5*8*4 + 6*8 + 8*4

func loadScoresToSlice() (res []*Score) {
	res = append(res, &PawnValue)
	res = append(res, &KnightValue)
	res = append(res, &BishopValue)
	res = append(res, &RookValue)
	res = append(res, &QueenValue)

	for i := 2; i <= 6; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				res = append(res, &pieceScores[i][y][x])
			}
		}
	}
	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			res = append(res, &pawnScores[y][x])
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ {
			res = append(res, &pawnsConnected[y][x])
		}
	}
	for y := 0; y < 9; y++ {
		res = append(res, &mobilityBonus[0][y])
	}
	for y := 0; y < 14; y++ {
		res = append(res, &mobilityBonus[1][y])
	}
	for y := 0; y < 15; y++ {
		res = append(res, &mobilityBonus[2][y])
	}
	for y := 0; y < 28; y++ {
		res = append(res, &mobilityBonus[3][y])
	}
	for y := 0; y < 7; y++ {
		res = append(res, &passedRank[y])
	}
	for y := 0; y < 8; y++ {
		res = append(res, &passedFile[y])
	}
	res = append(res, &isolated)
	res = append(res, &doubled)
	res = append(res, &backward)
	res = append(res, &backwardOpen)
	res = append(res, &bishopPair)
	res = append(res, &bishopOutpostUndefendedBonus)
	res = append(res, &bishopOutpostDefendedBonus)
	res = append(res, &knightOutpostUndefendedBonus)
	res = append(res, &knightOutpostDefendedBonus)
	res = append(res, &minorBehindPawn)
	res = append(res, &tempo)
	res = append(res, &rookOnFile[0])
	res = append(res, &rookOnFile[1])
	res = append(res, &pawnShieldBonus[0])
	res = append(res, &pawnShieldBonus[1])

	return
}
