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
	pos.Colours[utils.BoolToInt(unicode.IsUpper(piece))] |= bit
	switch byte(unicode.ToLower(piece)) {
	case 'p':
		pos.Pieces[Pawn] |= bit
	case 'r':
		pos.Pieces[Rook] |= bit
	case 'n':
		pos.Pieces[Knight] |= bit
	case 'b':
		pos.Pieces[Bishop] |= bit
	case 'q':
		pos.Pieces[Queen] |= bit
	case 'k':
		pos.Pieces[King] |= bit
	}
}
