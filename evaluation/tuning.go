package evaluation

import "path/filepath"
import "bufio"
import "os"
import . "github.com/mhib/combusken/backend"

const (
	MG = iota
	EG
)

type tuneEntry struct {
	result    float64
	phase     float64
	midFactor float64
	endFactor float64
}

type thread struct {
	stack [127]StackEntry
}

type StackEntry struct {
	position Position
	PV
	moves [256]EvaledMove
}

type PV struct {
	size  int
	items [127]Move
}

func (pv *PV) clear() {
	pv.size = 0
}

func (pv *PV) assign(m Move, child *PV) {
	pv.size = 1 + child.size
	pv.items[0] = m
	copy(pv.items[1:], child.Moves())
}

func (pv *PV) Moves() []Move {
	return pv.items[:pv.size]
}

func (t *thread) quiescence(alpha, beta, height int, inCheck bool) int {
	t.stack[height].PV.clear()
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
			t.stack[height].PV.assign(evaled[i].Move, &t.stack[height+1].PV)
		}
	}

	if moveCount == 0 && inCheck {
		return -Mate + height
	}
	return alpha
}

func LoadEntries() (result []tuneEntry) {
	absPath, _ := filepath.Abs("./games.fen")
	file, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fen := scanner.Text()
	}

	return
}
