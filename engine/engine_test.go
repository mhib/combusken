package engine

import (
	"context"
	"fmt"
	. "github.com/mhib/combusken/backend"
	"strings"
	"testing"
)

func TestWAC(t *testing.T) {
	var good, bad int
	engine := NewEngine()
	engine.Threads.Val = 1
	engine.Hash.Val = 256
	for _, entry := range loadEPD("./test_positions/WinAtChess.epd") {
		engine.NewGame()
		result := engine.Search(context.Background(), SearchParams{Positions: []Position{entry.Position}, Limits: LimitsType{MoveTime: 1000}}) // search for 1 second
		found := false
		for _, move := range entry.bestMoves {
			if ParseMoveSAN(&entry.Position, move) == result {
				found = true
				break
			}
		}
		if found {
			good++
			fmt.Printf("#%v correct\n", entry.id)
		} else {
			//t.Errorf("#%v expected %v, got %v\n", entry.id, strings.Join(entry.bestMoves, " or "), result)
			fmt.Printf("#%v expected %v, got %v\n", entry.id, strings.Join(entry.bestMoves, " or "), result)
			bad++
		}
	}
	if bad != 0 {
		//t.Errorf("Failed %d out of %d", bad, good+bad) //
		fmt.Printf("Failed %d out of %d\n", bad, good+bad)
	}
}

// SAN parsing from CounterGO
func moveToSAN(pos *Position, ml []EvaledMove, mv Move) string {
	const PieceNames = "NBRQK"
	if mv == WhiteKingSideCastle || mv == BlackKingSideCastle {
		return "O-O"
	}
	if mv == WhiteQueenSideCastle || mv == BlackQueenSideCastle {
		return "O-O-O"
	}
	var strPiece, strCapture, strFrom, strTo, strPromotion string
	if pos.TypeOnSquare(SquareMask[mv.From()]) != Pawn {
		strPiece = string(PieceNames[pos.TypeOnSquare(SquareMask[mv.From()])-Knight])
	}
	strTo = SquareString[mv.To()]
	if pos.TypeOnSquare(SquareMask[mv.To()]) != None {
		strCapture = "x"
		if pos.TypeOnSquare(SquareMask[mv.From()]) == Pawn {
			strFrom = SquareString[mv.From()][:1]
		}
	}
	if mv.IsPromotion() {
		strPromotion = "=" + string(PieceNames[mv.PromotedPiece()-Knight])
	}
	var ambiguity = false
	var uniqCol = true
	var uniqRow = true
	for _, evaled := range ml {
		mv1 := evaled.Move
		if mv1.From() == mv.From() {
			continue
		}
		if mv1.To() != mv.To() {
			continue
		}

		if pos.TypeOnSquare(SquareMask[mv.From()]) != pos.TypeOnSquare(SquareMask[mv1.From()]) {
			continue
		}
		ambiguity = true
		if File(mv1.From()) == File(mv.From()) {
			uniqCol = false
		}
		if Rank(mv1.From()) == Rank(mv.From()) {
			uniqRow = false
		}
	}
	if ambiguity {
		if uniqCol {
			strFrom = SquareString[mv.From()][:1]
		} else if uniqRow {
			strFrom = SquareString[mv.From()][1:2]
		} else {
			strFrom = SquareString[mv.From()]
		}
	}
	return strPiece + strFrom + strCapture + strTo + strPromotion
}

func ParseMoveSAN(pos *Position, san string) Move {
	var index = strings.IndexAny(san, "+#?!")
	if index >= 0 {
		san = san[:index]
	}
	var ml = pos.GenerateAllLegalMoves()
	for _, mv := range ml {
		if san == moveToSAN(pos, ml, mv.Move) {
			return mv.Move
		}
	}
	return NullMove
}
