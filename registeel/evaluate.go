package registeel

import (
	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

const sideSize = (48 + 64*5)
const castlingRightsSize = 16

type RegisteelNetwork struct {
	inputBuffer [65]int16
	hiddenLayer [8]float32
	output      [2]float32
}

func (network *RegisteelNetwork) CorrectEvaluation(pos *Position) Score {
	for i := range network.hiddenLayer {
		network.hiddenLayer[i] = 0
	}
	inputSize := network.FillInput(pos)
	for i := 0; i < inputSize; i++ {
		input := network.inputBuffer[i]
		for output := range network.hiddenLayer {
			network.hiddenLayer[output] += firstLayerWeights[input][output]
		}

	}
	network.output[0] = outputBias[0]
	network.output[1] = outputBias[1]
	for input := range network.hiddenLayer {
		network.hiddenLayer[input] += firstLayerBias[input]
		if network.hiddenLayer[input] > 0 {
			network.output[0] += secondLayerWeights[input][0] * network.hiddenLayer[input]
			network.output[1] += secondLayerWeights[input][1] * network.hiddenLayer[input]
		}
	}

	return S(int16(network.output[0]), int16(network.output[1]))
}

func (network *RegisteelNetwork) FillInput(pos *Position) (bufferSize int) {
	var fromBB uint64
	for colour := Black; colour <= White; colour++ {
		sideOffset := sideSize * colour
		for fromBB = pos.Colours[colour] & pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			network.inputBuffer[bufferSize] = int16(sideOffset + BitScan(fromBB) - 8)
			bufferSize++
		}
		piece := Knight
		offset := sideOffset + 48
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			network.inputBuffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Bishop
		offset = sideOffset + 48 + 64
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			network.inputBuffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Rook
		offset = sideOffset + 48 + 64*2
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			network.inputBuffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		piece = Queen
		offset = sideOffset + 48 + 64*3
		for fromBB = pos.Colours[colour] & pos.Pieces[piece]; fromBB != 0; fromBB &= (fromBB - 1) {
			network.inputBuffer[bufferSize] = int16(offset + BitScan(fromBB))
			bufferSize++
		}

		offset = sideOffset + 48 + 64*4
		network.inputBuffer[bufferSize] = int16(offset + BitScan(pos.Colours[colour]&pos.Pieces[King]))
		bufferSize++
	}
	network.inputBuffer[bufferSize] = int16(sideSize*2 + uint(pos.Flags))
	bufferSize++
	return
}
