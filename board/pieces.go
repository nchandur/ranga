package board

type Piece int

// w -> white b -> black
// P -> Pawn
// N -> Knight
// B -> Bishop
// R -> Rook
// Q -> Queen
// K -> King
const (
	Empty Piece = iota // empty square aka no piece
	wP
	wN
	wB
	wR
	wQ
	wK
	bP
	bN
	bB
	bR
	bQ
	bK
)

// checks if a given piece is a knight
func (piece *Piece) isKnight() bool {
	return *piece == wN || *piece == bN
}

// check if a given piece is a bishop
func (piece *Piece) isBishop() bool {
	return *piece == wB || *piece == bB
}

// check if a given piece is a rook
func (piece *Piece) isRook() bool {
	return *piece == wR || *piece == bR
}

// check if a given piece is a queen
func (piece *Piece) isQueen() bool {
	return *piece == wQ || *piece == bQ
}

// check if a given piece is a king
func (piece *Piece) isKing() bool {
	return *piece == wK || *piece == bK
}

// check if a given piece is a pawn
func (piece *Piece) isPawn() bool {
	return *piece == wP || *piece == bP
}
