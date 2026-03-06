package board

// checks if a given piece is a knight
func (b *Board) isKnight(piece Piece) bool {
	return piece == wN || piece == bN
}

// check if a given piece is a bishop
func (b *Board) isBishop(piece Piece) bool {
	return piece == wB || piece == bB
}

// check if a given piece is a rook
func (b *Board) isRook(piece Piece) bool {
	return piece == wR || piece == bR
}

// check if a given piece is a queen
func (b *Board) isQueen(piece Piece) bool {
	return piece == wQ || piece == bQ
}

// check if a given piece is a king
func (b *Board) isKing(piece Piece) bool {
	return piece == wK || piece == bK
}

// check if a given piece is a pawn
func (b *Board) isPawn(piece Piece) bool {
	return piece == wP || piece == bP
}
