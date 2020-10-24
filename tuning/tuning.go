package tuning

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	. "github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/engine"
	. "github.com/mhib/combusken/evaluation"
	. "github.com/mhib/combusken/utils"
)

type tuneEntry struct {
	Position
	result float64
}

type thread struct {
	stack [127]stackEntry
}

type stackEntry struct {
	engine.MoveProvider
	position Position
	pv
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

var moveHistory engine.MoveHistory

// Copy of quiescence search to extract quiet position
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
		t.stack[height].MoveProvider.InitSingular()
	} else {
		if val >= beta {
			return beta
		}
		if alpha < val {
			alpha = val
		}

		t.stack[height].MoveProvider.InitQs()
	}

	for i := range evaled {
		move := t.stack[height].GetNextMove(pos, &moveHistory, 128, height)
		if move == NullMove {
			break
		}
		if !pos.MakeMove(move, child) {
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
	k                       float64
	weights                 []EvaluationValue
	entries                 []tuneEntry
	bestWeights             []EvaluationValue
	bestError               float64
	bestErrorRegularization float64
	done                    bool
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
	t := &tuner{done: false}
	for entry := range resultChan {
		t.entries = append(t.entries, entry)
	}
	fmt.Println("Number of entries:")
	fmt.Println(len(t.entries))
	t.calculateOptimalK()
	fmt.Printf("Optimal k: %.17g\n", t.k)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Printf("\nBest values; error: %.17g; regularization: %.17g\n", t.bestError, t.bestErrorRegularization)
		fmt.Println(t.bestWeights)
		t.done = true
	}()

	t.weights = loadScoresToSlice()
	t.saveEvaluationValues()
	for {
		res := t.coordinateDescent()
		if t.done {
			return
		}
		// After coordinate descent current weights are the best, so no need to reload weights
		res = res || t.gradientDescent()
		if t.done {
			return
		}
		if !res {
			break
		}
		// After gradient descent current weights are probably not best weights
		t.loadEvaluationValues()
	}
	fmt.Printf("\nBest values; error: %.17g; regularization: %.17g\n", t.bestError, t.bestErrorRegularization)
	fmt.Println(t.bestWeights)

}

func (t *tuner) computeError(entriesCount int) float64 {
	numCPU := runtime.NumCPU()
	results := make([]float64, numCPU)
	wg := &sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var c, sum float64
			for y := idx; y < entriesCount; y += numCPU {
				entry := t.entries[y]
				evaluation := float64(Evaluate(&entry.Position))
				if entry.Position.SideToMove == Black {
					evaluation *= -1
				}
				diff := entry.result - sigmoid(t.k, evaluation)

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
	return sum / float64(entriesCount)
}

func (t *tuner) calculateOptimalK() {
	start := -10.0
	end := 10.0
	delta := 1.0
	t.k = start
	best := t.computeError(len(t.entries))
	for i := 0; i < 10; i++ {
		t.k = start - delta
		for t.k < end {
			t.k += delta
			err := t.computeError(len(t.entries))
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
	alpha := 0.2e-8
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

func (t *tuner) coordinateDescent() bool {
	anyImprovements := false
	t.bestError = t.computeError(len(t.entries))
	t.bestErrorRegularization = t.regularization()
	fmt.Printf("Initial values; error: %.17g; regularization: %.17g\n", t.bestError, t.bestErrorRegularization)
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
				if t.done {
					return false
				}
				oldValue := score.get(phase)
				bestValue := oldValue
				// try increasing
				for i := int16(1); i <= 64; i *= 2 {
					score.set(phase, oldValue+i)
					if idx < releventToPieceSquaresCount {
						LoadScoresToPieceSquares()
					}
					newError := t.computeError(len(t.entries))
					newErrorRegularization := t.regularization()
					// First compare to prevent decreasing parameter just to lower regularization(as some parameters may be irrelevant in test positions)
					if newError < t.bestError && newError+newErrorRegularization < t.bestError+t.bestErrorRegularization {
						bestValue = oldValue + i
						t.bestError = newError
						t.bestErrorRegularization = newErrorRegularization
						improved = true
						anyImprovements = true
					} else {
						score.set(phase, bestValue)
						if idx < releventToPieceSquaresCount {
							LoadScoresToPieceSquares()
						}
						break
					}
				}
				// try decreasing
				if bestValue == oldValue {
					for i := int16(1); i <= 64; i *= 2 {
						score.set(phase, oldValue-i)
						if idx < releventToPieceSquaresCount {
							LoadScoresToPieceSquares()
						}
						newError := t.computeError(len(t.entries))
						newErrorRegularization := t.regularization()
						if newError < t.bestError && newError+newErrorRegularization < t.bestError+t.bestErrorRegularization {
							bestValue = oldValue - i
							t.bestError = newError
							t.bestErrorRegularization = newErrorRegularization
							improved = true
							anyImprovements = true
						} else {
							score.set(phase, bestValue)
							if idx < releventToPieceSquaresCount {
								LoadScoresToPieceSquares()
							}
							break
						}
					}
				}
			}
		}
		t.saveEvaluationValues()
		fmt.Printf("Iteration %d; error: %.17g; regularization: %.17g\n", iter+1, t.bestError, t.bestErrorRegularization)
		fmt.Println(t.weights)
		if !improved {
			break
		}
	}
	return anyImprovements
}

func (t *tuner) symmetricDerivative(batchSize, idx, phase, h int) float64 {
	weight := t.weights[idx]
	oldValue := weight.get(phase)

	weight.set(phase, oldValue+1)
	if idx < releventToPieceSquaresCount {
		LoadScoresToPieceSquares()
	}
	newError1 := t.computeError(batchSize) + t.regularization()

	weight.set(phase, oldValue-1)
	if idx < releventToPieceSquaresCount {
		LoadScoresToPieceSquares()
	}
	newError2 := t.computeError(batchSize) + t.regularization()

	weight.set(phase, oldValue)
	if idx < releventToPieceSquaresCount {
		LoadScoresToPieceSquares()
	}

	return (newError1 - newError2) / (2.0 * float64(h))
}

func (t *tuner) richardsExtrapolationDerivative(batchSize, idx, phase, h int) float64 {
	return (4.0*t.symmetricDerivative(batchSize, idx, phase, h) -
		t.symmetricDerivative(batchSize, idx, phase, 2*h)) /
		3.0
}

func (t *tuner) calculateGradient(batchSize int) []GradientVariable {
	res := make([]GradientVariable, 0, len(t.weights))
	for idx, weight := range t.weights {
		if t.done {
			break
		}
		gradient := newGradient(weight.phaseCount())
		for phase := 0; phase < weight.phaseCount(); phase++ {
			gradient.phases[phase] = t.symmetricDerivative(batchSize, idx, phase, 1)
		}
		res = append(res, gradient)
	}

	return res
}

func (t *tuner) gradientDescent() bool {
	anyImprovements := false
	const momentumRatio = 0.1
	prevGradient := make([]GradientVariable, len(t.weights))
	for idx, weight := range t.weights {
		prevGradient[idx].phases = make([]float64, weight.phaseCount())
	}
	t.bestError = t.computeError(len(t.entries))
	t.bestErrorRegularization = t.regularization()
	t.saveEvaluationValues()
	fmt.Printf("Initial values; error: %.17g; regularization: %.17g\n", t.bestError, t.bestErrorRegularization)
	fmt.Println(t.weights)
	batchSize := len(t.entries) / 10
	iterationsSinceImprovement := 0

	for iter := 0; iter < 20000; iter++ {
		rand.Shuffle(len(t.entries), func(i, j int) {
			t.entries[i], t.entries[j] = t.entries[j], t.entries[i]
		})
		gradient := t.calculateGradient(batchSize)
		anyChanges := false

		if t.done {
			break
		}

		max := 0.0
		for idx, weight := range t.weights {
			for phase := 0; phase < weight.phaseCount(); phase++ {
				max = math.Max(max, math.Abs(gradient[idx].phases[phase]))
			}
		}
		learningRate := 2.0 / max

		fmt.Println(max, learningRate)
		for idx, weight := range t.weights {
			for phase := 0; phase < weight.phaseCount(); phase++ {
				oldValue := weight.get(phase)
				diff := prevGradient[idx].phases[phase]*momentumRatio + learningRate*(1-momentumRatio)*gradient[idx].phases[phase]
				weight.set(phase, oldValue-int16(math.RoundToEven(diff)))
				prevGradient[idx].phases[phase] = diff
				if oldValue != weight.get(phase) {
					anyChanges = true
				}
			}
		}
		if t.done {
			break
		}
		iterationsSinceImprovement++
		currentError := t.computeError(len(t.entries))
		currentRegularization := t.regularization()
		fmt.Printf("Iteration %d; error: %.17g; regularization: %.17g\n", iter+1, currentError, currentRegularization)
		fmt.Println(t.weights)

		if currentError+currentRegularization < t.bestError+t.bestErrorRegularization {
			t.bestError = currentError
			t.bestErrorRegularization = currentRegularization
			t.saveEvaluationValues()
			iterationsSinceImprovement = 0
			anyImprovements = true
		}

		if !anyChanges || iterationsSinceImprovement > 50 {
			break
		}
	}
	fmt.Printf("error: %.17g; regularization: %.17g\n", t.bestError, t.bestErrorRegularization)
	fmt.Println(t.bestWeights)
	return anyImprovements
}

const releventToPieceSquaresCount = 5 + 5*8*4 + 6*8 + 8*4

type EvaluationValue interface {
	phaseCount() int
	set(phase int, value int16)
	get(phase int) int16
	regularized() bool
}

func (t *tuner) saveEvaluationValues() {
	res := make([]EvaluationValue, len(t.weights))
	for idx := range t.weights {
		tmp, _ := copyEvaluationValue(t.weights[idx])
		res[idx] = tmp
	}
	t.bestWeights = res
}

func (t *tuner) loadEvaluationValues() {
	for idx, weight := range t.weights {
		for phase := 0; phase < weight.phaseCount(); phase++ {
			weight.set(phase, t.bestWeights[idx].get(phase))
		}
	}
	LoadScoresToPieceSquares()
}

func copyEvaluationValue(ev EvaluationValue) (EvaluationValue, error) {
	switch v := ev.(type) {
	case ScoreValue:
		tmp := S(ev.get(0), ev.get(1))
		return ScoreValue{&tmp}, nil
	case SingleValue:
		val := ev.get(0)
		return SingleValue{&val}, nil
	default:
		return nil, fmt.Errorf("Unknown type %s", v)
	}
}

type GradientVariable struct {
	phases []float64
}

func newGradient(counter int) GradientVariable {
	return GradientVariable{
		phases: make([]float64, counter),
	}
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
		*sv.Score = S(value, sv.Score.End())
	} else {
		*sv.Score = S(sv.Score.Middle(), value)
	}
}

func (sv ScoreValue) get(phase int) int16 {
	if phase == 0 {
		return sv.Score.Middle()
	}
	return sv.Score.End()
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

	for i := Knight; i <= King; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 4; x++ {
				res = append(res, ScoreValue{&PieceScores[i][y][x]})
			}
		}
	}
	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			res = append(res, ScoreValue{&PawnScores[y][x]})
		}
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 4; x++ {
			res = append(res, ScoreValue{&PawnsConnected[y][x]})
		}
	}
	for y := 0; y < 9; y++ {
		res = append(res, ScoreValue{&MobilityBonus[0][y]})
	}
	for y := 0; y < 14; y++ {
		res = append(res, ScoreValue{&MobilityBonus[1][y]})
	}
	for y := 0; y < 15; y++ {
		res = append(res, ScoreValue{&MobilityBonus[2][y]})
	}
	for y := 0; y < 28; y++ {
		res = append(res, ScoreValue{&MobilityBonus[3][y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&PassedFriendlyDistance[y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&PassedEnemyDistance[y]})
	}
	for y := 0; y < 7; y++ {
		res = append(res, ScoreValue{&PassedRank[y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&PassedFile[y]})
	}
	for y := 0; y < 8; y++ {
		res = append(res, ScoreValue{&PassedStacked[y]})
	}
	res = append(res, ScoreValue{&Isolated})
	res = append(res, ScoreValue{&Doubled})
	res = append(res, ScoreValue{&Backward})
	res = append(res, ScoreValue{&BackwardOpen})
	res = append(res, ScoreValue{&BishopPair})
	res = append(res, ScoreValue{&BishopRammedPawns})
	res = append(res, ScoreValue{&BishopOutpostUndefendedBonus})
	res = append(res, ScoreValue{&BishopOutpostDefendedBonus})
	res = append(res, ScoreValue{&LongDiagonalBishop})
	res = append(res, ScoreValue{&KnightOutpostUndefendedBonus})
	res = append(res, ScoreValue{&KnightOutpostDefendedBonus})
	for y := 0; y < 4; y++ {
		res = append(res, ScoreValue{&DistantKnight[y]})
	}
	res = append(res, ScoreValue{&MinorBehindPawn})
	res = append(res, ScoreValue{&RookOnFile[0]})
	res = append(res, ScoreValue{&RookOnFile[1]})
	res = append(res, ScoreValue{&RookOnQueenFile})
	for y := 0; y < 12; y++ {
		res = append(res, ScoreValue{&KingDefenders[y]})
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, ScoreValue{&KingShelter[x][y][z]})
			}
		}
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 8; z++ {
				res = append(res, ScoreValue{&KingStorm[x][y][z]})
			}
		}
	}
	res = append(res, ScoreValue{&KingOnPawnlessFlank})
	for x := Knight; x <= Queen; x++ {
		res = append(res, SingleValue{&KingSafetyAttacksWeights[x]})
	}
	res = append(res, SingleValue{&KingSafetyAttackValue})
	res = append(res, SingleValue{&KingSafetyWeakSquares})
	res = append(res, SingleValue{&KingSafetyFriendlyPawns})
	res = append(res, SingleValue{&KingSafetyNoEnemyQueens})
	res = append(res, SingleValue{&KingSafetySafeQueenCheck})
	res = append(res, SingleValue{&KingSafetySafeRookCheck})
	res = append(res, SingleValue{&KingSafetySafeBishopCheck})
	res = append(res, SingleValue{&KingSafetySafeKnightCheck})
	res = append(res, SingleValue{&KingSafetyAdjustment})
	res = append(res, SingleValue{&KingSafetyMiddleDivisor})
	res = append(res, SingleValue{&KingSafetyEndDivisor})
	res = append(res, SingleValue{&Tempo})
	res = append(res, ScoreValue{&Hanging})
	res = append(res, ScoreValue{&ThreatByKing})
	for x := Pawn; x <= King; x++ {
		res = append(res, ScoreValue{&ThreatByMinor[x]})
	}
	for x := Pawn; x <= King; x++ {
		res = append(res, ScoreValue{&ThreatByRook[x]})
	}

	return
}
