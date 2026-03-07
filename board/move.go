package board

import "fmt"

type Move struct {
	move  int
	score int
}

// 0000 0000 0000 0000 0000 0111 1111 -> From Square (7 bits) 0x7F
// 0000 0000 0000 0011 1111 1000 0000 -> To Square (7 bits) >> 7, 0x7F
// 0000 0000 0011 1100 0000 0000 0000 -> Captured piece (4 bits) >> 14, 0xF
// 0000 0000 0100 0000 0000 0000 0000 -> Enpassant move (1 bit) 0x40000
// 0000 0000 1000 0000 0000 0000 0000 -> Pawn start (1 bit) 0x80000
// 0000 1111 1000 0000 0000 0000 0000 -> Promotion piece (4 bit) >> 20, 0xF
// 0001 0000 0000 0000 0000 0000 0000 -> Castle move (1 bit) 0x1000000

// set from square
func (m *Move) SetFromSq(from Square) {
	m.move |= int(from)
}

// set to square
func (m *Move) SetToSq(to Square) {
	m.move |= int(to) << 7
}

// set captured piece
func (m *Move) SetCapturedPiece(piece Piece) {
	m.move |= int(piece) << 14
}

// set promoted piece
func (m *Move) SetPromotedPiece(piece Piece) {
	m.move |= int(piece) << 20
}

// set enpassant move
func (m *Move) SetEnPassant() {
	m.move |= 0x40000
}

// set pawn start
func (m *Move) SetPawnStart() {
	m.move |= 0x80000
}

// set castle move
func (m *Move) SetCapture() {
	m.move |= 0x1000000
}

// returns 120-based index for from square
func (m *Move) FromSq() int {
	return m.move & 0x7F
}

// returns 120-based index for to square
func (m *Move) ToSq() int {
	return (m.move >> 7) & 0x7F
}

// returns the piece that was captured
func (m *Move) CapturedPiece() int {
	return (m.move >> 14) & 0xF
}

// returns the promoted piece
func (m *Move) PromotedPiece() int {
	return (m.move >> 20) & 0xF
}

// returns if en passant happened
func (m *Move) IsEnPassant() bool {
	return m.move&0x40000 == 1
}

// returns if pawn started
func (m *Move) IsPawnStart() bool {
	return m.move&0x80000 == 1
}

// returns if castle
func (m *Move) IsCastle() bool {
	return m.move&0x1000000 == 1
}

// returns if capture happened
func (m *Move) IsCapture() bool {
	return m.move&0x7c000 == 1
}

// returns if promotion
func (m *Move) IsPromotion() bool {
	return m.move&0xF00000 == 1
}

func (m *Move) String() string {

	from := Square(m.FromSq())
	to := Square(m.ToSq())

	promoted := Piece(m.PromotedPiece())

	if promoted != 1 {
		pieceChar := 'q'

		if promoted.isKnight() {
			pieceChar = 'n'
		}

		if promoted.isBishop() {
			pieceChar = 'b'
		}

		if promoted.isRook() {
			pieceChar = 'r'
		}
		return fmt.Sprintf("%s%s%c", from.String(), to.String(), pieceChar)

	}

	return fmt.Sprintf("%s%s", from.String(), to.String())

}
