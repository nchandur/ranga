package board

// checks if a given piece is a knight
func isKnight(piece Piece) bool {
	return piece == wN || piece == bN
}

// check if a given piece is a bishop
func isBishop(piece Piece) bool {
	return piece == wB || piece == bB
}

// check if a given piece is a rook
func isRook(piece Piece) bool {
	return piece == wR || piece == bR
}

// check if a given piece is a queen
func isQueen(piece Piece) bool {
	return piece == wQ || piece == bQ
}

// check if a given piece is a king
func isKing(piece Piece) bool {
	return piece == wK || piece == bK
}

// check if a given piece is a pawn
func isPawn(piece Piece) bool {
	return piece == wP || piece == bP
}
