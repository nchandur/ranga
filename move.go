package main

import (
	"fmt"
	"strings"
)

// 0000 0000 0000 0000 0000 0111 1111 -> From Square (7 bits) 0x7F
// 0000 0000 0000 0011 1111 1000 0000 -> To Square (7 bits) >> 7, 0x7F
// 0000 0000 0011 1100 0000 0000 0000 -> Captured piece (4 bits) >> 14, 0xF
// 0000 0000 0100 0000 0000 0000 0000 -> Enpassant move (1 bit) 0x40000
// 0000 0000 1000 0000 0000 0000 0000 -> Pawn start (1 bit) 0x80000
// 0000 1111 1000 0000 0000 0000 0000 -> Promotion piece (4 bit) >> 20, 0xF
// 0001 0000 0000 0000 0000 0000 0000 -> Castle move (1 bit) 0x1000000

const (
	MFLAGEP   = 0x40000   // flag for enpassant capture
	MFLAGPS   = 0x80000   // flag for pawn start
	MFLAGCA   = 0x1000000 // flag for castle
	MFLAGCAP  = 0x7C000   // flag for capture
	MFLAGPROM = 0xF00000  // flag for promotion
	MNONE     = 0         // flag for no move
)

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

// string
func (m *Move) String() string {
	var b strings.Builder

	fromSq := m.FromSquare()
	toSq := m.ToSquare()

	fmt.Fprintf(&b, "%s%s", fromSq.String(), toSq.String())

	prom := m.Promoted()

	promStr := ""

	if prom != Empty {

		if isKnight[prom] {
			promStr = "n"
		} else if isBishop[prom] {
			promStr = "b"
		} else if isRook[prom] {
			promStr = "r"
		} else if isQueen[prom] {
			promStr = "q"
		}

	}

	fmt.Fprintf(&b, "%s", promStr)

	return b.String()
}
