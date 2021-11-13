package evaluation

import (
	"fmt"
	"testing"

	. "github.com/mhib/combusken/chess"
)

// From zurichess
var testFENs = []string{
	// Initial position
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	// Kiwipete: https://chessprogramming.wikispaces.com/Perft+Results
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	// Duplain: https://chessprogramming.wikispaces.com/Perft+Results
	"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	// Underpromotion: http://www.stmintz.com/ccc/index.php?id=366606
	"8/p1P5/P7/3p4/5p1p/3p1P1P/K2p2pp/3R2nk w - - 0 1",
	// Enpassant: http://www.10x8.net/chess/PerfT.html
	"8/7p/p5pb/4k3/P1pPn3/8/P5PP/1rB2RK1 b - d3 0 28",
	//// http://www.talkchess.com/forum/viewtopic.php?t=48609
	"1K1k4/8/5n2/3p4/8/1BN2B2/6b1/7b w - - 0 1",
	// http://www.talkchess.com/forum/viewtopic.php?t=51272
	"6k1/5ppp/3r4/8/3R2b1/8/5PPP/R3qB1K b - - 0 1",
	// http://www.stmintz.com/ccc/index.php?id=206056
	"2rqkb1r/p1pnpppp/3p3n/3B4/2BPP3/1QP5/PP3PPP/RN2K1NR w KQk - 0 1",
	// http://www.stmintz.com/ccc/index.php?id=60880
	"1rr3k1/4ppb1/2q1bnp1/1p2B1Q1/6P1/2p2P2/2P1B2R/2K4R w - - 0 1",
	// https://chessprogramming.wikispaces.com/SEE+-+The+Swap+Algorithm
	"1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1",
	"1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1",
	// http://www.talkchess.com/forum/viewtopic.php?topic_view=threads&p=419315&t=40054
	"8/8/3p4/4r3/2RKP3/5k2/8/8 b - - 0 1",
	// Pinned piece can give check: https://groups.google.com/forum/#!topic/fishcooking/S_4E_Xs5HaE
	"r2qk2r/pppb1ppp/2np4/1Bb5/4n3/5N2/PPP2PPP/RNBQR1K1 b kq - 1 1",
	// SAN test position: http://talkchess.com/forum/viewtopic.php?t=61393
	"Bn1N3R/ppPpNR1r/BnBr1NKR/k3pP2/3PR2R/N7/3P2P1/4Q2R w - e6 0 1",
	// zurichess: various
	"8/K5p1/1P1k1p1p/5P1P/2R3P1/8/8/8 b - - 0 78",
	"8/1P6/5ppp/3k1P1P/6P1/8/1K6/8 w - - 0 78",
	"1K6/1P6/5ppp/3k1P1P/6P1/8/8/8 w - - 0 1",
	"r1bqkb1r/ppp1pp2/2n3P1/3p4/3Pn3/5N1P/PPP1PPB1/RNBQK2R b KQkq - 0 1",
	"r1bqkb1r/ppp2p2/2n1p1pP/3p4/3Pn3/2N2N1P/PPP1PPB1/R1BQK2R b KQkq - 0 1",
	"r3kb2/ppp2pp1/6n1/7Q/8/2P1BN1b/1q2PPB1/3R1K1R b q - 0 1",
	"r7/1p4p1/2p2kb1/3r4/3N3n/4P2P/1p2BP2/3RK1R1 w - - 0 1",
	"r7/1p4p1/5k2/8/6P1/3Nn3/1p3P2/3BK3 w - - 0 1",
	"8/1p2k1p1/4P3/8/1p2N3/4P3/5P2/3BK3 b - - 0 1",
	"r1bk3r/ppp2p1p/4pp2/4n3/1b2P3/2N5/PPP2PPP/R3KBNR w KQ - 0 9",
	"rnb1kbnr/pp1ppppp/8/1q6/2PpP3/5N2/PP3PPP/RNBQ1K1R b kq c3 0 6",
	"1r2k2r/p5bp/4p1p1/q2pB1N1/6P1/6QP/1P6/2KR3R b k - 0 1",
	// zurichess: many captures
	"6k1/Qp1r1pp1/p1rP3p/P3q3/2Bnb1P1/1P3PNP/4p1K1/R1R5 b - - 0 1",
	"3r2k1/2Q2pb1/2n1r3/1p1p4/pB1PP3/n1N2p2/B1q2P1R/6RK b - - 0 1",
	"2r3k1/5p1n/6p1/pp3n2/2BPp2P/4P2P/q1rN1PQb/R1BKR3 b - - 0 1",
	"r3r3/bpp1Nk1p/p1bq1Bp1/5p2/PPP3n1/R7/3QBPPP/5RK1 w - - 0 1",
	"4r1q1/1p4bk/2pp2np/4N2n/2bp2pP/PR3rP1/2QBNPB1/4K2R b K - 0 1",
	// crafted:
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	"7k/8/8/8/1RRNN3/1BBQQ3/1KQQQ3/1QQQQ3 b - - 0 1",
	"rr2r1k1/ppBb1ppp/8/4p1NQ/8/1qB3B1/PP4PP/R5K1 w - - 0 1",
	// ECM
	"7r/1p2k3/2bpp3/p3np2/P1PR4/2N2PP1/1P4K1/3B4 b - - 0 1",
	"4k3/p1P3p1/2q1np1p/3N4/8/1Q3PP1/6KP/8 w - - 0 1",
	"3q4/pp3pkp/5npN/2bpr1B1/4r3/2P2Q2/PP3PPP/R4RK1 w - - 0 1",
	"4k3/p1P3p1/2q1np1p/3N4/8/1Q3PP1/7P/5K2 b - - 1 1",
	// Theban Chess
	"1p6/2p3kn/3p2pp/4pppp/5ppp/8/PPPPPPPP/PPPPPPKN w - - 0 1",
}

func see(pos *Position, mv Move) int {
	var from = mv.From()
	var to = mv.To()
	var piece = mv.MovedPiece()
	var side = pos.SideToMove
	var score = 0
	// special case for ep and castling
	if mv.Type() == EPCapture || mv.IsCastling() {
		return 0
	}
	if mv.CapturedPiece() != None {
		score += SEEValues[mv.CapturedPiece()]
	}
	if mv.IsPromotion() {
		piece = mv.PromotedPiece()
		score += SEEValues[piece] - SEEValues[Pawn]
	}
	pieces := (pos.Colours[Black] ^ pos.Colours[White] ^ SquareMask[from]) | SquareMask[to]
	score -= seeRec(pos, side^1, to, pieces, piece)
	return score
}

func seeRec(pos *Position, side int, to int, pieces uint64, lastPiece int) int {
	var maxScore = 0 // 0 if more captures are unprofitable
	var piece, from = getLeastValuableAttacker(pos, to, side, pieces)
	if from != NoSquare {
		var score = SEEValues[lastPiece]
		if lastPiece != King {
			score -= seeRec(pos, side^1, to, pieces&^SquareMask[from], piece)
		}
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

func calculatedSee(pos *Position, mv Move, expected int) int {
	for SeeAbove(pos, mv, expected) {
		expected++
	}
	return expected
}

func TestSEE(t *testing.T) {
	var buffer [256]EvaledMove
	good, bad := 0, 0
	for i, fen := range testFENs {
		pos := ParseFen(fen)
		quietSize := GenerateQuiet(&pos, buffer[:])
		noisySize := GenerateNoisy(&pos, buffer[quietSize:])
		for _, m := range buffer[:quietSize+noisySize] {
			expected := see(&pos, m.Move)
			shouldTrue := SeeAbove(&pos, m.Move, expected)
			shouldFalse := SeeAbove(&pos, m.Move, expected+1)

			if !shouldTrue || shouldFalse {
				t.Errorf("failed on on #%d %v", i, fen)
				fmt.Println(m.Move, expected, shouldTrue, shouldFalse, calculatedSee(&pos, m.Move, expected))
				bad++
			} else {
				good++
			}
		}
	}

	if bad != 0 {
		t.Errorf("Failed %d out of %d", bad, good+bad)
	}
}

// A benchmark position from http://www.stmintz.com/ccc/index.php?id=60880
var seeBench = "1rr3k1/4ppb1/2q1bnp1/1p2B1Q1/6P1/2p2P2/2P1B2R/2K4R w - - 0 1"

func BenchmarkSEESlow(b *testing.B) {
	var buffer [256]EvaledMove
	pos := ParseFen(seeBench)
	quietSize := GenerateQuiet(&pos, buffer[:])
	noisySize := GenerateNoisy(&pos, buffer[quietSize:])
	for i := 0; i < b.N; i++ {
		for _, m := range buffer[:quietSize+noisySize] {
			see(&pos, m.Move)
		}

	}
}

// pass SEE result as cutoff to seeAbove
var seeResults = []int{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	-1161,
	-459,
	-1289,
	-1289,
	0,
	0,
	0,
	-830,
	-702,
	-702,
	0,
	0,
	0,
	-830,
	-48,
	-830,
	459,
	-2401,
	0,
	0,
	0,
	-2529,
	-2401,
	-1747,
	-2401,
	-1699,
}

func BenchmarkSEEFast(b *testing.B) {
	var buffer [256]EvaledMove
	pos := ParseFen(seeBench)
	quietSize := GenerateQuiet(&pos, buffer[:])
	noisySize := GenerateNoisy(&pos, buffer[quietSize:])
	for i := 0; i < b.N; i++ {
		for idx, m := range buffer[:quietSize+noisySize] {
			SeeAbove(&pos, m.Move, seeResults[idx])
		}
	}
}
