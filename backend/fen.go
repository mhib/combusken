package backend

import (
	"github.com/mhib/combusken/utils"
	"strconv"
	"strings"
	"unicode"
)

const InitialPositionFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func ParseFen(input string) Position {
	var res Position

	for i := 0; i < 64; i++ {
		res.Squares[i] = NoPiece
	}

	slices := strings.Split(input, " ")
	res.Flags |= WhiteKingSideCastleFlag | WhiteQueenSideCastleFlag |
		BlackKingSideCastleFlag | BlackQueenSideCastleFlag

	y := uint(7)
	x := uint(0)
	for _, char := range slices[0] {
		if char == '/' {
			y--
			x = 0
		} else if unicode.IsDigit(char) {
			num, _ := strconv.Atoi(string(char))
			x += uint(num)
		} else {
			insertPiece(&res, char, 1<<(y*8+x))
			x++
		}
	}

	res.SideToMove = utils.BoolToInt(slices[1] == "w")

	for _, char := range slices[2] {
		switch char {
		case 'K':
			res.Flags ^= WhiteKingSideCastleFlag
		case 'Q':
			res.Flags ^= WhiteQueenSideCastleFlag
		case 'q':
			res.Flags ^= BlackQueenSideCastleFlag
		case 'k':
			res.Flags ^= BlackKingSideCastleFlag
		}
	}

	if slices[3] != "-" {
		square := (int(slices[3][0]) - int('a')) + (int(slices[3][1])-int('1'))*8
		if res.SideToMove == White {
			res.EpSquare = square - 8
		} else {
			res.EpSquare = square + 8
		}
	}

	if len(slices) >= 5 {
		parsed, _ := strconv.Atoi(slices[4])
		res.FiftyMove = parsed
	}

	HashPosition(&res)

	return res
}

func insertPiece(pos *Position, piece rune, bit uint64) {
	var intSide int
	if unicode.IsUpper(piece) {
		intSide = 1
		pos.Colours[1] |= bit
	} else {
		intSide = 0
		pos.Colours[0] |= bit
	}
	switch byte(unicode.ToLower(piece)) {
	case 'p':
		pos.Squares[BitScan(bit)] = NewPiece(Pawn, intSide)
		pos.Pieces[Pawn] |= bit
	case 'r':
		pos.Squares[BitScan(bit)] = NewPiece(Rook, intSide)
		pos.Pieces[Rook] |= bit
	case 'n':
		pos.Squares[BitScan(bit)] = NewPiece(Knight, intSide)
		pos.Pieces[Knight] |= bit
	case 'b':
		pos.Squares[BitScan(bit)] = NewPiece(Bishop, intSide)
		pos.Pieces[Bishop] |= bit
	case 'q':
		pos.Squares[BitScan(bit)] = NewPiece(Queen, intSide)
		pos.Pieces[Queen] |= bit
	case 'k':
		pos.Squares[BitScan(bit)] = NewPiece(King, intSide)
		pos.Pieces[King] |= bit
	}
}
