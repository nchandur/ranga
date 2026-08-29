package nnue

import "ranga/internal/board"

// maps a piece and square to an index 0-767 based on perspective
func getFeatureIndex(piece board.Piece, sq board.Square, perspective board.Color) int {
	p := int(piece)
	s := int(sq)

	// black perspective
	if perspective == board.Black {
		s = s ^ 56

		if p < 6 {
			p += 6
		} else {
			p -= 6
		}
	}

	// 12 pieces * 64 squares = 768
	return p*64 + s
}

// applies Clipped ReLU activation
func crelu(x int16) int16 {
	if x < 0 {
		return 0
	}
	if x > ActivationLimit {
		return ActivationLimit
	}
	return x
}
