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
func (p *Piece) isKnight() bool {
	return *p == wN || *p == bN
}

// check if a given piece is a bishop
func (p *Piece) isBishop() bool {
	return *p == wB || *p == bB
}

// check if a given piece is a rook
func (p *Piece) isRook() bool {
	return *p == wR || *p == bR
}

// check if a given piece is a queen
func (p *Piece) isQueen() bool {
	return *p == wQ || *p == bQ
}

// check if a given piece is a king
func (p *Piece) isKing() bool {
	return *p == wK || *p == bK
}

// check if a given piece is a pawn
func (p *Piece) isPawn() bool {
	return *p == wP || *p == bP
}
