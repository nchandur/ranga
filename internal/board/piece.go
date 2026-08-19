package board

type Piece uint8

const (
	WP Piece = iota
	WN
	WB
	WR
	WQ
	WK
	BP
	BN
	BB
	BR
	BQ
	BK

	Empty
)
