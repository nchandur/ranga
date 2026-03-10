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

// returns true the piece slides across the board
func (p *Piece) isSlider() bool {
	return *p == wB || *p == bB || *p == wR || *p == bR || *p == wQ || *p == bQ
}

// returns true if piece is major
func (p *Piece) isMajor() bool {
	return p.isQueen() || p.isRook()
}

// returns true if piece is minor
func (p *Piece) isMinor() bool {
	return p.isBishop() || p.isKnight()
}

// return material value of piece
func (p *Piece) Value() int {
	switch *p {
	case wK, bK:
		return 50000
	case wQ, bQ:
		return 1000
	case wR, bR:
		return 500
	case wN, bN, wB, bB:
		return 300
	case wP, bP:
		return 100

	default:
		return 0
	}
}

// returns color of piece
func (p *Piece) Color() Color {

	if *p > 0 && *p < 7 {
		return White
	}

	if *p >= 7 {
		return Black
	}

	return Both
}
