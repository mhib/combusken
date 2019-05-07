package engine

import (
	"bufio"
	"context"
	"fmt"
	. "github.com/mhib/combusken/backend"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type positionTest struct {
	Position
	bestMoves []string
	id        string
}

func loadEPD(fileLocation string) (res []positionTest) {
	path, _ := filepath.Abs(fileLocation)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry positionTest
		positionAndRest := strings.Split(line, " bm ")
		entry.Position = ParseFen(positionAndRest[0])
		bestMovesAndID := strings.Split(positionAndRest[1], ";")
		entry.bestMoves = strings.Split(bestMovesAndID[0], " ")
		entry.id = bestMovesAndID[1]
		res = append(res, entry)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

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
	if mv.MovedPiece() != Pawn {
		strPiece = string(PieceNames[mv.MovedPiece()-Knight])
	}
	strTo = SquareString[mv.To()]
	if mv.CapturedPiece() != None {
		strCapture = "x"
		if mv.MovedPiece() == Pawn {
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
		if mv1.MovedPiece() != mv.MovedPiece() {
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
	var buffer [256]EvaledMove
	var ml = pos.GenerateAllLegalMoves(buffer[:])
	for _, mv := range ml {
		if san == moveToSAN(pos, ml, mv.Move) {
			return mv.Move
		}
	}
	return NullMove
}
