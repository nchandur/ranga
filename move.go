package main

// 0000 0000 0000 0000 0000 0111 1111 -> From Square (7 bits) 0x7F
// 0000 0000 0000 0011 1111 1000 0000 -> To Square (7 bits) >> 7, 0x7F
// 0000 0000 0011 1100 0000 0000 0000 -> Captured piece (4 bits) >> 14, 0xF
// 0000 0000 0100 0000 0000 0000 0000 -> Enpassant move (1 bit) 0x40000
// 0000 0000 1000 0000 0000 0000 0000 -> Pawn start (1 bit) 0x80000
// 0000 1111 1000 0000 0000 0000 0000 -> Promotion piece (4 bit) >> 20, 0xF
// 0001 0000 0000 0000 0000 0000 0000 -> Castle move (1 bit) 0x1000000

// flag for enpassant capture
var MFLAGEP = 0x40000

// flag for pawn start
var MFLAGPS = 0x80000

// flag for castle
var MFLAGCA = 0x100000

// flag for capture
var MFLAGCAP = 0x7c000

// flag for promotion
var MFLAGPROM = 0xf00000

// flag for no move
var MNONE = 0

type Move struct {
	move  int
	score int
}

func NewMove(from, to Square, captured, promoted Piece, flag int) Move {

	m := Move{}

	m.move = int(from) | (int(to) << 7) | (int(captured) << 14) | (int(promoted) << 20) | flag
	m.score = 0

	return m
}

// returns square from which piece moved
func (m *Move) FromSquare() Square {
	return Square(m.move & 0x7f)
}

// returns square to which piece moved
func (m *Move) ToSquare() Square {
	return Square((m.move >> 7) & 0x7f)
}

// returns piece captured
func (m *Move) Captured() Piece {
	return Piece((m.move >> 14) & 0xF)
}

// returns piece to which pawn was promoted
func (m *Move) Promoted() Piece {
	return Piece((m.move >> 20) & 0xF)
}
