package backend

import (
	"testing"
)

var visited map[int]struct{}

func testPosition(t *testing.T, pos *Position) {
	visited[int(pos.BishopFlag)] = struct{}{}
	if pos.Pieces[Bishop] == 0 {
		if pos.BishopFlag != 0 {
			t.Errorf("Expected none but got %d", pos.BishopFlag)
			t.FailNow()
		}
		return
	}

	if pos.Pieces[Bishop]&pos.Colours[Black]&BLACK_SQUARES != 0 {
		if pos.BishopFlag&BlackBlackSquareBishopFlag == 0 {
			t.Errorf("Expected flag %d but got %d", BlackBlackSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	} else {
		if pos.BishopFlag&BlackBlackSquareBishopFlag != 0 {
			t.Errorf("Expected not %d but got %d", BlackBlackSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	}

	if pos.Pieces[Bishop]&pos.Colours[Black]&WHITE_SQUARES != 0 {
		if pos.BishopFlag&BlackWhiteSquareBishopFlag == 0 {
			t.Errorf("Expected flag %d but got %d", BlackWhiteSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	} else {
		if pos.BishopFlag&BlackWhiteSquareBishopFlag != 0 {
			t.Errorf("Expected not %d but got %d", BlackWhiteSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	}

	if pos.Pieces[Bishop]&pos.Colours[White]&BLACK_SQUARES != 0 {
		if pos.BishopFlag&WhiteBlackSquareBishopFlag == 0 {
			t.Errorf("Expected flag %d but got %d", WhiteBlackSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	} else {
		if pos.BishopFlag&WhiteBlackSquareBishopFlag != 0 {
			t.Errorf("Expected not %d but got %d", WhiteBlackSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	}

	if pos.Pieces[Bishop]&pos.Colours[White]&WHITE_SQUARES != 0 {
		if pos.BishopFlag&WhiteWhiteSquareBishopFlag == 0 {
			t.Errorf("Expected flag %d but got %d", WhiteWhiteSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	} else {
		if pos.BishopFlag&WhiteWhiteSquareBishopFlag != 0 {
			t.Errorf("Expected not %d but got %d", WhiteWhiteSquareBishopFlag, pos.BishopFlag)
			t.FailNow()
		}
	}

}
func testAllInDepth(t *testing.T, pos *Position, depth int) {
	testPosition(t, pos)
	if depth == 1 || pos.Pieces[Bishop] == 0 {
		return
	}

	var child Position

	var buffer [1000]EvaledMove
	noisySize := GenerateNoisy(pos, buffer[:])
	quietsSize := GenerateQuiet(pos, buffer[noisySize:])

	for _, move := range buffer[:noisySize+quietsSize] {
		if pos.MakeMove(move.Move, &child) {
			testAllInDepth(t, &child, depth-1)
		}
	}
}

func TestBishopFlag(t *testing.T) {
	fen := "8/8/3k4/3bb3/3BB3/3K4/8/8 w - - 0 1"

	p := ParseFen(fen)
	visited = make(map[int]struct{})
	testAllInDepth(t, &p, 6)
	if len(visited) != 16 {
		t.Errorf("Some possibilities were not tested")
	}
}
