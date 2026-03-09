package main

type Piece int

const (
	Empty Piece = iota
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

// returns true if piece is pawn
func (p *Piece) isPawn() bool {
	return *p == wP || *p == bP
}

// returns true if piece is knight
func (p *Piece) isKnight() bool {
	return *p == wN || *p == bN
}

// returns true if piece is bishop
func (p *Piece) isBishop() bool {
	return *p == wB || *p == bB
}

// returns true if piece is rook
func (p *Piece) isRook() bool {
	return *p == wR || *p == bR
}

// returns true if piece is queen
func (p *Piece) isQueen() bool {
	return *p == wQ || *p == bQ
}

// returns true if piece is king
func (p *Piece) isKing() bool {
	return *p == wK || *p == bK
}
