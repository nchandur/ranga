package board

import (
	"fmt"
)

//       binary move bits                               hexidecimal constants

// 0000 0000 0000 0000 0011 1111    source square       0x3f
// 0000 0000 0000 1111 1100 0000    target square       0xfc0
// 0000 0000 1111 0000 0000 0000    capture piece       0xf000
// 0000 1111 0000 0000 0000 0000    promoted piece      0xf0000
// 0001 0000 0000 0000 0000 0000    capture flag        0x100000
// 0010 0000 0000 0000 0000 0000    double push flag    0x200000
// 0100 0000 0000 0000 0000 0000    enpassant flag      0x400000
// 1000 0000 0000 0000 0000 0000    castling flag       0x800000

type Move uint32

const (
	SourceMask   = 0x3f
	TargetMask   = 0xfc0
	PieceMask    = 0xf000
	PromotedMask = 0xf0000
	CaptureFlag  = 0x100000
	DoubleFlag   = 0x200000
	EnPassFlag   = 0x400000
	CastleFlag   = 0x800000
)

const NOMOVE = Move(0)
const MAX_MOVES = 512

func NewMove(source, target Square, piece, promoted Piece, isCapture, isDouble, isEnpass, isCastle bool) Move {

	val := int(source)
	val |= int(target) << 6
	val |= int(piece) << 12
	val |= int(promoted) << 16

	if isCapture {
		val |= (1 << 20)
	}

	if isDouble {
		val |= (1 << 21)
	}

	if isEnpass {
		val |= (1 << 22)
	}

	if isCastle {
		val |= (1 << 23)
	}

	return Move(val)
}

// returns square from where the piece was moved
func (m Move) Source() Square {
	return Square(m & SourceMask)
}

// returns square to which the piece is moved
func (m Move) Target() Square {
	return Square((m & TargetMask) >> 6)
}

// returns the piece that was moved
func (m Move) Piece() Piece {
	return Piece((m & PieceMask) >> 12)
}

// returns the piece to which pawn promoted
func (m Move) Promoted() Piece {
	return Piece((m & PromotedMask) >> 16)
}

// returns if move was a capture
func (m Move) IsCapture() bool {
	return (m & CaptureFlag) != 0
}

// returns if move was a double pawN push
func (m Move) IsDouble() bool {
	return (m & DoubleFlag) != 0
}

// returns if move was an enpassant
func (m Move) IsEnpass() bool {
	return (m & EnPassFlag) != 0
}

// returns if move was a castle
func (m Move) IsCastle() bool {
	return (m & CastleFlag) != 0
}

// print move
func (m Move) String() string {

	if m == 0 {
		return ""
	}

	res := fmt.Sprintf("%s%s", m.Source(), m.Target())

	switch m.Promoted() {
	case 1, 7:
		res += "n"
	case 2, 8:
		res += "b"
	case 3, 9:
		res += "r"
	case 4, 10:
		res += "q"
	}

	return res
}
