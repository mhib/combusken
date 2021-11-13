package registeel

import (
	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

const sideSize = (48 + 64*5)
const castlingRightsSize = 16

type RegisteelNetwork struct {
	hiddenLayer [16]int64
	output      [2]int64
}

func (network *RegisteelNetwork) Initialize(pos *Position) {
	for i := range network.hiddenLayer {
		network.hiddenLayer[i] = int64(firstLayerQuantizedBias[i])
	}

	var inputBuffer [65]int16
	inputSize := fillInput(pos, inputBuffer[:])

	for i := 0; i < inputSize; i++ {
		input := inputBuffer[i]
		for output := range network.hiddenLayer {
			network.hiddenLayer[output] += int64(firstLayerQuatizedWeights[input][output])
		}
	}
}

func roundDivide(number, divisor int64) int64 {
	if number < 0 {
		return (number - divisor/2) / divisor
	} else {
		return (number + divisor/2) / divisor
	}
}

func (network *RegisteelNetwork) CorrectEvaluation(pos *Position) Score {
	network.output[0] = int64(outputQuantizedBias[0])
	network.output[1] = int64(outputQuantizedBias[1])
	for input := range network.hiddenLayer {
		if network.hiddenLayer[input] > 0 {
			network.output[0] += int64(secondLayerQuantizedWeights[input][0]) * network.hiddenLayer[input]
			network.output[1] += int64(secondLayerQuantizedWeights[input][1]) * network.hiddenLayer[input]
		}
	}
	return S(int16(roundDivide(network.output[0], outputDivisor)), int16(roundDivide(network.output[1], outputDivisor)))
}

// func freshEvaluation(pos *Position) Score {
// 	var hidden [8]float32
// 	var output [2]float32
// 	for i := range hidden {
// 		hidden[i] = firstLayerBias[i]
// 	}
// 	var network RegisteelNetwork
// 	inputSize := network.FillInput(pos)
// 	for i := 0; i < inputSize; i++ {
// 		input := network.inputBuffer[i]
// 		for output := range hidden {
// 			hidden[output] += firstLayerWeights[input][output]
// 		}
// 	}
// 	output[0] = outputBias[0]
// 	output[1] = outputBias[1]
// 	for input := range hidden {
// 		if hidden[input] > 0 {
// 			output[0] += secondLayerWeights[input][0] * hidden[input]
// 			output[1] += secondLayerWeights[input][1] * hidden[input]
// 		}
// 	}
// 	return S(int16(math.RoundToEven(float64(output[0]))), int16(math.RoundToEven(float64(output[1]))))
// }

func (network *RegisteelNetwork) RemoveInput(input int) {
	for output := 0; output < len(network.hiddenLayer); output++ {
		network.hiddenLayer[output] -= int64(firstLayerQuatizedWeights[input][output])
	}
}
func (network *RegisteelNetwork) AddInput(input int) {
	for output := 0; output < len(network.hiddenLayer); output++ {
		network.hiddenLayer[output] += int64(firstLayerQuatizedWeights[input][output])
	}
}

func flagsIdx(flag uint8) int {
	return sideSize*2 + int(flag)
}

func pieceIdx(square, piece, colour int) int {
	sideOffset := sideSize * colour
	if piece == Pawn {
		return sideOffset + square - 8
	}
	return sideOffset + 48 + 64*(piece-Knight) + square
}

func (network *RegisteelNetwork) ApplyMove(move Move, from, to *Position) {
	side := from.SideToMove
	if !move.IsPromotion() {
		network.RemoveInput(pieceIdx(move.From(), move.MovedPiece(), side))
		network.AddInput(pieceIdx(move.To(), move.MovedPiece(), side))
		switch move.Type() {
		case Capture:
			network.RemoveInput(pieceIdx(move.To(), move.CapturedPiece(), side^1))
		case KingCastle:
			if side == White {
				network.RemoveInput(pieceIdx(H1, Rook, side))
				network.AddInput(pieceIdx(F1, Rook, side))
			} else {
				network.RemoveInput(pieceIdx(H8, Rook, side))
				network.AddInput(pieceIdx(F8, Rook, side))
			}
		case QueenCastle:
			if side == White {
				network.RemoveInput(pieceIdx(A1, Rook, side))
				network.AddInput(pieceIdx(D1, Rook, side))
			} else {
				network.RemoveInput(pieceIdx(A8, Rook, side))
				network.AddInput(pieceIdx(D8, Rook, side))
			}
		case EPCapture:
			network.RemoveInput(pieceIdx(from.EpSquare, Pawn, side^1))
		default:
		}
	} else {
		network.RemoveInput(pieceIdx(move.From(), Pawn, side))
		network.AddInput(pieceIdx(move.To(), move.PromotedPiece(), side))
		if move.IsCapture() {
			network.RemoveInput(pieceIdx(move.To(), move.CapturedPiece(), side^1))
		}
	}
	if from.Flags != to.Flags {
		network.RemoveInput(flagsIdx(from.Flags))
		network.AddInput(flagsIdx(to.Flags))
	}
}

func (network *RegisteelNetwork) RevertMove(move Move, from, to *Position) {
	side := from.SideToMove
	if !move.IsPromotion() {
		network.AddInput(pieceIdx(move.From(), move.MovedPiece(), side))
		network.RemoveInput(pieceIdx(move.To(), move.MovedPiece(), side))
		switch move.Type() {
		case Capture:
			network.AddInput(pieceIdx(move.To(), move.CapturedPiece(), side^1))
		case KingCastle:
			if side == White {
				network.AddInput(pieceIdx(H1, Rook, side))
				network.RemoveInput(pieceIdx(F1, Rook, side))
			} else {
				network.AddInput(pieceIdx(H8, Rook, side))
				network.RemoveInput(pieceIdx(F8, Rook, side))
			}
		case QueenCastle:
			if side == White {
				network.AddInput(pieceIdx(A1, Rook, side))
				network.RemoveInput(pieceIdx(D1, Rook, side))
			} else {
				network.AddInput(pieceIdx(A8, Rook, side))
				network.RemoveInput(pieceIdx(D8, Rook, side))
			}
		case EPCapture:
			network.AddInput(pieceIdx(from.EpSquare, Pawn, side^1))
		default:
		}
	} else {
		network.AddInput(pieceIdx(move.From(), Pawn, side))
		network.RemoveInput(pieceIdx(move.To(), move.PromotedPiece(), side))
		if move.IsCapture() {
			network.AddInput(pieceIdx(move.To(), move.CapturedPiece(), side^1))
		}
	}
	if from.Flags != to.Flags {
		network.AddInput(flagsIdx(from.Flags))
		network.RemoveInput(flagsIdx(to.Flags))
	}
}

func fillInput(pos *Position, buffer []int16) (bufferSize int) {
	var fromBB uint64
	for colour := Black; colour <= White; colour++ {
		sideOffset := sideSize * colour
		for fromBB = pos.Colours[colour] & pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			buffer[bufferSize] = int16(sideOffset + BitScan(fromBB) - 8)
			bufferSize++
		}
		piece := Knight
		offset := sideOffset + 48
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			buffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Bishop
		offset = sideOffset + 48 + 64
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			buffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Rook
		offset = sideOffset + 48 + 64*2
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			buffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Queen
		offset = sideOffset + 48 + 64*3
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			buffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		offset = sideOffset + 48 + 64*4
		buffer[bufferSize] = int16(offset + BitScan(pos.Colours[colour]&pos.Pieces[King]))
		bufferSize++
	}
	buffer[bufferSize] = int16(sideSize*2 + uint(pos.Flags))
	bufferSize++
	return
}
