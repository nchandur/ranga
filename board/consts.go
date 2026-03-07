package board


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


type Castling uint8

const (
	WKSide Castling = 1 << iota
	WQSide
	BKSide
	BQSide
)
