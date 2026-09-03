package nnue

import "ranga/internal/board"

// tracks the hidden layer state from both perspectives
type Accumulator struct {
	White [HiddenSize]int16
	Black [HiddenSize]int16
}

// computes accumulator from scratch for given board state.
func (acc *Accumulator) Refresh(nn *Network, b *board.Board) {
	for i := range HiddenSize {
		acc.White[i] = nn.FeatureBias[i]
		acc.Black[i] = nn.FeatureBias[i]
	}

	for sq := range board.Square(64) {
		piece := b.Mailbox[sq]

		if piece != board.Empty {
			acc.AddPiece(nn, piece, sq)
		}
	}
}

// AddPiece adds the weights of a newly placed piece to the accumulator.
func (acc *Accumulator) AddPiece(nn *Network, piece board.Piece, sq board.Square) {
	wIdx := getFeatureIndex(piece, sq, board.White)
	bIdx := getFeatureIndex(piece, sq, board.Black)

	for i := range HiddenSize {
		acc.White[i] += nn.FeatureWeights[wIdx][i]
		acc.Black[i] += nn.FeatureWeights[bIdx][i]
	}
}

// RemovePiece subtracts the weights of a removed piece from the accumulator.
func (acc *Accumulator) RemovePiece(nn *Network, piece board.Piece, sq board.Square) {
	wIdx := getFeatureIndex(piece, sq, board.White)
	bIdx := getFeatureIndex(piece, sq, board.Black)

	for i := range HiddenSize {
		acc.White[i] -= nn.FeatureWeights[wIdx][i]
		acc.Black[i] -= nn.FeatureWeights[bIdx][i]
	}
}
