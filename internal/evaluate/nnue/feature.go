package nnue

import "ranga/internal/board"

// computes the input feature index for one piece from one perspective
func FeatureIndex(perspectiveIsWhite, isUs bool, piece board.Piece, square board.Square, kSq board.Square) int {

	kFile := int(kSq) % 8

	if kFile >= 4 {
		square = board.Square(int(square) ^ 7)
	}

	relSquare := int(square)

	if perspectiveIsWhite {
		relSquare = relSquare ^ 56
	}

	offset := 0
	if !isUs {
		offset = 384
	}

	pieceType := int(piece) % 6

	return offset + pieceType*64 + relSquare
}
