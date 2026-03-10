package main

const (
	BOARD_SQUARE_NUM = 120
)

type Board struct {
	Pieces     []Piece // represents pieces on the 120-index board. the square on which a piece exists has its corresponding value
	SideToMove Color
	FiftyMove  int
	Ply        int
	HistoryPly int
	Castling   CastleBit 
}
