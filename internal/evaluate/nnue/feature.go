package nnue

import "ranga/internal/board"

// computes the input feature index for one piece from one perspective
func FeatureIndex(perspectiveIsWhite, isUs bool, piece board.Piece, square board.Square) int {
	relSquare := square
	if !perspectiveIsWhite {
		relSquare = square ^ 56
	}
	offset := 0
	if !isUs {
		offset = 384
	}

	pieceType := int(piece) % 6

	return offset + pieceType*64 + int(relSquare)
}
