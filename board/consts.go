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

type File int

const (
	A File = iota
	B
	C
	D
	E
	F
	G
	H
	FNone
)

type Rank int

const (
	One Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	RNone
)

type Color int

const (
	White Color = iota
	Black
	Both
)

type Square int

const (
	A1 Square = 21 + iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	_
	_

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	_
	_

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	_
	_

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	_
	_

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	_
	_

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	_
	_

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	_
	_

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	NoSquare
)

type Castling uint8

const (
	WKSide Castling = 1 << iota
	WQSide
	BKSide
	BQSide
)
