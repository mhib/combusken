import chess

import numpy as np
import pandas as pd

SIDE_SIZE = (48 + 64 * 5)
CASTLING_RIGHTS_SIZE = 16
SIZE = SIDE_SIZE * 2 + CASTLING_RIGHTS_SIZE

def rank(square):
	return square >> 3

def file(square):
	return square & 7

def distance(left, right):
	return max(abs(rank(left) - rank(right)), abs(file(left) - file(right)))

def castling_idx(board):
	idx = 0
	if (board.castling_rights & chess.BB_H1) == 0:
		idx |= 1
	if (board.castling_rights & chess.BB_A1) == 0:
		idx |= 2
	if (board.castling_rights & chess.BB_H8) == 0:
		idx |= 4
	if (board.castling_rights & chess.BB_A8) == 0:
		idx |= 8
	return idx + SIDE_SIZE * 2

fen_idx = 0
def vectorize_fen(fen):
	global fen_idx
	fen_idx += 1
	if fen_idx % 100_000 == 0:
		print(fen_idx)
	board = chess.Board(fen)
	res = []
	for square, piece in board.piece_map().items():
		piece_type = piece.piece_type
		offset = SIDE_SIZE * piece.color
		if piece_type == chess.PAWN:
			res.append(offset + (square - 8))
			continue
		offset += 48
		res.append(offset + 64 * (piece_type - 2) + square)
	res.append(castling_idx(board))
	return res


def idxs_to_vector(idxs, vector_size):
	vector = [0] * vector_size
	for el in idxs:
		vector[el] = 1
	return np.array(vector, dtype=bool)

if __name__ == '__main__':
	df = pd.read_csv('nn_input_with_phase_and_scale.csv', header=None, delimiter=';', names=['position', 'result', 'static', 'phase', 'scale'])
	df['position'] = df['position'].map(lambda x: np.array(idxs_to_vector(vectorize_fen(x), SIZE)))
	df.to_pickle('nn_input_with_phase_and_scale.gzip')
